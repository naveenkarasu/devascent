package content

import (
	"os/exec"
	"testing"

	"devascent/internal/grader"
	"devascent/internal/toolchain"
)

// Known-correct reference solutions for every authored task. The test below
// runs each YAML task's tests against these via the real grader, so any wrong
// expected value or bad YAML transcription fails loudly.
var refSolutions = map[string]string{
	"greet":          "def greet():\n    return \"Hello, world!\"\n",
	"announce":       "def announce():\n    return \"DevAscent online\"\n",
	"double_score":   "def double_score(score):\n    return score * 2\n",
	"total_cost":     "def total_cost(price, quantity):\n    return price * quantity\n",
	"can_enter":      "def can_enter(age):\n    return \"allowed\" if age >= 18 else \"denied\"\n",
	"sign_of":        "def sign_of(n):\n    return \"positive\" if n > 0 else (\"negative\" if n < 0 else \"zero\")\n",
	"count_up":       "def count_up(n):\n    r = 0\n    for i in range(1, n + 1):\n        r += i\n    return r\n",
	"repeat_char":    "def repeat_char(ch, times):\n    return ch * times\n",
	"add":            "def add(a, b):\n    return a + b\n",
	"rectangle_area": "def rectangle_area(width, height):\n    return width * height\n",
	"sum_list":       "def sum_list(numbers):\n    t = 0\n    for n in numbers:\n        t += n\n    return t\n",
	"largest":        "def largest(numbers):\n    m = numbers[0]\n    for n in numbers:\n        if n > m:\n            m = n\n    return m\n",
	"get_price":      "def get_price(menu, item):\n    return menu[item]\n",
	"count_letters":  "def count_letters(word):\n    d = {}\n    for c in word:\n        d[c] = d.get(c, 0) + 1\n    return d\n",
	"is_even":        "def is_even(n):\n    return n % 2 == 0\n",
	"count_passing":  "def count_passing(scores):\n    c = 0\n    for s in scores:\n        if s >= 60:\n            c += 1\n    return c\n",
	"safe_first":     "def safe_first(items):\n    return None if len(items) == 0 else items[0]\n",
	"average":        "def average(numbers):\n    return 0 if len(numbers) == 0 else sum(numbers) / len(numbers)\n",
	"fizz_report":    "def fizz_report(n):\n    out = []\n    for i in range(1, n + 1):\n        if i % 15 == 0:\n            out.append(\"FizzBuzz\")\n        elif i % 3 == 0:\n            out.append(\"Fizz\")\n        elif i % 5 == 0:\n            out.append(\"Buzz\")\n        else:\n            out.append(str(i))\n    return out\n",
	"long_words":     "def long_words(words):\n    return [w for w in words if len(w) >= 4]\n",
	"reverse_text":   "def reverse_text(s):\n    return s[::-1]\n",
	"count_evens":    "def count_evens(nums):\n    return sum(1 for x in nums if x % 2 == 0)\n",
	// coding-floor variants
	"subtract":     "def subtract(a, b):\n    return a - b\n",
	"multiply":     "def multiply(a, b):\n    return a * b\n",
	"max_two":      "def max_two(a, b):\n    return a if a > b else b\n",
	"min_two":      "def min_two(a, b):\n    return a if a < b else b\n",
	"square":       "def square(n):\n    return n * n\n",
	"triple":       "def triple(n):\n    return n * 3\n",
	"last_element": "def last_element(items):\n    return items[-1]\n",
	"first_char":   "def first_char(s):\n    return s[0]\n",
	"negate":       "def negate(n):\n    return -n\n",
	// coding-mid variants
	"count_vowels":    "def count_vowels(s):\n    return sum(1 for c in s.lower() if c in 'aeiou')\n",
	"count_positives": "def count_positives(nums):\n    return sum(1 for n in nums if n > 0)\n",
	// coding-probe variants
	"second_largest": "def second_largest(nums):\n    return sorted(set(nums))[-2]\n",
	"is_palindrome":  "def is_palindrome(s):\n    return s == s[::-1]\n",
	"factorial":      "def factorial(n):\n    r = 1\n    for i in range(2, n + 1):\n        r *= i\n    return r\n",
	"sum_evens":      "def sum_evens(nums):\n    return sum(n for n in nums if n % 2 == 0)\n",
	"count_words":    "def count_words(s):\n    return len(s.split())\n",
}

