package mentor

// Deterministic output enforcement — the engine-side half of the guardrails.
// A response that violates its tier's contract is REJECTED (template fallback
// + token refund), never trimmed into compliance: a mentor that ignores the
// rules once will ignore them again, and silently sanitized output teaches
// the player nothing.

import (
	"fmt"
	"regexp"
	"strings"
)

var fenceRe = regexp.MustCompile("(?s)```.*?(```|$)")

// statementLike marks lines that look like real code rather than prose.
var statementLike = regexp.MustCompile(`(?:[;{}]\s*$|^\s*(?:def |func |class |return |for\s*\(|while\s*\(|if\s*\(.*\)\s*{)|=\s*[^=].*[;)]\s*$)`)

func wordCount(s string) int {
	return len(strings.Fields(s))
}

// codeBlockiness counts the longest run of consecutive statement-like lines.
func codeBlockiness(s string) int {
	run, worst := 0, 0
	for _, line := range strings.Split(s, "\n") {
		if statementLike.MatchString(strings.TrimSpace(line)) {
			run++
			if run > worst {
				worst = run
			}
		} else {
			run = 0
		}
	}
	return worst
}

// Validate enforces the per-kind output contract; a non-nil error means
// "fall back to templates and refund".
func Validate(kind Kind, lang, text string) error {
	t := strings.TrimSpace(text)
	if t == "" {
		return fmt.Errorf("empty response")
	}
	switch kind {
	case KindStrategy:
		if fenceRe.MatchString(t) {
			return fmt.Errorf("strategy hint contains a code block")
		}
		if codeBlockiness(t) >= 3 {
			return fmt.Errorf("strategy hint contains code-like lines")
		}
		if wordCount(t) > 200 { // contract says 120; reject only at hard runaway
			return fmt.Errorf("strategy hint too long (%d words)", wordCount(t))
		}
	case KindWalkthrough:
		for _, m := range fenceRe.FindAllString(t, -1) {
			tag := strings.ToLower(strings.TrimSpace(strings.SplitN(strings.TrimPrefix(m, "```"), "\n", 2)[0]))
			if tag != "" && langMatches(tag, lang) {
				return fmt.Errorf("walkthrough contains a %s code block", lang)
			}
		}
		if wordCount(t) > 600 {
			return fmt.Errorf("walkthrough too long (%d words)", wordCount(t))
		}
	case KindFollowup:
		if !strings.Contains(t, "?") {
			return fmt.Errorf("follow-up is not a question")
		}
		if wordCount(t) > 60 {
			return fmt.Errorf("follow-up too long (%d words)", wordCount(t))
		}
		if fenceRe.MatchString(t) {
			return fmt.Errorf("follow-up contains a code block")
		}
	case KindReview:
		if wordCount(t) > 150 {
			return fmt.Errorf("review too long (%d words)", wordCount(t))
		}
	case KindStandup:
		if fenceRe.MatchString(t) || codeBlockiness(t) >= 3 {
			return fmt.Errorf("standup contains code")
		}
		if wordCount(t) > 250 {
			return fmt.Errorf("standup too long (%d words)", wordCount(t))
		}
	case KindDiscuss:
		if fenceRe.MatchString(t) {
			return fmt.Errorf("plan contains a code block")
		}
		if wordCount(t) > 120 {
			return fmt.Errorf("plan too long (%d words)", wordCount(t))
		}
	}
	return nil
}

// langMatches maps fence tags to the game's language IDs.
func langMatches(tag, lang string) bool {
	aliases := map[string][]string{
		"python":     {"python", "py", "python3"},
		"go":         {"go", "golang"},
		"java":       {"java"},
		"csharp":     {"csharp", "cs", "c#"},
		"javascript": {"javascript", "js", "node"},
		"typescript": {"typescript", "ts"},
		"rust":       {"rust", "rs"},
	}
	for _, a := range aliases[lang] {
		if tag == a {
			return true
		}
	}
	return false
}
