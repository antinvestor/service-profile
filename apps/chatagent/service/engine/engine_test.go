package engine_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/antinvestor/service-profile/apps/chatagent/service/engine"
)

func placementLikeContext() engine.ContextDef {
	return engine.ContextDef{
		Key:     "stawi.placement.intake",
		Purpose: "Collect placement preferences and qualifications for opportunity matching.",
		Fields: []engine.FieldDef{
			{Name: "target_job_title", Type: engine.FieldString, Required: true, Priority: 1, Ask: "What role are you targeting?", Why: "drives matching"},
			{Name: "capabilities", Type: engine.FieldString, Required: true, Priority: 2, MinLength: 80, Ask: "Please share your CV or work history.", Why: "qualifications", EvidenceHints: []string{"document"}},
			{Name: "job_types", Type: engine.FieldStringList, Required: true, Priority: 3, Enum: []string{"Full-time", "Part-time", "Contract"}, Ask: "Which job types?"},
			{Name: "salary_min", Type: engine.FieldNumber, Required: true, Priority: 4, Ask: "What salary minimum?"},
			{Name: "preferred_countries", Type: engine.FieldStringList, Required: true, Priority: 5, Ask: "Which countries to search?"},
			{Name: "experience_level", Type: engine.FieldString, Required: true, Priority: 6, Enum: []string{"entry", "junior", "mid", "senior", "lead", "executive"}, Ask: "Experience level?"},
			{Name: "linkedin", Type: engine.FieldString, Required: false, Priority: 99, Ask: "LinkedIn (optional)"},
		},
		ReplyPolicy: engine.ReplyPolicy{MaxSentences: 3, AskOneMissingOnly: true},
	}
}

func TestAssess_ReadyWhenRequiredPresent(t *testing.T) {
	t.Parallel()
	def := placementLikeContext()
	fields := engine.Fields{
		"target_job_title":    "Backend Engineer",
		"capabilities":        strings.Repeat("Experience building APIs. Education skills employment. ", 3),
		"job_types":           []string{"Full-time"},
		"salary_min":          50000.0,
		"preferred_countries": []string{"KE"},
		"experience_level":    "senior",
	}
	_, missing, ready := engine.Assess(def, fields)
	require.True(t, ready)
	require.Empty(t, missing)
}

func TestAssess_MissingPriorityOrder(t *testing.T) {
	t.Parallel()
	def := placementLikeContext()
	_, missing, ready := engine.Assess(def, engine.Fields{})
	require.False(t, ready)
	require.Equal(t, "target_job_title", missing[0])
	require.Equal(t, "capabilities", missing[1])
}

func TestApplyDocuments_FillsCapabilitiesFromCV(t *testing.T) {
	t.Parallel()
	def := placementLikeContext()
	cv := "John Doe\nExperience: 5 years backend\nEducation: BSc CS\nSkills: Go, PostgreSQL\n" + strings.Repeat("worked at companies with responsibilities. ", 5)
	fields := engine.ApplyDocuments(def, engine.Fields{}, []engine.Document{
		{Name: "cv", Kind: "cv", Text: cv},
	})
	caps, _ := fields["capabilities"].(string)
	require.Contains(t, caps, "Experience")
	// also stores under kind/name keys
	require.NotEmpty(t, fields["cv"])
}

// Prefer Evaluate path; merge is covered through Turn.
func TestTurn_EvidenceFirstDoesNotAskForKnownCV(t *testing.T) {
	t.Parallel()
	def := placementLikeContext()
	cv := "Jane Doe\nExperience building distributed systems for 6 years.\nEducation: MSc\nSkills: Go, Kubernetes\nEmployment history at Acme Corp with responsibilities leading teams. "
	agent := engine.NewAgent(nil)
	res, err := agent.Evaluate(context.Background(), def, nil, engine.Evidence{
		SeedFields: engine.Fields{
			"target_job_title":    "Platform Engineer",
			"job_types":           []string{"Full-time"},
			"salary_min":          80000.0,
			"preferred_countries": []string{"KE", "UG"},
			"experience_level":    "senior",
		},
		Documents: []engine.Document{{Name: "cv", Kind: "cv", Text: cv}},
	})
	require.NoError(t, err)
	require.True(t, res.Ready, "missing=%v fields=%v", res.Missing, res.Fields)
	require.Empty(t, res.Missing)
	require.NotContains(t, strings.ToLower(res.Reply), "cv")
}

func TestTurn_AsksOnlyNextMissing(t *testing.T) {
	t.Parallel()
	def := placementLikeContext()
	agent := engine.NewAgent(nil)
	res, err := agent.Turn(context.Background(), def, engine.TurnInput{
		Message: "I want full-time senior backend roles in Kenya",
		Evidence: engine.Evidence{
			SeedFields: engine.Fields{
				"target_job_title": "Backend Engineer",
				"job_types":        []string{"Full-time"},
				"experience_level": "senior",
			},
		},
		PriorFields: nil,
	})
	require.NoError(t, err)
	require.False(t, res.Ready)
	require.Equal(t, "capabilities", res.Missing[0])
	require.Contains(t, strings.ToLower(res.Reply), "cv")
}

type stubLLM struct {
	reply string
}

func (s stubLLM) Complete(_ context.Context, _ string) (string, error) {
	return s.reply, nil
}

func TestTurn_LLMMergeAndNoFalseReady(t *testing.T) {
	t.Parallel()
	def := placementLikeContext()
	payload, _ := json.Marshal(map[string]any{
		"fields": map[string]any{
			"target_job_title": "Data Engineer",
			"salary_min":       100000,
		},
		"reply": "You're all set — everything I need is here!",
	})
	agent := engine.NewAgent(stubLLM{reply: string(payload)})
	res, err := agent.Turn(context.Background(), def, engine.TurnInput{
		Message: "I'm a data engineer, min salary 100k",
		Evidence: engine.Evidence{
			Messages: []engine.Message{{Role: "user", Content: "hello"}},
		},
	})
	require.NoError(t, err)
	require.False(t, res.Ready)
	require.Equal(t, "Data Engineer", res.Fields["target_job_title"])
	require.NotContains(t, strings.ToLower(res.Reply), "all set")
}

func TestMergeFields_DoesNotWipeWithEmpty(t *testing.T) {
	t.Parallel()
	base := engine.Fields{"target_job_title": "SRE", "salary_min": 50.0}
	out := engine.MergeFields(base, engine.Fields{"target_job_title": "", "salary_min": nil, "job_types": []string{"Full-time"}})
	require.Equal(t, "SRE", out["target_job_title"])
	require.InDelta(t, 50.0, out["salary_min"], 0.001)
	require.Equal(t, []string{"Full-time"}, out["job_types"])
}

func TestSanitize_EnumEnforcement(t *testing.T) {
	t.Parallel()
	def := placementLikeContext()
	out := engine.Sanitize(def, engine.Fields{
		"experience_level": "SENIOR",
		"job_types":        []any{"full-time", "gig"},
	})
	// enum match is case-insensitive but stores canonical enum value
	require.Equal(t, "senior", out["experience_level"])
	// job_types enum has "Full-time" - full-time should match
	require.Equal(t, []string{"Full-time"}, out["job_types"])
}
