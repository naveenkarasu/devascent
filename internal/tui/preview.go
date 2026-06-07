package tui

import (
	"fmt"
	"sort"
	"strings"

	"devascent/internal/content"
)

// previewCategoryOrder is the canonical category order for the primer preview
// (matches the bench pattern progression).
var previewCategoryOrder = []string{
	"Arrays & Hashing", "Two Pointers & Sliding Window", "Binary Search", "Stack",
	"Strings", "Math & Bit", "Greedy", "Dynamic Programming", "Backtracking",
	"Heap / Priority Queue", "Trees & Graphs", "Intervals", "Linked List", "Tries",
	"Advanced Graphs",
}

// PreviewPrimers renders every primer for a language — all sections + the worked
// example, syntax-highlighted exactly as the in-game Learn panel shows them — to
// one string. It's a QA aid: pipe it to your terminal to eyeball a whole language
// at once without walking the game. Errors if the language has no primers.
func PreviewPrimers(lang string) (string, error) {
	cat, err := content.Load()
	if err != nil {
		return "", err
	}
	if !cat.PrimerLangs()[lang] {
		have := make([]string, 0)
		for l := range cat.PrimerLangs() {
			have = append(have, l)
		}
		sort.Strings(have)
		return "", fmt.Errorf("no primers for language %q (available: %s)", lang, strings.Join(have, ", "))
	}

	m := Model{cat: cat, lang: lang}
	var b strings.Builder
	count := 0
	for _, category := range previewCategoryOrder {
		pr, ok := cat.PrimerByCategoryAndLang(category, lang)
		if !ok {
			continue
		}
		count++
		m.primer = pr
		b.WriteString(titleStyle.Render("════ "+pr.Title+" ════") + "\n\n")
		for _, pg := range m.primerPages() {
			b.WriteString(titleStyle.Render("• "+pg.heading) + "\n\n")
			b.WriteString(pg.body + "\n\n")
		}
		b.WriteString("\n")
	}
	b.WriteString(dimStyle.Render(fmt.Sprintf("(%d primers for %s)", count, langLabel(lang))) + "\n")
	return b.String(), nil
}

// PreviewAdvancedTopics renders a language's Stage-2 Advanced Topics (explainer +
// each exercise's broken/fixed code, syntax-highlighted) to one string — the spike
// preview for the reference / spot-the-bug UX, reusing the primer highlighter.
func PreviewAdvancedTopics(lang string) (string, error) {
	cat, err := content.Load()
	if err != nil {
		return "", err
	}
	topics := cat.AdvancedTopicsByLang(lang)
	if len(topics) == 0 {
		return "", fmt.Errorf("no advanced topics for language %q", lang)
	}
	var b strings.Builder
	for _, at := range topics {
		b.WriteString(titleStyle.Render("════ "+at.Title+" ════") + "  " + dimStyle.Render("["+at.Group+" · "+at.Tag+"]") + "\n\n")
		if strings.TrimSpace(at.Summary) != "" {
			b.WriteString(strings.TrimRight(at.Summary, "\n") + "\n\n")
		}
		for i, ex := range at.Exercises {
			b.WriteString(titleStyle.Render(fmt.Sprintf("Exercise %d — %s", i+1, ex.Kind)) + "\n")
			b.WriteString(strings.TrimRight(ex.Prompt, "\n") + "\n\n")
			b.WriteString(dimStyle.Render("  broken:") + "\n")
			b.WriteString(highlightCode(indentLines(ex.BrokenCode, "      "), lang) + "\n\n")
			b.WriteString(dimStyle.Render("  ▸ reveal — the bug:") + "\n")
			b.WriteString("    " + strings.ReplaceAll(strings.TrimRight(ex.Bug, "\n"), "\n", "\n    ") + "\n\n")
			b.WriteString(dimStyle.Render("  ▸ reveal — the fix:") + "\n")
			b.WriteString(highlightCode(indentLines(ex.FixedCode, "      "), lang) + "\n")
			grade := "reveal-only"
			if ex.Check == "tests" {
				grade = "auto-graded now (tests)"
			} else if ex.Check != "" && ex.Check != "none" {
				grade = "auto-grades when " + lang + " toolchain lands (" + ex.Check + ": " + ex.Signal + ")"
			}
			b.WriteString(dimStyle.Render("    grading: "+grade) + "\n\n")
		}
	}
	return b.String(), nil
}