// TestBenchProblemsGradePass runs every bench problem's reference solution
// through the REAL grader against its tests — the validation gate for the mined
// pool (catches grader/equality mismatches the standalone harness can't).
func TestBenchProblemsGradePass(t *testing.T) {
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python not found on PATH")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Problems) < 100 {
		t.Fatalf("expected 100+ bench problems, got %d", len(c.Problems))
	}
	g := grader.NewNativePython()
	fails := 0
	for _, p := range c.Problems {
		if p.Solution == "" {
			t.Errorf("%s: no reference solution", p.ID)
			continue
		}
		v, err := g.Run("python", p.Solution, p.FuncName, p.Tests, p.GraderShape())
		if err != nil {
			t.Errorf("%s: grader error: %v", p.ID, err)
			continue
		}
		if !v.Passed {
			fails++
			t.Errorf("%s (%s): reference solution did NOT pass: %+v", p.ID, p.Difficulty, v.Results)
		}
	}
	if fails == 0 {
		t.Logf("all %d bench problems graded clean", len(c.Problems))
	}
}

// TestLocalToolchainBackendMatchesGate is the ADR-0007 dual-run gate: every bench
// reference solution must pass through the NEW default backend (LocalToolchain,
// which shells out to the player's installed Python) exactly as it does through
// the NativePython oracle. Proves the default flip didn't regress grading.
// Slow (a python subprocess per problem) → opt-in, skipped in -short.
func TestLocalToolchainBackendMatchesGate(t *testing.T) {
	if testing.Short() {
		t.Skip("local-toolchain dual-run gate is slow; skipped in -short")
	}
	if _, err := exec.LookPath("python"); err != nil {
		if _, err3 := exec.LookPath("python3"); err3 != nil {
			t.Skip("python not found on PATH")
		}
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	oracle := grader.NewNativePython()
	g := grader.NewLocalToolchain(toolchain.New())
	mismatches := 0
	for _, p := range c.Problems {
		if p.Solution == "" {
			continue
		}
		want, _ := oracle.Run("python", p.Solution, p.FuncName, p.Tests, p.GraderShape())
		got, err := g.Run("python", p.Solution, p.FuncName, p.Tests, p.GraderShape())
		if err != nil {
			t.Errorf("%s: LocalToolchain error: %v", p.ID, err)
			continue
		}
		if !got.Passed {
			t.Errorf("%s: reference solution did NOT pass through LocalToolchain: %+v", p.ID, got.Results)
			continue
		}
		if got.Passed != want.Passed {
			mismatches++
			t.Errorf("%s: LocalToolchain (%v) disagrees with NativePython (%v)", p.ID, got.Passed, want.Passed)
		}
	}
	if mismatches == 0 {
		t.Logf("LocalToolchain matched the native oracle across %d bench problems", len(c.Problems))
	}
}

// TestListSlugInvariants enforces honest curated-list counts: exactly 75
// distinct Blind-75 canonical slugs and 150 NeetCode-150 slugs (variants don't
// inflate), and no duplicate problem IDs.
func TestListSlugInvariants(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	blind := map[string]bool{}
	neet := map[string]bool{}
	for _, p := range c.Problems {
		if ids[p.ID] {
			t.Errorf("duplicate problem id: %s", p.ID)
		}
		ids[p.ID] = true
		for _, l := range p.Lists {
			switch l {
			case "blind75":
				blind[p.CanonicalSlug()] = true
			case "neetcode150":
				neet[p.CanonicalSlug()] = true
			}
		}
	}
	if len(blind) != 75 {
		t.Errorf("Blind 75 distinct slugs = %d, want 75", len(blind))
	}
	if len(neet) != 150 {
		t.Errorf("NeetCode 150 distinct slugs = %d, want 150", len(neet))
	}
}

// TestCuratedProblemsHaveEnoughTests: any problem in a curated list (Blind 75 /
// NeetCode 150) must carry at least 3 tests, so the grader can't be satisfied by
// a thin/over-fit case set. Untagged mined "extra practice" is exempt.
func TestCuratedProblemsHaveEnoughTests(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range c.Problems {
		if len(p.Lists) > 0 && len(p.Tests) < 3 {
			t.Errorf("%s (curated: %v) has %d tests, want >= 3", p.ID, p.Lists, len(p.Tests))
		}
	}
}

// TestEveryBenchCategoryHasPrimer: every category that has bench problems must
// have a Learn-panel primer (so [L] always has content to show).
func TestEveryBenchCategoryHasPrimer(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	cats := map[string]bool{}
	for _, p := range c.Problems {
		if p.Category != "" {
			cats[p.Category] = true
		}
	}
	for cat := range cats {
		if _, ok := c.PrimerByCategory(cat); !ok {
			t.Errorf("category %q has bench problems but no primer", cat)
		}
	}
}

// TestEveryPrimerHasSections: every loaded primer must have at least one section
// (the legacy-ops load shim guarantees this even for un-migrated flat files), so
// the section pager always has something to render.
func TestEveryPrimerHasSections(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range c.Primers {
		if len(p.Sections) == 0 {
			t.Errorf("primer %q (%s) has no sections", p.Category, p.Lang)
		}
	}
}

// TestEveryPrimerLangCoversEveryCategory: once a language is offered (it has at
// least one primer), it must cover EVERY bench category — a player who picks that
// session language must never hit "No <lang> primer for this category yet" on the
// Learn panel. This keeps each language's fan-out complete before it ships.
func TestEveryPrimerLangCoversEveryCategory(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	cats := map[string]bool{}
	for _, p := range c.Problems {
		if p.Category != "" {
			cats[p.Category] = true
		}
	}
	for lang := range c.PrimerLangs() {
		for cat := range cats {
			if _, ok := c.PrimerByCategoryAndLang(cat, lang); !ok {
				t.Errorf("language %q is offered but has no primer for category %q", lang, cat)
			}
		}
	}
}

// TestScaffoldBypassesRejected guards the deep-copy / round-trip scaffolds: an
// identity "solution" (return the input unchanged) must FAIL. It reuses each
// problem's real solution + an identity override (Python late-binds, so the
// YAML wrapper calls the override), so this also catches future drift if a
// boundary assertion is ever removed from a scaffold.
func TestScaffoldBypassesRejected(t *testing.T) {
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python not found on PATH")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]Problem{}
	for _, p := range c.Problems {
		byID[p.ID] = p
	}
	overrides := map[string]string{
		"nc-encode-decode-strings": "\ndef encode(strs):\n    return strs\ndef decode(s):\n    return s\n",
		"nc-serialize-tree":        "\ndef serialize(root):\n    return root\ndef deserialize(data):\n    return data\n",
		"nc-clone-graph":           "\ndef clone(node):\n    return node\n",
		"nc-copy-random-list":      "\ndef copy_random_list(head):\n    return head\n",
	}
	g := grader.NewNativePython()
	for id, override := range overrides {
		p, ok := byID[id]
		if !ok {
			t.Errorf("%s: problem not found", id)
			continue
		}
		wrong := p.Solution + override
		v, err := g.Run("python", wrong, p.FuncName, p.Tests, p.GraderShape())
		if err != nil {
			t.Errorf("%s: grader error: %v", id, err)
			continue
		}
		if v.Passed {
			t.Errorf("%s: identity-bypass PASSED — scaffold does not enforce the real algorithm", id)
		}
	}
}

func TestAllAuthoredTasksGradePass(t *testing.T) {
	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python not found on PATH")
	}
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	g := grader.NewNativePython()

	var tasks []*Task
	for i := range c.Diagnostics {
		if c.Diagnostics[i].Task != nil {
			tasks = append(tasks, c.Diagnostics[i].Task)
		}
	}
	for _, l := range c.Lessons {
		// Only the Python lessons are graded here against the Python refSolutions
		// map; per-language lesson variants (Go/Rust/Java/C#/JS/TS) are validated
		// through their own toolchains by TestLessonsGradeNativeRoundTrip.
		if l.Lang != "" && l.Lang != "python" {
			continue
		}
		for _, s := range l.Stages {
			if s.Task != nil {
				tasks = append(tasks, s.Task)
			}
		}
	}
	if len(tasks) < 20 {
		t.Fatalf("expected 20+ authored tasks, got %d", len(tasks))
	}

	for _, task := range tasks {
		sol, ok := refSolutions[task.FuncName]
		if !ok {
			t.Errorf("%s: no reference solution registered", task.FuncName)
			continue
		}
		v, err := g.Run("python", sol, task.FuncName, task.Tests, grader.Shape{})
		if err != nil {
			t.Errorf("%s: grader error: %v", task.FuncName, err)
			continue
		}
		if !v.Passed {
			t.Errorf("%s: reference solution did NOT pass authored tests: %+v", task.FuncName, v.Results)
		}
	}
}
