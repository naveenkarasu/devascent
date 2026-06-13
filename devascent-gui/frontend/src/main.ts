import './style.css';
import './workbench.css';
import brandLogo from './assets/brand_logo.png'; // channel/studio logo — loading splash only
import appIcon from './assets/appicon.png'; // app logo — in-app chrome (titlebar)
import { applyEditorTheme } from './monaco-setup';
import { langIcon } from './langicons';
import { mountHome } from './home';
import { mountBench } from './workbench';
import { mountOrientation } from './orientation';
import { mountTutorial } from './tutorial';
import { mountDevLit } from './devlit';
import { mountAdvanced } from './advanced';
import { GetLanguages, GetLangPresence, GetProgress, GetVersion } from '../wailsjs/go/main/App';
import { guiapi } from '../wailsjs/go/models';

// ── Boot splash ───────────────────────────────────────────────
// Mounted before any view (and before Monaco wakes up); dismissed once the
// first render has completed AND at least 1s has elapsed (no strobe), with a
// 300ms fade. The overlay swallows all interaction while visible.
const splashAt = Date.now();
const splash = document.createElement('div');
splash.className = 'splash';
splash.innerHTML = `
  <img class="splash-logo" src="${brandLogo}" alt="" draggable="false">
  <div class="splash-name">DevAscent</div>
  <div class="splash-bar"><div class="splash-shimmer"></div></div>`;
document.body.appendChild(splash);
let splashDismissed = false;
function dismissSplash(): void {
  if (splashDismissed) return;
  splashDismissed = true;
  // Hold the loading screen for a standard ~3s splash so the brand moment
  // registers and never flashes by on a fast boot, then fade out.
  const wait = Math.max(0, 3000 - (Date.now() - splashAt));
  window.setTimeout(() => {
    splash.classList.add('out');
    window.setTimeout(() => splash.remove(), 320);
  }, wait);
}

// ── Theme registry ────────────────────────────────────────────
type ThemeId = 'prismatic' | 'dracula' | 'monokai';
const THEMES: { id: ThemeId; label: string; sw: string }[] = [
  { id: 'prismatic', label: 'Prismatic', sw: '#61AFEF' },
  { id: 'dracula', label: 'Dracula', sw: '#bd93f9' },
  { id: 'monokai', label: 'Monokai', sw: '#a6e22e' },
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

// language catalog (graded + reference-only); loaded once at boot
let langs: guiapi.LangInfo[] = [];
const isGraded = (id: string): boolean => langs.find((l) => l.id === id)?.graded !== false;
const optLabel = (l: guiapi.LangInfo): string => (l.graded ? l.label : `${l.label} (reference)`);

// ── Progression locks (hard-gate the header nav) ──────────────
// view → lock reason; missing = unlocked. Starts fully locked (only Home +
// Entrance Test) until the first save read lands — no flash of open nav.
const LOCK_INTAKE = 'Finish the Entrance Test first';
let locks: Partial<Record<ViewId, string>> = {
  tutorial: LOCK_INTAKE,
  devlit: LOCK_INTAKE,
  bench: LOCK_INTAKE,
  advanced: LOCK_INTAKE,
};

// computeLocks derives the per-entry gate from the save. Placement-aware, not
// strictly linear:
//  - no placement   → only Home + Entrance Test
//  - tutorial-full  → + Tutorial
//  - dev-literacy   → + Tutorial AND Dev-Literacy
//  - test-out       → everything (they aced the intake)
//  - Bench/Advanced → once the run shows bench progress (banked > 0) or a
//    returning player (step0Met). A brand-new tutorial player reaches the
//    bench through the tutorial's own hand-off — programmatic nav (custom
//    devascent-nav events from views) is deliberately not gated.
function computeLocks(pr: guiapi.Progress): Partial<Record<ViewId, string>> {
  const out: Partial<Record<ViewId, string>> = {};
  // Reference-only language (C++): nothing gradeable can run, so the graded
  // tracks lock outright; the Bench/Advanced stay open for browsing.
  if (!isGraded(lang)) {
    const why = 'C++ is reference-only — switch to a graded language to play this track';
    out.orientation = why;
    out.tutorial = why;
    out.devlit = why;
    return out;
  }
  if (!pr.placement) {
    out.tutorial = LOCK_INTAKE;
    out.devlit = LOCK_INTAKE;
    out.bench = LOCK_INTAKE;
    out.advanced = LOCK_INTAKE;
    return out;
  }
  if (pr.placement === 'test-out') return out;
  if (pr.placement !== 'dev-literacy')
    out.devlit = 'Unlocks with a Dev-Literacy placement on the Entrance Test';
  if (!(pr.banked > 0 || pr.step0Met)) {
    const why =
      pr.placement === 'dev-literacy'
        ? 'Finish the Dev-Literacy brush-up first — it hands off to the Bench'
        : 'Finish Tutorial Island first — it hands off to the Bench';
    out.bench = why;
    out.advanced = why;
  }
  return out;
}

function applyLocks(): void {
  document.querySelectorAll<HTMLButtonElement>('.nav button').forEach((b) => {
    const reason = locks[b.dataset.view as ViewId] || '';
    b.disabled = !!reason;
    b.title = reason;
  });
}

let lockSeq = 0; // stale-read guard across rapid language switches
let locksLang = ''; // language the current `locks` were computed for
async function refreshLocks(): Promise<void> {
  const seq = ++lockSeq;
  // On a LANGUAGE SWITCH, lock the gated tracks pessimistically and
  // SYNCHRONOUSLY (before any click can land) so the previous language's
  // unlocked nav can't carry over while the new language's progress loads.
  // Same-language renders skip this (no flicker).
  if (lang !== locksLang) {
    locks = { tutorial: LOCK_INTAKE, devlit: LOCK_INTAKE, bench: LOCK_INTAKE, advanced: LOCK_INTAKE };
    applyLocks();
  }
  const pr = await GetProgress(lang);
  if (seq !== lockSeq) return;
  locksLang = lang;
  locks = computeLocks(pr);
  applyLocks();
}

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
  else dispose = mountBench(root, lang, isGraded(lang));
  document
    .querySelectorAll<HTMLButtonElement>('.nav button')
    .forEach((b) => b.classList.toggle('on', b.dataset.view === view));
  // live status segments
  const navItem = NAV.find((n) => n.id === view);
  const sv = document.getElementById('stView');
  if (sv) sv.textContent = navItem ? navItem.label : view;
  const sl = document.getElementById('stLang');
  if (sl) sl.textContent = lang;
  // language marks track the active language (header + status bar)
  const hi = document.getElementById('hdrLangIco');
  if (hi) hi.innerHTML = langIcon(lang);
  const si = document.getElementById('stLangIco');
  if (si) si.innerHTML = langIcon(lang);
  // recompute nav locks on every render — covers boot, language switches and
  // the post-orientation hand-off (placement is saved before its nav fires)
  void refreshLocks();
}

