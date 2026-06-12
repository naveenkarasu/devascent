import './style.css';
import './workbench.css';
import { applyEditorTheme } from './monaco-setup';
import { mountHome } from './home';
import { mountBench } from './workbench';
import { mountOrientation } from './orientation';
import { mountTutorial } from './tutorial';
import { mountDevLit } from './devlit';
import { mountAdvanced } from './advanced';
import { GradedLanguages, GetLangPresence, GetVersion } from '../wailsjs/go/main/App';

// ── Theme registry ────────────────────────────────────────────
type ThemeId = 'prismatic' | 'tokyo' | 'gruvbox';
const THEMES: { id: ThemeId; label: string; sw: string }[] = [
  { id: 'prismatic', label: 'Prismatic', sw: '#61AFEF' },
  { id: 'tokyo', label: 'Tokyo', sw: '#7aa2f7' },
  { id: 'gruvbox', label: 'Gruvbox', sw: '#fabd2f' },
];
const STORE_KEY = 'devascent-theme';

function currentTheme(): ThemeId {
  const q = new URLSearchParams(location.search).get('theme') as ThemeId | null;
  if (q && THEMES.some((x) => x.id === q)) return q;
  const t = localStorage.getItem(STORE_KEY) as ThemeId | null;
  return t && THEMES.some((x) => x.id === t) ? t : 'prismatic';
}
function applyTheme(id: ThemeId) {
  document.documentElement.dataset.theme = id;
  localStorage.setItem(STORE_KEY, id);
  document
    .querySelectorAll<HTMLButtonElement>('.themesw button')
    .forEach((b) => b.classList.toggle('on', b.dataset.theme === id));
  applyEditorTheme();
}

// ── View router ───────────────────────────────────────────────
type ViewId = 'home' | 'orientation' | 'tutorial' | 'devlit' | 'bench' | 'advanced';
const NAV: { id: ViewId; label: string }[] = [
  { id: 'home', label: 'Home' },
  { id: 'orientation', label: 'Entrance Test' },
  { id: 'tutorial', label: 'Tutorial Island' },
  { id: 'devlit', label: 'Dev-Literacy' },
  { id: 'bench', label: 'Bench' },
  { id: 'advanced', label: 'Advanced' },
];
let view: ViewId = 'home';
let lang = 'python';
let dispose: (() => void) | null = null;

function render() {
  if (dispose) {
    dispose();
    dispose = null;
  }
  const root = document.getElementById('view') as HTMLElement;
  root.innerHTML = '';
  if (view === 'home') dispose = mountHome(root, lang);
  else if (view === 'orientation') dispose = mountOrientation(root, lang);
  else if (view === 'tutorial') dispose = mountTutorial(root, lang);
  else if (view === 'devlit') dispose = mountDevLit(root);
  else if (view === 'advanced') dispose = mountAdvanced(root, lang);
  else dispose = mountBench(root, lang);
  document
    .querySelectorAll<HTMLButtonElement>('.nav button')
    .forEach((b) => b.classList.toggle('on', b.dataset.view === view));
  // live status segments
  const navItem = NAV.find((n) => n.id === view);
  const sv = document.getElementById('stView');
  if (sv) sv.textContent = navItem ? navItem.label : view;
  const sl = document.getElementById('stLang');
  if (sl) sl.textContent = lang;
}

// ── Shell ─────────────────────────────────────────────────────
document.querySelector('#app')!.innerHTML = `
  <div class="prompt">
    <div class="seg lead"><span class="ico">&#9608;</span> DEVASCENT</div>
    <div class="nav">
      ${NAV.map((n) => `<button data-view="${n.id}">${n.label}</button>`).join('')}
    </div>
    <div class="seg cap"><span class="ico" style="color:var(--blue)">&#955;</span>
      <select id="langsel" class="langsel"></select></div>
    <div class="spacer"></div>
    <div class="themesw"><span class="tlabel">Theme</span>
      ${THEMES.map(
        (t) =>
          `<button data-theme="${t.id}"><span class="sw" style="background:${t.sw}"></span>${t.label}</button>`,
      ).join('')}
    </div>
  </div>

  <div id="view" class="viewhost"></div>

  <div class="status"><span class="s-accent">&#9670; <span id="stView">Home</span></span>
    <span>&#955; <span id="stLang">python</span></span><span>graded by your local toolchain</span>
    <div class="right"><span>UTF-8</span><span id="stVer"></span></div></div>
`;

// version (ldflags-injected at release; "dev" locally)
void (async () => {
  const v = await GetVersion();
  const sv = document.getElementById('stVer');
  if (sv) sv.textContent = v;
})();

// theme
applyTheme(currentTheme());
document
  .querySelectorAll<HTMLButtonElement>('.themesw button')
  .forEach((b) => b.addEventListener('click', () => applyTheme(b.dataset.theme as ThemeId)));

// nav
document.querySelectorAll<HTMLButtonElement>('.nav button').forEach((b) =>
  b.addEventListener('click', () => {
    view = b.dataset.view as ViewId;
    render();
  }),
);
window.addEventListener('devascent-nav', (e) => {
  view = (e as CustomEvent).detail as ViewId;
  render();
});
// profile picker → switch the whole app to that language's slot
window.addEventListener('devascent-lang', (e) => {
  lang = (e as CustomEvent).detail as string;
  const s = document.getElementById('langsel') as HTMLSelectElement | null;
  if (s) s.value = lang;
  render();
});

// language selector (drives all views; re-mounts the current view on change).
// Options carry the toolchain presence mark (✓ installed / ✗ not detected) —
// browsing stays open for every language; grading surfaces gate per-view.
const sel = document.getElementById('langsel') as HTMLSelectElement;
void (async () => {
  const langs = await GradedLanguages();
  sel.innerHTML = langs.map((l) => `<option value="${l}">${l}</option>`).join('');
  sel.value = lang;
  render();
  const statuses = await GetLangPresence();
  statuses.forEach((s) => {
    const opt = sel.querySelector(`option[value="${s.lang}"]`) as HTMLOptionElement | null;
    if (opt) opt.textContent = `${s.lang} ${s.status === 'available' ? '✓' : '✗'}`;
  });
})();
sel.addEventListener('change', () => {
  lang = sel.value;
  render();
});
