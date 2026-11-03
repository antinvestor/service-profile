package business

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pitabwire/frame/v3/data"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/antinvestor/service-profile/apps/default/service/models"
	"github.com/antinvestor/service-profile/internal/outbox"
)

// A handle is lowercase letters and digits in words joined by single dots,
// e.g. "peter.bwire". It is shown with a leading "@" but stored without it.
const (
	handleMinLen   = 3
	handleMaxLen   = 40
	handleBaseMax  = 24
	handleSuffix   = 8
	handleFallback = "user"

	// profileHandleIndex is the unique index on profiles.handle (see
	// models.Profile).
	profileHandleIndex = "idx_profiles_handle"
	pgUniqueViolation  = "23505"
)

var (
	handleSeparator = regexp.MustCompile(`[^a-z0-9]+`)
	handlePattern   = regexp.MustCompile(`^[a-z0-9]+(\.[a-z0-9]+)*$`)

	// ErrInvalidHandle is returned for a handle that cannot exist.
	ErrInvalidHandle = errors.New("invalid handle")
)

// NormalizeHandle turns a handle as a person types it ("@Peter.Bwire ")
// into its stored form ("peter.bwire"), or returns ErrInvalidHandle.
func NormalizeHandle(handle string) (string, error) {
	h := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(handle), "@"))
	if len(h) < handleMinLen || len(h) > handleMaxLen || !handlePattern.MatchString(h) {
		return "", ErrInvalidHandle
	}
	return h, nil
}

// slugHandle reduces a display name to handle words: accents are dropped,
// anything that is not a letter or digit separates words, and the result is
// cut to handleBaseMax characters.
func slugHandle(name string) string {
	folded, _, err := transform.String(
		transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC), name)
	if err != nil {
		folded = name
	}
	words := handleSeparator.Split(strings.ToLower(folded), -1)
	slug := strings.Join(nonEmpty(words), ".")
	if len(slug) > handleBaseMax {
		slug = strings.TrimRight(slug[:handleBaseMax], ".")
	}
	return slug
}

func nonEmpty(words []string) []string {
	out := words[:0]
	for _, w := range words {
		if w != "" {
			out = append(out, w)
		}
	}
	return out
}

// handleBase is the handle a profile would like: its name as a slug, or
// "user" when it has no usable name.
func handleBase(props data.JSONMap) string {
	// The properties a new handle is built from, in order.
	for _, key := range []string{"name", "au_name"} {
		name, ok := props[key].(string)
		if !ok {
			continue
		}
		if slug := slugHandle(name); len(slug) >= handleMinLen {
			return slug
		}
	}
	return handleFallback
}

// handleCandidates lists the handles to try for a new profile, in order:
// the bare name, the name with the tail of the profile id, and the profile
// id itself, which no other profile can hold. A profile without a name
// never gets a bare "user".
func handleCandidates(base, profileID string) []string {
	suffix := handleSeparator.ReplaceAllString(strings.ToLower(profileID), "")
	if len(suffix) > handleSuffix {
		suffix = suffix[len(suffix)-handleSuffix:]
	}
	suffixed := []string{base + "." + suffix, profileID}
	if base == handleFallback {
		return suffixed
	}
	return append([]string{base}, suffixed...)
}

// createWithHandle gives p the first free handle candidate and creates it
// (with its account and facts) in one transaction. A candidate another
// request claimed in between fails on the unique index; the next one is
// tried.
func (pb *profileBusiness) createWithHandle(
	ctx context.Context,
	p *models.Profile,
	account *models.ProfileAccount,
	facts func(*models.Profile) []*outbox.Event,
) error {
	p.GenID(ctx)

	var lastErr error
	for _, handle := range handleCandidates(handleBase(p.Properties), p.GetID()) {
		taken, err := pb.profileRepo.HandleTaken(ctx, handle)
		if err != nil {
			return err
		}
		if taken {
			continue
		}

		p.Handle = handle
		err = pb.profileRepo.CreateWithAccount(ctx, p, account, facts)
		if !isHandleConflict(err) {
			return err
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = errors.New("no free handle for the profile")
	}
	return lastErr
}

func isHandleConflict(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation && pgErr.ConstraintName == profileHandleIndex
}