// ── Shell ─────────────────────────────────────────────────────
document.querySelector('#app')!.innerHTML = `
  <div class="prompt">
    <div class="seg lead"><img class="brandlogo-sm" src="${appIcon}" alt="" draggable="false"> DEVASCENT</div>
    <div class="nav">
      ${NAV.map((n) => `<button data-view="${n.id}">${n.label}</button>`).join('')}
    </div>
    <div class="seg cap"><span class="ico" id="hdrLangIco">${langIcon(lang)}</span>
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
    <span><span class="ico" id="stLangIco">${langIcon(lang)}</span> <span id="stLang">python</span></span><span>graded by your local toolchain</span>
    <div class="right"><span>UTF-8</span><span id="stVer"></span></div></div>
`;
applyLocks(); // start hard-locked; the save read in render() opens what's earned

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

// nav (locked entries are disabled buttons — clicks never land; the guard
// covers synthetic clicks too)
document.querySelectorAll<HTMLButtonElement>('.nav button').forEach((b) =>
  b.addEventListener('click', () => {
    const v = b.dataset.view as ViewId;
    if (locks[v]) return;
    view = v;
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
// Sourced from GetLanguages — reference-only rows render as "<label>
// (reference)" and stay selectable (browse works; grading is gated per-view).
// Options carry the toolchain presence mark (✓ installed / ✗ not detected).
const sel = document.getElementById('langsel') as HTMLSelectElement;
void (async () => {
  try {
    langs = (await GetLanguages()) || [];
    sel.innerHTML = langs
      .map((l) => `<option value="${l.id}">${optLabel(l)}</option>`)
      .join('');
    sel.value = lang;
    render();
  } finally {
    dismissSplash(); // first render is up (or boot failed) — fade the splash
  }
  const statuses = await GetLangPresence();
  statuses.forEach((s) => {
    const opt = sel.querySelector(`option[value="${s.lang}"]`) as HTMLOptionElement | null;
    const info = langs.find((l) => l.id === s.lang);
    if (opt && info) opt.textContent = `${optLabel(info)} ${s.status === 'available' ? '✓' : '✗'}`;
  });
})();
sel.addEventListener('change', () => {
  lang = sel.value;
  render();
});
