package engine

import "strings"

// Place computes the intake placement from the score (percentage of items
// passed), then applies the self-report band clamp. Single source of truth for
// routing, shared by the TUI and the GUI facade.
//
// Bands: ≥80% → "test-out" · 40–79% → "dev-literacy" · <40% → "tutorial-full".
// Clamp: level "never" always lands in tutorial-full (acing EASY ≠ bench-ready);
// level "regularly" never lands in beginner tutorial (upgraded to dev-literacy).
func Place(passed, total int, level string) string {
	if total == 0 {
		total = 1
	}
	pct := float64(passed) / float64(total)
	var place string
	switch {
	case pct >= 0.8:
		place = "test-out"
	case pct >= 0.4:
		place = "dev-literacy"
	default:
		place = "tutorial-full"
	}
	switch level {
	case "never":
		place = "tutorial-full"
	case "regularly":
		if place == "tutorial-full" {
			place = "dev-literacy"
		}
	}
	return place
}

// SpecMatch grades a free-text spec answer: each Required entry is a synonym
// GROUP ("list|array|collection") and the answer must contain at least one
// synonym from every group. Avoids false-negatives on correct paraphrases.
func SpecMatch(ans string, required []string) bool {
	low := strings.ToLower(ans)
	for _, group := range required {
		hit := false
		for _, syn := range strings.Split(group, "|") {
			syn = strings.TrimSpace(strings.ToLower(syn))
			if syn != "" && strings.Contains(low, syn) {
				hit = true
				break
			}
		}
		if !hit {
			return false
		}
	}
	return true
}

// DevMatch grades a typed command (Dev-Literacy track): PASS if the first token
// is one of commands AND every flag appears — OR the whole line (normalised)
// equals an accept form. Lenient on args/quoting (a checker, not a shell).
func DevMatch(ans string, commands, flags, accept []string) bool {
	norm := func(s string) string { return strings.Join(strings.Fields(strings.ToLower(s)), " ") }
	toks := strings.Fields(strings.ToLower(ans))
	if len(toks) > 0 {
		has := func(tok string) bool {
			for _, x := range toks {
				if x == tok {
					return true
				}
			}
			return false
		}
		for _, c := range commands {
			if toks[0] != strings.ToLower(strings.TrimSpace(c)) {
				continue
			}
			ok := true
			for _, fl := range flags {
				if !has(strings.ToLower(strings.TrimSpace(fl))) {
					ok = false
					break
				}
			}
			if ok {
				return true
			}
		}
	}
	a := norm(ans)
	for _, acc := range accept {
		if a == norm(acc) {
			return true
		}
	}
	return false
}
