package tui

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// codeTheme is the chroma style used to syntax-highlight primer code. It is a
// dark-background theme so token colors read well on a black terminal. Changing
// this one constant reskins every primer snippet.
const codeTheme = "monokai"

var (
	chromaStyle     = styles.Get(codeTheme)
	chromaFormatter = formatters.Get("terminal256") // 256-color ANSI; foreground only (won't paint the bg)
	// hlCache memoizes highlighted snippets. Primer content is static, so the same
	// (lang, code) is re-highlighted on every render otherwise. Safe without a
	// mutex: bubbletea runs View and Update on a single goroutine.
	hlCache = map[string]string{}
)

// chromaLexerName maps a DevAscent language key to a chroma lexer name. Empty means
// "no highlighter" → fall back to the flat cyan codeStyle.
func chromaLexerName(lang string) string {
	switch lang {
	case "python":
		return "python"
	case "java":
		return "java"
	case "cpp":
		return "cpp"
	case "csharp":
		return "csharp"
	case "javascript":
		return "javascript"
	case "typescript":
		return "typescript"
	case "rust":
		return "rust"
	case "go":
		return "go"
	}
	return ""
}

// highlightCode returns code with ANSI syntax coloring for the given language.
// It falls back to the flat cyan codeStyle whenever a lexer/formatter/style is
// missing or highlighting errors — a primer must never render worse than plain.
func highlightCode(code, lang string) string {
	name := chromaLexerName(lang)
	if name == "" || chromaStyle == nil || chromaFormatter == nil {
		return codeStyle.Render(code)
	}
	key := lang + "\x00" + code
	if out, ok := hlCache[key]; ok {
		return out
	}
	lexer := lexers.Get(name)
	if lexer == nil {
		return codeStyle.Render(code)
	}
	lexer = chroma.Coalesce(lexer)
	it, err := lexer.Tokenise(nil, code)
	if err != nil {
		return codeStyle.Render(code)
	}
	var buf bytes.Buffer
	if err := chromaFormatter.Format(&buf, chromaStyle, it); err != nil {
		return codeStyle.Render(code)
	}
	out := strings.TrimRight(buf.String(), "\n")
	hlCache[key] = out
	return out
}

// indentLines prefixes every line of s with prefix (used to inset op snippets
// before highlighting, so the indentation is preserved as plain whitespace).
func indentLines(s, prefix string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i := range lines {
		lines[i] = prefix + lines[i]
	}
	return strings.Join(lines, "\n")
}

// codeIndentThreshold is the relative indent (after YAML block-scalar stripping)
// that separates a worked-example's prose (indent ≤ 2) from its code block
// (indent ≥ 4). See splitExample.
const codeIndentThreshold = 4

// isCodeLine reports whether a worked-example line belongs to the code block
// (indented at least codeIndentThreshold spaces). Blank lines are not code.
func isCodeLine(ln string) bool {
	if strings.TrimSpace(ln) == "" {
		return false
	}
	return len(ln)-len(strings.TrimLeft(ln, " ")) >= codeIndentThreshold
}

// splitExample separates a worked-example string into (intro prose, code block,
// outro prose). Each example has exactly one indented code block; the block is
// the span from the first to the last code line (blank lines inside it included).
// If no code block is detected, the whole thing is returned as intro.
func splitExample(example string) (intro, code, outro string) {
	lines := strings.Split(strings.TrimRight(example, "\n"), "\n")
	start, end := -1, -1
	for i, ln := range lines {
		if isCodeLine(ln) {
			if start == -1 {
				start = i
			}
			end = i
		}
	}
	if start == -1 {
		return example, "", ""
	}
	return strings.Join(lines[:start], "\n"),
		strings.Join(lines[start:end+1], "\n"),
		strings.Join(lines[end+1:], "\n")
}

// renderExampleBody returns the worked example with its code block syntax-
// highlighted and its prose left plain.
func renderExampleBody(example, lang string) string {
	intro, code, outro := splitExample(example)
	if strings.TrimSpace(code) == "" {
		return strings.TrimRight(example, "\n") // no code block detected; render as-is
	}
	var parts []string
	if strings.TrimSpace(intro) != "" {
		parts = append(parts, strings.TrimRight(intro, "\n"))
	}
	parts = append(parts, highlightCode(code, lang))
	if strings.TrimSpace(outro) != "" {
		parts = append(parts, strings.TrimLeft(strings.TrimRight(outro, "\n"), "\n"))
	}
	return strings.Join(parts, "\n\n")
}
