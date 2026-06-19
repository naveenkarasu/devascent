package mentor

// Context-pack construction: the ONLY place prompts are built, so the
// guardrails table has a single enforcement point. BuildPrompt is also what
// the "here's exactly what gets sent" preview shows the player.

import (
	"fmt"
	"regexp"
	"strings"
)

// absPathRe scrubs machine paths (Windows drives and Unix homes) from
// anything that goes over the wire — player code can contain them in
// comments, and they're nobody's business.
var absPathRe = regexp.MustCompile(`(?i)([a-z]:\\[^\s"']+|/(?:home|users)/[^\s"']+)`)

func scrub(s string) string {
	return absPathRe.ReplaceAllString(s, "<path>")
}

// BuildPrompt renders the tier-scoped context pack. Everything the AI will
// ever see of the player's run flows through here.
func BuildPrompt(req Request) string {
	switch req.Kind {
	case KindStandup:
		return buildStandupPrompt(req)
	case KindDiscuss:
		return buildDiscussPrompt(req)
	}
	var b strings.Builder
	switch req.Kind {
	case KindStrategy:
		b.WriteString("You are a senior developer mentoring an apprentice on a coding problem.\n")
		b.WriteString("Describe the APPROACH that solves it: the key insight and which data structure or technique to reach for, in at most 120 words.\n")
		b.WriteString("HARD RULES: no code, no code blocks, no step-by-step recipe — point the direction, don't walk the path. Plain text only.\n")
	case KindWalkthrough:
		b.WriteString("You are a senior developer mentoring an apprentice on a coding problem.\n")
		b.WriteString("Give a numbered step-by-step walkthrough of the solution. Pseudocode is allowed.\n")
		fmt.Fprintf(&b, "HARD RULES: do NOT write compilable %s code — the apprentice must still type the real solution themselves. Plain text and pseudocode only.\n", req.Lang)
	case KindFollowup:
		b.WriteString("You are a senior developer reviewing an apprentice's write-up of their accepted solution.\n")
		b.WriteString("Ask exactly ONE short probing question (at most 40 words) that tests whether they really understand their own code. Output only the question.\n")
	case KindReview:
		b.WriteString("You are a senior developer giving a quick review of an apprentice's accepted solution.\n")
		b.WriteString("In at most 80 words: name ONE strength and ONE concrete improvement. Plain text only.\n")
	}

	fmt.Fprintf(&b, "\nProblem: %s (%s, %s)\n", req.Title, req.Category, req.Difficulty)
	b.WriteString("Statement:\n")
	b.WriteString(scrub(req.Prompt))
	b.WriteString("\n")
	if req.PlayerCode != "" {
		fmt.Fprintf(&b, "\nThe apprentice's current %s code:\n%s\n", req.Lang, scrub(req.PlayerCode))
	}
	if req.Kind == KindStrategy || req.Kind == KindWalkthrough {
		if req.FailedRuns > 0 {
			fmt.Fprintf(&b, "\nThey have failed %d grading attempts", req.FailedRuns)
			if req.FirstFail != "" {
				fmt.Fprintf(&b, "; the first failing test is named %q", scrub(req.FirstFail))
			}
			b.WriteString(".\n")
		}
	}
	if req.Writeup != "" && (req.Kind == KindFollowup || req.Kind == KindReview) {
		fmt.Fprintf(&b, "\nTheir write-up:\n%s\n", scrub(req.Writeup))
	}
	return b.String()
}

// buildStandupPrompt renders the team's pre-computed public status lines as a
// natural morning standup. The AI only ever sees the status facts we hand it —
// never tickets, tests, or solutions.
func buildStandupPrompt(req Request) string {
	var b strings.Builder
	fmt.Fprintf(&b, "You are %s, running a short daily standup at a software studio on day %d.\n", personaOr(req.Persona, "the engineering manager"), req.Day)
	b.WriteString("Turn the raw team status below into a crisp standup: ONE line per person (what they did, what's next, any blocker), then ONE short manager focus line.\n")
	b.WriteString("HARD RULES: at most 120 words, plain text, no code, no preamble. Keep each person to one line.\n\n")
	b.WriteString("Team status:\n")
	for _, s := range req.Status {
		b.WriteString("- " + scrub(s) + "\n")
	}
	return b.String()
}

// buildDiscussPrompt asks a delegated teammate for a short plan + estimate. They
// see the ticket's public Definition-of-Ready (title, description, acceptance),
// exactly what a real engineer gets — never the hidden grading tests.
func buildDiscussPrompt(req Request) string {
	var b strings.Builder
	fmt.Fprintf(&b, "You are %s, a teammate who has just been handed a ticket to own.\n", personaOr(req.Persona, "a software engineer"))
	b.WriteString("In 2-3 sentences, say how you'll approach it and give a rough estimate in days. Be concrete and brief.\n")
	b.WriteString("HARD RULES: at most 80 words, plain text, no code. Speak in the first person.\n")
	fmt.Fprintf(&b, "\nTicket: %s", req.Title)
	if req.Priority != "" {
		fmt.Fprintf(&b, " (priority: %s)", req.Priority)
	}
	b.WriteString("\n")
	if req.Prompt != "" {
		b.WriteString(scrub(req.Prompt) + "\n")
	}
	if len(req.Acceptance) > 0 {
		b.WriteString("Acceptance criteria:\n")
		for _, a := range req.Acceptance {
			b.WriteString("- " + scrub(a) + "\n")
		}
	}
	return b.String()
}

func personaOr(p, def string) string {
	if strings.TrimSpace(p) != "" {
		return p
	}
	return def
}
