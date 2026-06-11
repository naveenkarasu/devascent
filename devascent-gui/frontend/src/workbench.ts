import { monaco, monacoLang, makeEditor } from './monaco-setup';
import {
  ListProblems,
  GetProblem,
  Grade,
  GetProgress,
  NextProblem,
  GetPrimer,
  CheckLang,
  RecheckLang,
  GetInstallGuide,
} from '../wailsjs/go/main/App';
import { guiapi } from '../wailsjs/go/models';
import { esc, renderVerdict, renderInstallGuide } from './util';

// mountBench renders the bench workbench into root for lang and returns a
// disposer (call before unmounting to release the Monaco editor).
export function mountBench(root: HTMLElement, lang: string): () => void {
  root.innerHTML = `
    <div class="bench">
      <aside class="sidebar">
        <div id="bScore" class="scorecard"></div>
        <div class="chiprow" id="bLists">
          <button class="fchip on" data-list="">All</button>
          <button class="fchip" data-list="blind75">Blind 75</button>
          <button class="fchip" data-list="neetcode150">NeetCode 150</button>
        </div>
        <div class="chiprow" id="bDiffs">
          <button class="fchip d-easy on" data-diff="easy">E</button>
          <button class="fchip d-medium on" data-diff="medium">M</button>
          <button class="fchip d-hard on" data-diff="hard">H</button>
        </div>
        <input id="bSearch" class="search" placeholder="Search problems…" />
        <div id="bCount" class="pcount"></div>
        <div id="bList" class="plist"></div>
      </aside>
      <main class="work">
        <div class="probhead">
          <span id="bTitle" class="ptitle-h">Loading…</span>
          <span id="bDiff" class="tag"></span>
          <span id="bCat" class="tag cat"></span>
          <span class="spacer"></span>
          <button id="bLearn" class="btn small">&#9656; Learn</button>
        </div>
        <div id="bDrawer" class="drawer hidden"></div>
        <div id="bPrompt" class="prompt-text"></div>
        <div id="bEditor" class="editor"></div>
        <div class="workbar">
          <button id="bGrade" class="btn primary">&#9654; Grade</button>
          <button id="bReset" class="btn">&#8634; Reset starter</button>
        </div>
        <div id="bVerdict" class="verdict"><div class="vmuted">Pick a problem to begin.</div></div>
      </main>
    </div>`;

  const el = (id: string) => root.querySelector('#' + id) as HTMLElement;
  const editor = makeEditor(el('bEditor'), lang, '');
  let problems: guiapi.ProblemSummary[] = [];
  let currentId = '';
  let currentCategory = '';
  let disposed = false; // set by the disposer; async continuations must bail
  let listFilter = ''; // '' | 'blind75' | 'neetcode150'
  const diffOn: Record<string, boolean> = { easy: true, medium: true, hard: true };
  let toolchainOK = true; // optimistic until the capability canary reports
  let toolchainReason = '';

  // ── Capability gate ──────────────────────────────────────────
  function showGateBanner(): void {
    (el('bGrade') as HTMLButtonElement).disabled = true;
    el('bVerdict').innerHTML = `
      <div class="vhead vfail">✖ ${esc(lang)} toolchain not available</div>
      <div class="cdetail">${esc(toolchainReason)} — browsing, prompts and primers still work; grading needs the real toolchain.</div>
      <div class="workbar">
        <button id="bInstall" class="btn">&#9656; Install guide</button>
        <button id="bRecheck" class="btn primary">&#8634; Re-check</button>
      </div>`;
    (el('bInstall') as HTMLButtonElement).addEventListener('click', () => void openInstallDrawer());
    (el('bRecheck') as HTMLButtonElement).addEventListener('click', () => void recheck());
  }

  function applyGate(st: guiapi.LangStatus): void {
    toolchainOK = st.status === 'available';
    toolchainReason = st.reason || `${lang} toolchain ${st.status}`;
    if (!toolchainOK) {
      showGateBanner();
    } else {
      (el('bGrade') as HTMLButtonElement).disabled = false;
    }
  }

  async function recheck(): Promise<void> {
    el('bVerdict').innerHTML = '<div class="vmuted">Re-checking the toolchain (compiles + runs a canary)…</div>';
    const st = await RecheckLang(lang);
    if (disposed) return;
    applyGate(st);
    if (toolchainOK) {
      el('bVerdict').innerHTML = `<div class="vhead vok">✓ ${esc(lang)} toolchain detected${st.version ? ' · ' + esc(st.version) : ''}</div><div class="cdetail">You're good — write your solution and press Grade.</div>`;
    }
  }

  async function openInstallDrawer(): Promise<void> {
    const g = await GetInstallGuide(lang);
    if (disposed) return;
    const d = el('bDrawer');
    d.innerHTML = `
      <div class="drawhead"><span>Install ${esc(g.label || lang)}</span><button class="btn small" id="bDrawX">✕</button></div>
      <div class="drawbody">${renderInstallGuide(g, toolchainReason)}
        <div class="workbar"><button id="bDrawRecheck" class="btn primary">&#8634; Re-check</button></div>
      </div>`;
    d.classList.remove('hidden');
    (root.querySelector('#bDrawX') as HTMLButtonElement).addEventListener('click', closeDrawer);
    (root.querySelector('#bDrawRecheck') as HTMLButtonElement).addEventListener('click', () => {
      closeDrawer();
      void recheck();
    });
  }

  function visible(): guiapi.ProblemSummary[] {
    const q = (el('bSearch') as HTMLInputElement).value.toLowerCase();
    return problems.filter(
      (p) =>
        (!listFilter || (p.lists || []).includes(listFilter)) &&
        diffOn[p.difficulty] !== false &&
        (!q || p.title.toLowerCase().includes(q) || p.id.toLowerCase().includes(q)),
    );
  }

  async function refreshScore(): Promise<void> {
    const pr = await GetProgress(lang);
    if (disposed) return;
    const pct = Math.min(100, Math.round((100 * pr.banked) / pr.bankTarget));
    el('bScore').innerHTML = `
      <div class="scrow">
        <span class="scval">${pr.banked}<span class="scmut">/${pr.bankTarget}</span></span><span class="sclabel">banked</span>
        <span class="scval">${pr.cats}<span class="scmut">/${pr.catTarget}</span></span><span class="sclabel">cats</span>
        <span class="scval">${pr.hard}<span class="scmut">/${pr.hardTarget}</span></span><span class="sclabel">hard</span>
      </div>
      <div class="scbar"><div class="scfill" style="width:${pct}%"></div></div>
      ${
        pr.step0Met
          ? `<div class="scmet">✓ Step 0 milestone met · ${esc(pr.track)}</div>`
          : `<div class="schint">Step 0: bank ${pr.bankTarget}, cover ${pr.catTarget} categories, ${pr.hardTarget} hard</div>`
      }`;
  }

  function renderList(): void {
    const items = visible();
    const banked = items.filter((p) => p.solved).length;
    el('bList').innerHTML = items
      .map(
        (p) => `<div class="prow ${p.id === currentId ? 'active' : ''}" data-id="${p.id}">
          <span class="ptick ${p.solved ? 'on' : ''}">${p.solved ? '✓' : ''}</span>
          <span class="ptitle">${esc(p.title)}</span>
          <span class="pdiff d-${esc(p.difficulty)}">${esc(p.difficulty)}</span></div>`,
      )
      .join('');
    el('bList')
      .querySelectorAll<HTMLElement>('.prow')
      .forEach((r) => r.addEventListener('click', () => void open(r.dataset.id as string)));
    el('bCount').textContent = `${items.length} problems · ${banked} banked`;
  }

  async function open(id: string): Promise<void> {
    currentId = id;
    const d = await GetProblem(id, lang);
    if (disposed || !d.found) return;
    currentCategory = d.category;
    closeDrawer();
    el('bTitle').textContent = d.title;
    const diff = el('bDiff');
    diff.textContent = d.difficulty;
    diff.className = `tag d-${d.difficulty}`;
    el('bCat').textContent = d.category;
    el('bPrompt').textContent = d.prompt;
    editor.setValue(d.starter);
    const model = editor.getModel();
    if (model) monaco.editor.setModelLanguage(model, monacoLang(lang));
    if (toolchainOK) {
      el('bVerdict').innerHTML = '<div class="vmuted">Write your solution and press Grade.</div>';
    } else {
      showGateBanner();
    }
    renderList();
  }

  // ── Learn drawer (category primer) ──────────────────────────
  function closeDrawer(): void {
    el('bDrawer').classList.add('hidden');
  }

  async function openDrawer(): Promise<void> {
    if (!currentCategory) return;
    const pv = await GetPrimer(currentCategory, lang);
    if (disposed) return;
    const d = el('bDrawer');
    if (!pv.found) {
      d.innerHTML = `<div class="drawhead"><span>Learn</span><button class="btn small" id="bDrawX">✕</button></div>
        <div class="vmuted" style="padding:12px">No primer for ${esc(currentCategory)} yet.</div>`;
    } else {
      const sections = (pv.sections || [])
        .map(
          (s) => `<div class="psec"><div class="psectitle">${esc(s.title)}</div>
          ${(s.ops || [])
            .map(
              (op) => `<div class="pop"><div class="poplabel">${esc(op.label)}</div><pre class="popcode">${esc(op.code.replace(/\s+$/, ''))}</pre></div>`,
            )
            .join('')}</div>`,
        )
        .join('');
      d.innerHTML = `
        <div class="drawhead"><span>${esc(pv.title)}</span><button class="btn small" id="bDrawX">✕</button></div>
        <div class="drawbody">
          <div class="psummary">${esc(pv.summary)}</div>
          ${sections}
          ${pv.example ? `<div class="psectitle">Worked example</div><pre class="popcode">${esc(pv.example.replace(/\s+$/, ''))}</pre>` : ''}
        </div>`;
    }
    d.classList.remove('hidden');
    (root.querySelector('#bDrawX') as HTMLButtonElement).addEventListener('click', closeDrawer);
  }

  async function grade(): Promise<void> {
    if (!currentId || !toolchainOK) return;
    const btn = el('bGrade') as HTMLButtonElement;
    btn.disabled = true;
    btn.textContent = 'Grading…';
    el('bVerdict').innerHTML = '<div class="vmuted">Compiling + running hidden tests…</div>';
    try {
      const res = await Grade(lang, currentId, editor.getValue());
      if (disposed) return;
      if (res.passed) {
        const item = problems.find((p) => p.id === currentId);
        if (item && !item.solved) {
          item.solved = true;
          renderList();
        }
        void refreshScore();
        const bankLine = res.newlyBanked
          ? '<div class="vhead vok">✓ Banked</div>'
          : '<div class="vhead vok">✓ Already banked</div>';
        const saveWarn = res.saveErr
          ? `<div class="cdetail">⚠ progress not saved: ${esc(res.saveErr)}</div>`
          : '';
        el('bVerdict').innerHTML =
          bankLine + saveWarn + renderVerdict(res) +
          `<div class="workbar"><button id="bNext" class="btn primary">Next problem →</button></div>`;
        (el('bNext') as HTMLButtonElement).addEventListener('click', async () => {
          const next = await NextProblem(lang, currentId);
          if (!disposed && next) void open(next);
        });
      } else {
        el('bVerdict').innerHTML = renderVerdict(res);
      }
    } catch (e) {
      if (disposed) return;
      el('bVerdict').innerHTML = `<div class="vhead vfail">✖ Grader error</div><pre class="vpre">${esc(String(e))}</pre>`;
    } finally {
      if (!disposed) {
        btn.disabled = false;
        btn.textContent = '▶ Grade';
      }
    }
  }

  el('bGrade').addEventListener('click', () => void grade());
  el('bReset').addEventListener('click', () => {
    if (currentId) void open(currentId);
  });
  el('bLearn').addEventListener('click', () => void openDrawer());
  const search = el('bSearch') as HTMLInputElement;
  search.addEventListener('input', () => renderList());
  el('bLists')
    .querySelectorAll<HTMLButtonElement>('.fchip')
    .forEach((c) =>
      c.addEventListener('click', () => {
        listFilter = c.dataset.list as string;
        el('bLists')
          .querySelectorAll<HTMLButtonElement>('.fchip')
          .forEach((x) => x.classList.toggle('on', x === c));
        renderList();
      }),
    );
  el('bDiffs')
    .querySelectorAll<HTMLButtonElement>('.fchip')
    .forEach((c) =>
      c.addEventListener('click', () => {
        const d = c.dataset.diff as string;
        diffOn[d] = !diffOn[d];
        c.classList.toggle('on', diffOn[d]);
        renderList();
      }),
    );

  void refreshScore();
  // authoritative toolchain check (cached after the first run per session)
  void (async () => {
    const st = await CheckLang(lang);
    if (!disposed) applyGate(st);
  })();
  void (async () => {
    problems = await ListProblems(lang);
    if (disposed) return;
    if (!problems || !problems.length) {
      el('bTitle').textContent = 'No problems loaded';
      el('bVerdict').innerHTML =
        '<div class="vmuted">The content catalog failed to load — restart the app.</div>';
      problems = [];
      return;
    }
    renderList();
    void open(problems[0].id);
  })();

  return () => {
    disposed = true;
    editor.dispose();
  };
}
