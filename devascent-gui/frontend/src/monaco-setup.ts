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

// Build the 'devascent' Monaco theme from the active CSS theme variables so the
// editor tracks the app theme. Call again after a theme switch.
export function applyEditorTheme(): void {
  const css = getComputedStyle(document.documentElement);
  const v = (name: string) => css.getPropertyValue(name).trim() || '#000000';
  const tok = (name: string) => v(name).replace('#', '');
  monaco.editor.defineTheme('devascent', {
    base: 'vs-dark',
    inherit: true,
    rules: [
      { token: 'comment', foreground: tok('--textMuted'), fontStyle: 'italic' },
      { token: 'keyword', foreground: tok('--magenta') },
      { token: 'string', foreground: tok('--green') },
      { token: 'number', foreground: tok('--yellow') },
      { token: 'type', foreground: tok('--cyan') },
      { token: 'identifier', foreground: tok('--textPrimary') },
    ],
    colors: {
      'editor.background': v('--base'),
      'editor.foreground': v('--textPrimary'),
      'editorLineNumber.foreground': v('--textMuted'),
      'editorCursor.foreground': v('--blue'),
      'editor.selectionBackground': v('--panelStrong'),
      'editor.lineHighlightBackground': v('--panel'),
      'editorWidget.background': v('--panel'),
      'editorWidget.border': v('--borderSubtle'),
    },
  });
  monaco.editor.setTheme('devascent');
}

export { monaco };
