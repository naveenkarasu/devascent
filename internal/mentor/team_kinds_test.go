package mentor

import (
	"strings"
	"testing"
)

// The standup/discuss kinds enforce the same output contract: no code, bounded length.
func TestValidate_StandupAndDiscuss(t *testing.T) {
	if err := Validate(KindStandup, "", "You: continuing PXF-1. Maya: on PXF-2, no blockers."); err != nil {
		t.Errorf("a plain standup should pass: %v", err)
	}
	if err := Validate(KindStandup, "", "```go\nfmt.Println()\n```"); err == nil {
		t.Error("a standup with a code block should be rejected")
	}
	if err := Validate(KindDiscuss, "", "I'll reproduce it with a test, then make the smallest fix. ~2 days."); err != nil {
		t.Errorf("a plain plan should pass: %v", err)
	}
	if err := Validate(KindDiscuss, "", "```\ncode\n```"); err == nil {
		t.Error("a plan with a code block should be rejected")
	}
	if err := Validate(KindDiscuss, "", strings.Repeat("word ", 200)); err == nil {
		t.Error("an over-long plan should be rejected")
	}
}

// Offline templates cover the new kinds (the standup echoes the status facts; the
// discuss plan reads like a plan).
func TestTemplateAnswer_TeamKinds(t *testing.T) {
	got := templateAnswer(Request{Kind: KindStandup, Status: []string{"You: x", "Maya: y"}})
	if !strings.Contains(got, "Maya: y") {
		t.Errorf("the standup template should include the status lines, got %q", got)
	}
	if templateAnswer(Request{Kind: KindStandup}) == "" {
		t.Error("an empty standup should still produce a fallback line")
	}
	if !strings.Contains(strings.ToLower(templateAnswer(Request{Kind: KindDiscuss})), "plan") {
		t.Error("the discuss template should read like a plan")
	}
}

// The discuss prompt shares the public Definition-of-Ready (title, acceptance) and
// the persona — but Request structurally cannot carry hidden tests or solutions.
func TestBuildPrompt_DiscussSharesDoR(t *testing.T) {
	p := BuildPrompt(Request{
		Kind:       KindDiscuss,
		Persona:    "Maya, a Junior Engineer",
		Title:      "Fix the paginator",
		Prompt:     "Off by one.",
		Priority:   "Major",
		Acceptance: []string{"page 1 returns items 1..size"},
	})
	for _, want := range []string{"Maya, a Junior Engineer", "Fix the paginator", "page 1 returns items 1..size"} {
		if !strings.Contains(p, want) {
			t.Errorf("the discuss prompt is missing %q", want)
		}
	}
}

// The standup prompt is built only from the status facts it is handed.
func TestBuildPrompt_StandupFromFacts(t *testing.T) {
	p := BuildPrompt(Request{Kind: KindStandup, Persona: "Sam", Day: 3, Status: []string{"Maya: on PXF-7"}})
	if !strings.Contains(p, "Maya: on PXF-7") {
		t.Errorf("the standup prompt should include the status facts, got:\n%s", p)
	}
}
