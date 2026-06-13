// Trimmed Monaco build: the core editor + all editor features, plus ONLY the 7
// graded languages' Monarch grammars (syntax highlighting). We deliberately skip
// editor.main — it bundles ~80 grammars and the ts/json/css/html language
// services, none of which we use (the real compiler grades the player's code).
import 'monaco-editor/esm/vs/editor/editor.all';
import 'monaco-editor/esm/vs/basic-languages/python/python.contribution';
import 'monaco-editor/esm/vs/basic-languages/go/go.contribution';
import 'monaco-editor/esm/vs/basic-languages/csharp/csharp.contribution';
import 'monaco-editor/esm/vs/basic-languages/java/java.contribution';
import 'monaco-editor/esm/vs/basic-languages/rust/rust.contribution';
import 'monaco-editor/esm/vs/basic-languages/javascript/javascript.contribution';
import 'monaco-editor/esm/vs/basic-languages/typescript/typescript.contribution';
import 'monaco-editor/esm/vs/basic-languages/cpp/cpp.contribution'; // reference-only viewing
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';

// The base editor worker is enough: syntax highlighting (Monarch) lives in the
// main bundle; the worker only backs the editor itself. We don't ship the
// per-language IntelliSense workers — the player's code is graded by the real
// compiler, not the browser.
self.MonacoEnvironment = {
  getWorker: () => new editorWorker(),
};

const LANG_MAP: Record<string, string> = {
  python: 'python',
  go: 'go',
  csharp: 'csharp',
  javascript: 'javascript',
  typescript: 'typescript',
  java: 'java',
  rust: 'rust',
  cpp: 'cpp',
};

export function monacoLang(lang: string): string {
  return LANG_MAP[lang] ?? 'plaintext';
}

// makeEditor creates a themed code editor in container, seeded with value.
export function makeEditor(
  container: HTMLElement,
  lang: string,
  value: string,
): monaco.editor.IStandaloneCodeEditor {
  return monaco.editor.create(container, {
    value,
    language: monacoLang(lang),
    theme: 'devascent',
    automaticLayout: true,
    fontSize: 14,
    fontFamily: 'Cascadia Mono, Consolas, monospace',
    minimap: { enabled: false },
    scrollBeyondLastLine: false,
    tabSize: 4,
    renderWhitespace: 'none',
    padding: { top: 10 },
  });
}

// ── Authored per-theme syntax palettes ───────────────────────────────────
// The editor is the biggest visual surface, so each theme ships its REAL
// upstream syntax palette instead of colors derived from the UI vars — the
// derivation made every theme's code look like the same theme dimmed.
// Editor CHROME (background, gutter, selection, cursor) still tracks the CSS
// vars in applyEditorTheme below; only token colors are authored here.
interface TokenPalette {
  keyword: string;
  string: string;
  number: string;
  comment: string;
  type: string; // types / classes
  func: string; // function & method names, annotations
  variable: string; // plain identifiers
  operator: string;
  constant: string;
}

const TOKEN_PALETTES: Record<string, TokenPalette> = {
  // One-Dark-ish: formalizes the derived look prismatic already had (its
  // accents are One Dark's #61AFEF/#98C379/#C678DD), plus authored function/
  // operator/constant colors the derived path never set.
  prismatic: {
    keyword: 'C678DD',
    string: '98C379',
    number: 'E5C07B',
    comment: '8794A8',
    type: '56B6C2',
    func: '61AFEF',
    variable: 'F5F7FA',
    operator: 'C8D1DD',
    constant: 'D19A66',
  },
  // Dracula (upstream: dracula/dracula-theme). Pink keywords, lime functions,
  // cyan types; comment lifted slightly off the dim #6272a4 toward readable.
  dracula: {
    keyword: 'FF79C6',
    string: 'F1FA8C',
    number: 'BD93F9',
    comment: '7081B9',
    type: '8BE9FD',
    func: '50FA7B',
    variable: 'F8F8F2',
    operator: 'FF79C6',
    constant: 'BD93F9',
  },
  // Monokai (classic Sublime tmTheme). Pink keywords, lime functions, cyan
  // types, purple constants; comment lifted slightly off #75715E.
  monokai: {
    keyword: 'F92672',
    string: 'E6DB74',
    number: 'AE81FF',
    comment: '88826C',
    type: '66D9EF',
    func: 'A6E22E',
    variable: 'F8F8F2',
    operator: 'F92672',
    constant: 'AE81FF',
  },
};

// Map a palette onto the token names the bundled Monarch grammars emit
// (Monaco matches rules by dotted prefix, so 'number' also covers
// 'number.float', 'string' covers 'string.escape', …).
function tokenRules(p: TokenPalette): monaco.editor.ITokenThemeRule[] {
  return [
    { token: 'comment', foreground: p.comment, fontStyle: 'italic' },
    { token: 'keyword', foreground: p.keyword },
    { token: 'string', foreground: p.string },
    { token: 'number', foreground: p.number },
    { token: 'type', foreground: p.type },
    { token: 'type.identifier', foreground: p.type },
    { token: 'namespace', foreground: p.type },
    { token: 'function', foreground: p.func },
    { token: 'annotation', foreground: p.func },
    { token: 'identifier', foreground: p.variable },
    { token: 'variable', foreground: p.variable },
    { token: 'operator', foreground: p.operator },
    { token: 'delimiter', foreground: p.variable },
    { token: 'constant', foreground: p.constant },
  ];
}

// Build the 'devascent' Monaco theme: chrome colors from the active CSS theme
// variables, token colors from the authored palette for the active data-theme
// (falling back to var-derived tokens for unknown themes). Call again after a
// theme switch.
export function applyEditorTheme(): void {
  const css = getComputedStyle(document.documentElement);
  const v = (name: string) => css.getPropertyValue(name).trim() || '#000000';
  const tok = (name: string) => v(name).replace('#', '');
  const palette = TOKEN_PALETTES[document.documentElement.dataset.theme ?? ''];
  const rules = palette
    ? tokenRules(palette)
    : [
        // var-derived fallback for themes without an authored palette
        { token: 'comment', foreground: tok('--textMuted'), fontStyle: 'italic' },
        { token: 'keyword', foreground: tok('--magenta') },
        { token: 'string', foreground: tok('--green') },
        { token: 'number', foreground: tok('--yellow') },
        { token: 'type', foreground: tok('--cyan') },
        { token: 'identifier', foreground: tok('--textPrimary') },
      ];
  monaco.editor.defineTheme('devascent', {
    base: 'vs-dark',
    inherit: true,
    rules,
    colors: {
      'editor.background': v('--base'),
      'editor.foreground': v('--textPrimary'),
      'editorLineNumber.foreground': v('--textMuted'),
      'editorCursor.foreground': v('--accent'),
      'editor.selectionBackground': v('--selection'),
      'editor.lineHighlightBackground': v('--panel'),
      'editorWidget.background': v('--panel'),
      'editorWidget.border': v('--borderSubtle'),
    },
  });
  monaco.editor.setTheme('devascent');
}

export { monaco };
