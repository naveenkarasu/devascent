// Package content loads DevAscent's data-driven catalog (lessons, diagnostic
// items, …) from embedded YAML. Mirrors the Unified Content Catalog design
// (ContentItem → Lesson / DiagnosticItem) so design and code stay in sync.
package content

import (
	"strings"

	"devascent/internal/grader"
)

// Slugify normalizes a title to a canonical slug (lowercase, non-alphanumeric
// runs → "-", trimmed). Used to group problem variants under one canonical slug.
func Slugify(s string) string {
	var b strings.Builder
	prevHyphen := false
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevHyphen = false
		} else if !prevHyphen {
			b.WriteByte('-')
			prevHyphen = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// Task is a code exercise: the player writes funcName; we grade by calling it.
// Solution is an optional reference implementation (in the lesson's language),
// never shown to the player — it exists so the native lesson-grading test can
// prove the task actually grades green in that language (mirrors Problem.Solution).
type Task struct {
	Prompt   string            `yaml:"prompt"`
	FuncName string            `yaml:"func_name"`
	Starter  string            `yaml:"starter"`
	Solution string            `yaml:"solution,omitempty"`
	Tests    []grader.TestCase `yaml:"tests"`
}

// LessonStage is one Gradual-Release stage: i_do (read) / we_do / you_do.
type LessonStage struct {
	Kind  string `yaml:"kind"`
	Title string `yaml:"title"`
	Body  string `yaml:"body"`
	Task  *Task  `yaml:"task"` // nil for i_do (read-only) stages
}

// Lesson is one Tutorial Island lesson. Lang ties a variant to the session
// language (empty = "python", the original teaching language). The engine plays
// ONE variant per lesson ID — the one matching m.lang, else the Python fallback —
// so a Go player reads Go syntax (func/:=) instead of Python (def). See
// LessonsForLang. Variants of one lesson share ID, Order, and Title; for most
// lessons the graded Task (FuncName/Tests) is identical too and only the Body +
// starter differ. The exceptions are lessons "dicts" and "read-the-crash": Python
// stays dynamic (map-as-argument; None for empty), but every statically-typed
// language needs harness-safe tasks — the in-language graders can't build a map
// argument or return None — so those variants deliberately diverge (a map used
// internally returning an int; a -1 sentinel for empty). The non-Python variants
// share one task shape across all languages, validated by the native round-trip.
type Lesson struct {
	ID     string        `yaml:"id"`
	Lang   string        `yaml:"lang"` // python | go | java | csharp | javascript | typescript | rust | cpp
	Order  int           `yaml:"order"`
	Title  string        `yaml:"title"`
	Stages []LessonStage `yaml:"stages"`
}

// Choice is one option of a multiple-choice diagnostic item.
type Choice struct {
	Text     string `yaml:"text"`
	Correct  bool   `yaml:"correct"`
	Feedback string `yaml:"feedback"`
	Value    string `yaml:"value"` // for the self-report intro: "never" | "a-little" | "regularly"
}

// Spec is a free-text "read the problem statement" item, graded deterministically:
// the player's answer passes if it contains every Required keyword (case-insensitive).
type Spec struct {
	Required []string `yaml:"required"` // all must appear in the answer
	Answer   string   `yaml:"answer"`   // model answer, shown after grading
}

// Diagnostic is one intake item. Measures: coding | machine | spec | self.
// "self" (self-report) is informational only — it never affects routing.
// Slot groups interchangeable variants: the engine plays ONE variant per slot
// (chosen at runtime) so the ladder differs run-to-run. All variants in a slot
// MUST share Kind and Measures (enforced in Load).
type Diagnostic struct {
	ID         string   `yaml:"id"`
	Slot       string   `yaml:"slot"`       // e.g. "coding-floor", "machine-terminal", "spec"
	Lang       string   `yaml:"lang"`       // "" (=python) | go | java | csharp | javascript | typescript | rust | cpp
	Order      int      `yaml:"order"`      // slot position in the ladder (shared by a slot's variants)
	Difficulty int      `yaml:"difficulty"` // 1..3 (tagged now; selection ignores it in v1)
	Measures   string   `yaml:"measures"`
	Kind       string   `yaml:"kind"` // "code" | "choice" | "spec"
	Prompt     string   `yaml:"prompt"`
	Task       *Task    `yaml:"task"`    // for kind: code
	Choices    []Choice `yaml:"choices"` // for kind: choice
	Spec       *Spec    `yaml:"spec"`    // for kind: spec
}

// DevTask is one Dev-Literacy micro-task: the player types a terminal command,
// graded by base-command + required-flags (a checker, not a real shell). The
// engine picks a few tasks across distinct Categories so each run varies.
type DevTask struct {
	ID         string   `yaml:"id"`
	Order      int      `yaml:"order"`
	Category   string   `yaml:"category"`   // navigation | files | inspect | git | text | archive | ...
	Difficulty int      `yaml:"difficulty"` // 1..3
	Title      string   `yaml:"title"`
	Prompt     string   `yaml:"prompt"`
	Commands   []string `yaml:"commands"` // acceptable base commands (first token), e.g. ["ls","dir"]
	Flags      []string `yaml:"flags"`    // tokens that must ALL appear (e.g. ["-r"]); usually empty
	Accept     []string `yaml:"accept"`   // explicit full-line accepted forms (exact, fallback)
	Hint       string   `yaml:"hint"`
	Success    string   `yaml:"success"`
}

// Problem is one Step 0 bench problem (from the mined DSA pool). Graded by the
// same function-call loop as everything else. Solution travels with it for the
// grade-validation test (and a future "reveal" feature); the TUI never shows it.
type Problem struct {
	ID         string            `yaml:"id"`
	Title      string            `yaml:"title"`
	Difficulty string            `yaml:"difficulty"` // easy | medium | hard
	Pattern    string            `yaml:"pattern"`
	Category   string            `yaml:"category"` // coarse grouping for the bench browser
	Lists      []string          `yaml:"lists"`    // curated-list membership (e.g. ["blind75"]); hand-verified only
	FuncName   string            `yaml:"func_name"`
	Prompt     string            `yaml:"prompt"`
	Starter    string            `yaml:"starter"`
	Solution   string            `yaml:"solution"` // reference; validated by tests, NOT shown to the player
	Tests      []grader.TestCase `yaml:"tests"`

	// Data-structure I/O (node problems). Empty = plain JSON I/O.
	Shape    string   `yaml:"shape"`     // "" | "linkedlist" | "tree"
	ArgKinds []string `yaml:"arg_kinds"` // per-arg: "node" else raw
	RetKind  string   `yaml:"ret_kind"`  // "node" else raw

	// Slug groups VARIANTS of the same canonical problem (e.g. a mined version and
	// the clean version). Curated-list counts are by distinct slug; serving a list
	// picks one random variant per slug. Defaults to Slugify(Title).
	Slug string `yaml:"slug"`
}

// CanonicalSlug returns the explicit Slug or, if empty, Slugify(Title).
func (p Problem) CanonicalSlug() string {
	if p.Slug != "" {
		return p.Slug
	}
	return Slugify(p.Title)
}

// GraderShape returns the grader.Shape for this problem (zero for plain I/O).
func (p Problem) GraderShape() grader.Shape {
	return grader.Shape{Kind: p.Shape, ArgKinds: p.ArgKinds, RetKind: p.RetKind}
}

// PrimerOp is one basic operation shown in a pattern primer (label + snippet).
type PrimerOp struct {
	Label string `yaml:"label"`
	Code  string `yaml:"code"`
}

// PrimerSection groups related ops under a heading (e.g. Declarations, Basic
// operations, Built-in functions, Conversions, Iteration & loops, Special
// operations). Sections are also the PAGING unit in the Learn panel — one
// section per screen — so a rich primer stays readable instead of scrolling off
// the top of the terminal.
type PrimerSection struct {
	Title string     `yaml:"title"`
	Ops   []PrimerOp `yaml:"ops"`
}

// Primer is a per-category, per-language coding refresher (the "Learn" panel):
// basic ops a player can reuse + one worked example with reasoning. Not graded —
// pure reference. Lang ties the primer to the session language (m.lang); an empty
// Lang is normalized to "python" at load time (see Load) for back-compat with the
// original Python-only primers.
type Primer struct {
	Category string          `yaml:"category"`
	Lang     string          `yaml:"lang"` // python | java | cpp | csharp | javascript | typescript | rust | go
	Title    string          `yaml:"title"`
	Summary  string          `yaml:"summary"`
	Sections []PrimerSection `yaml:"sections"` // grouped, paged ops (preferred)
	Ops      []PrimerOp      `yaml:"ops"`      // legacy flat list; load wraps it into one section
	Example  string          `yaml:"example"`
}

// Exercise is one Stage-2 Advanced-Topics task. It is REFERENCE-first (the player
// sees BrokenCode, attempts, then reveals Bug + FixedCode), but FORWARD-COMPATIBLE
// for grading: Check + Signal say HOW a grader validates it. Today only Check=="tests"
// runs (Python, via the existing function-call harness — FuncName/Tests/FixedCode);
// "compile-error"/"compiles"/"stdout" are validated once a per-language toolchain
// grader exists, with no content re-authoring. "none" = reveal-only.
type Exercise struct {
	Prompt     string `yaml:"prompt"`
	Kind       string `yaml:"kind"`        // fix-it | spot-the-bug | predict-output
	BrokenCode string `yaml:"broken_code"` // shown to the player (highlighted)
	FixedCode  string `yaml:"fixed_code"`  // model solution; revealed, and graded when Check=="tests"
	Bug        string `yaml:"bug"`         // what's wrong + why (the spot-the-bug answer)
	Check      string `yaml:"check"`       // tests | compile-error | stdout | compiles | none
	Signal     string `yaml:"signal"`      // value Check compares (error code / expected stdout); "" for tests

	// Check=="tests": grade FixedCode (the solution) through the function-call harness.
	FuncName string            `yaml:"func_name,omitempty"`
	Tests    []grader.TestCase `yaml:"tests,omitempty"`
	Shape    string            `yaml:"shape,omitempty"`
	ArgKinds []string          `yaml:"arg_kinds,omitempty"`
	RetKind  string            `yaml:"ret_kind,omitempty"`
}

// GraderShape returns the grader.Shape for a Check=="tests" exercise.
func (e Exercise) GraderShape() grader.Shape {
	return grader.Shape{Kind: e.Shape, ArgKinds: e.ArgKinds, RetKind: e.RetKind}
}

// AdvancedTopic is one Stage-2 language-specific topic (e.g. "Ownership & Borrowing"
// for Rust). The explainer reuses the primer Sections renderer; Exercises carry the
// problems. Reached from the bench "Advanced Topics" browse, scoped to m.lang.
type AdvancedTopic struct {
	Lang      string          `yaml:"lang"`  // python | java | cpp | ...
	Group     string          `yaml:"group"` // taxonomy category, e.g. "Ownership & Borrowing"
	Title     string          `yaml:"title"`
	Tag       string          `yaml:"tag"`     // E | C | P | gotcha
	Summary   string          `yaml:"summary"` // prose explainer (the "what & why")
	Sections  []PrimerSection `yaml:"sections"`
	Exercises []Exercise      `yaml:"exercises"`
}

// Catalog is the loaded content set, ordered.
type Catalog struct {
	Intro          []Diagnostic // warm-up items shown before the counted intake (not scored)
	Diagnostics    []Diagnostic
	Lessons        []Lesson
	DevTasks       []DevTask
	Problems       []Problem       // Step 0 bench pool
	Primers        []Primer        // per-category coding refreshers (Learn panel)
	AdvancedTopics []AdvancedTopic // Stage-2 language-specific topics
	InstallGuides  []InstallGuide  // per-OS toolchain install instructions (ADR-0007)
}

// InstallGuide holds per-OS instructions to install one language's toolchain
// (ADR-0007: DevAscent bundles no runtimes — the player installs their own). The
// SAME data renders both in-game (the Install Help screen, shown when a player
// lands on an unavailable language) and to the repo INSTALL.md, so they never
// drift. OS keys are "windows" | "macos" | "linux".
type InstallGuide struct {
	Lang  string                  `yaml:"lang"`
	Label string                  `yaml:"label"`
	Notes string                  `yaml:"notes,omitempty"`
	OS    map[string]InstallSteps `yaml:"os"`
}

// InstallSteps is one OS's install path: a download link, ordered steps, and the
// command to verify success.
type InstallSteps struct {
	Link   string   `yaml:"link"`
	Steps  []string `yaml:"steps"`
	Verify string   `yaml:"verify,omitempty"`
}

// InstallGuideForLang returns the install guide for a language, or false.
func (c Catalog) InstallGuideForLang(lang string) (InstallGuide, bool) {
	for _, g := range c.InstallGuides {
		if g.Lang == lang {
			return g, true
		}
	}
	return InstallGuide{}, false
}

// AdvancedTopicsByLang returns the advanced topics for a language, in load order.
func (c Catalog) AdvancedTopicsByLang(lang string) []AdvancedTopic {
	var out []AdvancedTopic
	for _, t := range c.AdvancedTopics {
		if t.Lang == lang {
			out = append(out, t)
		}
	}
	return out
}

// PrimerByCategoryAndLang returns the primer for a (category, language) pair, or
// false. This is the lang-aware lookup the Learn panel uses; the session language
// (m.lang) decides which language's primer renders.
func (c Catalog) PrimerByCategoryAndLang(cat, lang string) (Primer, bool) {
	if lang == "" {
		lang = "python"
	}
	for _, p := range c.Primers {
		if p.Category == cat && p.Lang == lang {
			return p, true
		}
	}
	return Primer{}, false
}

// PrimerByCategory returns the Python primer for a category, or false. Kept as a
// convenience for the default (playable) language; lang-aware callers should use
// PrimerByCategoryAndLang.
func (c Catalog) PrimerByCategory(cat string) (Primer, bool) {
	return c.PrimerByCategoryAndLang(cat, "python")
}

// normLang treats an empty language as "python" (the original default), so the
// per-language accessors can assume a concrete value.
func normLang(lang string) string {
	if lang == "" {
		return "python"
	}
	return lang
}

// LessonsForLang returns one variant per lesson (by ID), in load (Order) order:
// the variant whose Lang matches lang, else the Python fallback. This is how
// Tutorial Island renders in the session language — a Go beginner sees Go syntax,
// a Rust beginner sees Rust, and an unauthored language degrades to the Python
// lesson rather than vanishing. Variants of one lesson share ID and Order; only
// the explanatory Body and starter hints differ (the graded Task is identical).
func (c Catalog) LessonsForLang(lang string) []Lesson {
	lang = normLang(lang)
	byID := map[string][]Lesson{}
	var order []string
	for _, l := range c.Lessons {
		if _, ok := byID[l.ID]; !ok {
			order = append(order, l.ID)
		}
		byID[l.ID] = append(byID[l.ID], l)
	}
	out := make([]Lesson, 0, len(order))
	for _, id := range order {
		variants := byID[id]
		var pick *Lesson
		for i := range variants {
			if normLang(variants[i].Lang) == lang {
				pick = &variants[i]
				break
			}
		}
		if pick == nil { // no language-specific variant → Python fallback
			for i := range variants {
				if normLang(variants[i].Lang) == "python" {
					pick = &variants[i]
					break
				}
			}
		}
		if pick != nil {
			out = append(out, *pick)
		}
	}
	return out
}

// DiagnosticsForLang returns the intake pool narrowed to the session language:
// for each slot, the variants whose Lang matches lang if any exist, else the
// Python (default) variants. Most slots are language-neutral (code tasks whose
// starter is regenerated per language, shell/spec prose) and carry only Python
// variants, so they fall through unchanged; a language-specific slot like
// machine-error supplies native variants (a Go compile error, not a Python
// IndentationError). The returned pool feeds selectIntake unchanged.
func (c Catalog) DiagnosticsForLang(lang string) []Diagnostic {
	lang = normLang(lang)
	bySlot := map[string][]Diagnostic{}
	var order []string
	for _, d := range c.Diagnostics {
		if _, ok := bySlot[d.Slot]; !ok {
			order = append(order, d.Slot)
		}
		bySlot[d.Slot] = append(bySlot[d.Slot], d)
	}
	var out []Diagnostic
	for _, slot := range order {
		var langMatch, pyFallback []Diagnostic
		for _, d := range bySlot[slot] {
			if normLang(d.Lang) == lang {
				langMatch = append(langMatch, d)
			}
			if normLang(d.Lang) == "python" {
				pyFallback = append(pyFallback, d)
			}
		}
		if len(langMatch) > 0 {
			out = append(out, langMatch...)
		} else {
			out = append(out, pyFallback...)
		}
	}
	return out
}

// PrimerLangs returns the set of languages that have at least one authored
// primer. The language-pick screen offers only these (so a player can't pick a
// language with no reference content yet).
func (c Catalog) PrimerLangs() map[string]bool {
	out := map[string]bool{}
	for _, p := range c.Primers {
		lang := p.Lang
		if lang == "" {
			lang = "python"
		}
		out[lang] = true
	}
	return out
}
