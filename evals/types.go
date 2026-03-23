package main

// EvalSuite is the top-level structure of evals.json.
type EvalSuite struct {
	Version int          `json:"version"`
	Skills  []SkillEvals `json:"skills"`
}

// SkillEvals groups evals for a single skill.
type SkillEvals struct {
	Name  string `json:"skill_name"`
	Path  string `json:"skill_path"`
	Evals []Eval `json:"evals"`
}

// Eval is a single test case.
type Eval struct {
	ID         int         `json:"id"`
	Prompt     string      `json:"prompt"`
	Expected   string      `json:"expected_output,omitempty"`
	Assertions []Assertion `json:"assertions,omitempty"`
}

// Assertion is a machine-checkable condition on the eval output.
type Assertion struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Text  string `json:"text,omitempty"`
}

// GradingResult holds the outcome of grading a single assertion.
type GradingResult struct {
	Text     string `json:"text"`
	Passed   bool   `json:"passed"`
	Evidence string `json:"evidence"`
}

// GradingGroup holds assertion stats for a group of runs.
type GradingGroup struct {
	TotalAssertions int     `json:"total_assertions"`
	Passed          int     `json:"passed"`
	Failed          int     `json:"failed"`
	PassRate        float64 `json:"pass_rate"`
}

// GradingSummary aggregates grading results across a workspace.
type GradingSummary struct {
	WithSkill    GradingGroup `json:"with_skill"`
	WithoutSkill GradingGroup `json:"without_skill"`
	Timestamp    string       `json:"timestamp"`
}
