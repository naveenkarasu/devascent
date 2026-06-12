package mentor

import (
	"context"
	"strings"
	"testing"
)

// fakeBackend scripts one reply (or error) for Service tests.
type fakeBackend struct {
	reply string
	err   error
}

func (f *fakeBackend) ID() string              { return "fake" }
func (f *fakeBackend) Name() string            { return "Fake" }
func (f *fakeBackend) Present() (bool, string) { return true, "test" }
func (f *fakeBackend) Ask(context.Context, string) (string, error) {
	return f.reply, f.err
}

func serviceWith(f *fakeBackend) *Service {
	s := NewService(Config{Backend: "fake"})
	s.backends["fake"] = f
	s.order = append(s.order, "fake")
	return s
}

func strategyReq() Request {
	return Request{Kind: KindStrategy, Lang: "python", Title: "Two Sum",
		Prompt: "find indices", Category: "Arrays & Hashing", Difficulty: "easy"}
}

func TestHint_TemplatesWhenNoBackend(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	s := NewService(Config{})
	r := s.Hint(context.Background(), strategyReq())
	if r.Source != "template" || r.FellBack {
		t.Fatalf("no-backend hint: %+v", r)
	}
	if !strings.Contains(r.Text, "hash map") {
		t.Fatalf("arrays strategy template missing: %q", r.Text)
	}
}

func TestHint_AIAnswerPassesValidation(t *testing.T) {
	s := serviceWith(&fakeBackend{reply: "Think about a single pass with a hash map that remembers complements."})
	r := s.Hint(context.Background(), strategyReq())
	if r.Source != "fake" || r.FellBack {
		t.Fatalf("clean AI answer rejected: %+v", r)
	}
}

func TestHint_GuardrailViolationFallsBack(t *testing.T) {
	bad := "Here you go:\n```python\ndef two_sum(nums, target):\n    seen = {}\n```"
	s := serviceWith(&fakeBackend{reply: bad})
	r := s.Hint(context.Background(), strategyReq())
	if r.Source != "template" || !r.FellBack {
		t.Fatalf("code-block strategy answer not rejected: %+v", r)
	}
}

func TestHint_BackendErrorFallsBack(t *testing.T) {
	s := serviceWith(&fakeBackend{err: context.DeadlineExceeded})
	r := s.Hint(context.Background(), strategyReq())
	if r.Source != "template" || !r.FellBack {
		t.Fatalf("backend error not handled: %+v", r)
	}
}

func TestHint_ScrubsPathsFromAIOutput(t *testing.T) {
	s := serviceWith(&fakeBackend{reply: `Look at the loop you wrote in C:\Users\someone\code\main.py carefully.`})
	r := s.Hint(context.Background(), strategyReq())
	if strings.Contains(r.Text, `C:\Users`) {
		t.Fatalf("machine path leaked: %q", r.Text)
	}
}

func TestValidate_PerKindContracts(t *testing.T) {
	cases := []struct {
		kind Kind
		lang string
		text string
		ok   bool
	}{
		{KindStrategy, "go", "Use a map.", true},
		{KindStrategy, "go", "```go\ncode\n```", false},
		{KindStrategy, "go", "x := 1;\ny := 2;\nz := x + y;\nreturn z;", false},
		{KindWalkthrough, "python", "1. Sort.\n2. Sweep with pseudocode:\n```\nfor item in list: keep max\n```", true},
		{KindWalkthrough, "python", "```python\ndef solve():\n    pass\n```", false},
		{KindFollowup, "go", "Why does the window never skip a valid answer?", true},
		{KindFollowup, "go", "Nice work, looks good.", false},
		{KindReview, "go", "Strength: clear naming. Improvement: hoist the invariant check.", true},
		{KindStrategy, "go", "", false},
	}
	for i, c := range cases {
		err := Validate(c.kind, c.lang, c.text)
		if (err == nil) != c.ok {
			t.Errorf("case %d (%s): err=%v, want ok=%v", i, c.kind, err, c.ok)
		}
	}
}

func TestBuildPrompt_GuardrailScope(t *testing.T) {
	req := strategyReq()
	req.PlayerCode = `# saved from D:\private\secret\notes.py
def two_sum(nums, target): pass`
	req.FailedRuns = 3
	req.FirstFail = "case-2"
	p := BuildPrompt(req)
	if strings.Contains(p, `D:\private`) {
		t.Fatal("machine path leaked into the context pack")
	}
	if !strings.Contains(p, "def two_sum") || !strings.Contains(p, "case-2") {
		t.Fatal("allowed context missing from the pack")
	}
	if !strings.Contains(p, "no code") {
		t.Fatal("strategy hard rules missing")
	}
}

func TestNudges_AllBenchCategoriesCovered(t *testing.T) {
	cats := []string{
		"Arrays & Hashing", "Two Pointers & Sliding Window", "Stack", "Binary Search",
		"Linked List", "Trees & Graphs", "Dynamic Programming", "Backtracking",
		"Heap / Priority Queue", "Intervals", "Greedy", "Math & Bit", "Strings",
		"Tries", "Advanced Graphs",
	}
	for _, c := range cats {
		if _, ok := nudges[c]; !ok {
			t.Errorf("category %q has no authored nudges", c)
		}
		if _, ok := strategyTemplates[c]; !ok {
			t.Errorf("category %q has no strategy template", c)
		}
		for attempt := 0; attempt < 4; attempt++ {
			if Nudge(c, attempt) == "" {
				t.Errorf("%s attempt %d: empty nudge", c, attempt)
			}
		}
	}
	if Nudge("Unknown Category", 1) == "" {
		t.Error("generic nudge fallback empty")
	}
}

func TestConfigRoundTrip(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	want := Config{Backend: "ollama", Endpoint: "http://localhost:1234/v1", Model: "qwen3:4b"}
	if err := SaveConfig(want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := LoadConfig()
	if err != nil || got != want {
		t.Fatalf("round trip: got %+v err %v", got, err)
	}
}

func TestSelect_TemplateAlwaysWorks(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	s := NewService(Config{Backend: "claude"})
	if err := s.Select(context.Background(), "template"); err != nil {
		t.Fatalf("selecting templates failed: %v", err)
	}
	if got, _ := LoadConfig(); got.Backend != "" {
		t.Fatalf("template selection persisted as %q", got.Backend)
	}
}

func TestProbe_SentinelChecked(t *testing.T) {
	s := serviceWith(&fakeBackend{reply: "sure, whatever"})
	if err := s.Probe(context.Background(), "fake"); err == nil {
		t.Fatal("probe accepted a reply without the sentinel")
	}
	s2 := serviceWith(&fakeBackend{reply: "DEVASCENT_OK"})
	if err := s2.Probe(context.Background(), "fake"); err != nil {
		t.Fatalf("probe rejected the sentinel: %v", err)
	}
	var st Status
	for _, row := range s2.Statuses() {
		if row.ID == "fake" {
			st = row
		}
	}
	if !st.Probed || !st.ProbeOK {
		t.Fatalf("probe result not cached in statuses: %+v", st)
	}
}
