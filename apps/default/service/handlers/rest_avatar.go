package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pitabwire/frame/v2/data"
	"github.com/pitabwire/frame/v2/security"
)

// Max avatar size accepted for data-URI storage on the profile row.
// Matches the widget's sanitizePictureUrl ceiling (~512 KiB data URLs).
const maxAvatarBytes = 400 * 1024

// RestUpdateAvatar handles PUT/POST multipart avatar uploads from @stawi/profile.
//
// Path (after gateway /profile prefix strip):
//
//	/profile.v1.ProfileService/UpdateAvatar/{profileId}
//
// Body: multipart form with a "file" part (auth-runtime upload() contract).
// Response: {"data": <ProfileObject>} where properties.au_avater_uri is set.
//
// Images are stored as data:image/*;base64,... on the profile property
// au_avater_uri. Larger durable storage via the files service can replace
// this later without changing the client contract.
func (ps *ProfileServer) RestUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ps.Service.Log(ctx).WithField("method", "RestUpdateAvatar")

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		w.Header().Set("Allow", "PUT, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims := security.ClaimsFromContext(ctx)
	if claims == nil {
		ps.writeError(ctx, w, fmt.Errorf("authorization required"), http.StatusUnauthorized)
		return
	}

	profileID := strings.TrimSpace(r.PathValue("profileId"))
	if profileID == "" {
		// Fallback for routers that don't set PathValue.
		const prefix = "/profile.v1.ProfileService/UpdateAvatar/"
		if i := strings.Index(r.URL.Path, prefix); i >= 0 {
			profileID = strings.Trim(r.URL.Path[i+len(prefix):], "/")
		}
	}
	if profileID == "" {
		ps.writeError(ctx, w, fmt.Errorf("profile id required"), http.StatusBadRequest)
		return
	}

	subject, _ := claims.GetSubject()
	if subject != profileID {
		if err := ps.checker.Check(ctx, "profile_update"); err != nil {
			ps.writeError(ctx, w, err, http.StatusForbidden)
			return
		}
	}

	if err := r.ParseMultipartForm(maxAvatarBytes + 64*1024); err != nil {
		log.WithError(err).Warn("avatar: multipart parse failed")
		ps.writeError(ctx, w, fmt.Errorf("could not parse multipart body"), http.StatusBadRequest)
		return
	}

	file, hdr, err := r.FormFile("file")
	if err != nil {
		ps.writeError(ctx, w, fmt.Errorf("file part 'file' is required"), http.StatusBadRequest)
		return
	}
	defer func() { _ = file.Close() }()

	raw, err := io.ReadAll(io.LimitReader(file, maxAvatarBytes+1))
	if err != nil {
		ps.writeError(ctx, w, fmt.Errorf("could not read file"), http.StatusBadRequest)
		return
	}
	if len(raw) == 0 {
		ps.writeError(ctx, w, fmt.Errorf("empty file"), http.StatusBadRequest)
		return
	}
	if len(raw) > maxAvatarBytes {
		ps.writeError(ctx, w, fmt.Errorf("avatar exceeds %d byte limit — use a smaller image", maxAvatarBytes), http.StatusRequestEntityTooLarge)
		return
	}

	contentType := hdr.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = http.DetectContentType(raw)
	}
	if !allowedAvatarMIME(contentType) {
		ps.writeError(ctx, w, fmt.Errorf("unsupported image type %q (png/jpeg/gif/webp only)", contentType), http.StatusUnsupportedMediaType)
		return
	}
	// Normalize MIME for the data URI (strip params).
	mime := strings.Split(contentType, ";")[0]
	mime = strings.TrimSpace(mime)

	dataURI := "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(raw)

	profileObj, err := ps.profileBusiness.UpdateProfileProperties(
		ctx,
		profileID,
		data.JSONMap{profileDefaultAvaterURI: dataURI},
		false,
	)
	if err != nil {
		log.WithError(err).Error("avatar: update properties failed")
		ps.writeError(ctx, w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": profileObj})
}

func allowedAvatarMIME(ct string) bool {
	ct = strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	switch ct {
	case "image/png", "image/jpeg", "image/jpg", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}
