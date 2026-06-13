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
  GetWallet,
  RequestHint,
  GetWriteup,
  SubmitWriteup,
  GetGate,
  GetMentorPreview,
} from '../wailsjs/go/main/App';
import { guiapi } from '../wailsjs/go/models';
import { esc, renderVerdict, renderInstallGuide } from './util';

// mountBench renders the bench workbench into root for lang and returns a
// disposer (call before unmounting to release the Monaco editor).
// graded=false marks a reference-only language (e.g. C++): browse/read works
// as-is, but Grade is disabled and the hint economy is hidden.
export function mountBench(root: HTMLElement, lang: string, graded = true): () => void {
  root.innerHTML = `
    <div class="bench">
      <aside class="sidebar">
        <div id="bScore" class="scorecard"></div>
        <div id="bGate" class="gatecard"></div>
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
          <button id="bHintBtn" class="btn small">✦ Hint</button>
          <button id="bLearn" class="btn small">&#9656; Learn</button>
        </div>
        <div id="bDrawer" class="drawer hidden"></div>
        <div id="bHintPanel" class="drawer hintpanel hidden"></div>
        <div id="bPrompt" class="prompt-text"></div>
        <div id="bEditor" class="editor"></div>
        <div class="workbar">
          <button id="bGrade" class="btn primary">&#9654; Grade</button>
          <button id="bReset" class="btn">&#8634; Reset starter</button>
        </div>
        <div id="bVerdict" class="verdict"><div class="vmuted">Pick a problem to begin.</div></div>
      </main>
      <div id="bModal" class="modal hidden"></div>
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

  // ── Reference-only language (no grader wired) ────────────────
  const refLabel = lang === 'cpp' ? 'C++' : lang;
  const refNote = `<div class="vmuted">Reference mode — browse the prompt and starter; grading for ${esc(refLabel)} isn't wired yet.</div>`;
  if (!graded) {
    const g = el('bGrade') as HTMLButtonElement;
    g.disabled = true;
    g.title = `${refLabel} is reference-only — its grader isn't wired yet`;
    el('bHintBtn').classList.add('hidden'); // no hint economy without a grader
  }

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
    closeHintPanel(); // the hint panel shares the right edge
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
    const [pr, gate] = await Promise.all([GetProgress(lang), GetGate(lang)]);
    if (disposed) return;
    renderGate(gate);
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

  // ── Graduation gate (Blind 75) section ──────────────────────
  let gateOpen = false; // expanded state survives re-renders within this mount

  function renderGate(g: guiapi.GateView): void {
    const pend = g.provisional
      ? ` <span class="scmut">(+${g.provisional} pending write-up)</span>`
      : '';
    const cats = (g.categories || [])
      .map((c) => {
        const ok = c.done >= c.required;
        const pct = Math.min(100, Math.round((100 * c.done) / Math.max(1, c.required)));
        return `<div class="gcat"><span class="gcname">${esc(c.category)}</span><span class="gcnum ${ok ? 'ok' : ''}">${c.done}/${c.required}</span>
          <div class="gcbar"><div class="gcfill ${ok ? 'ok' : ''}" style="width:${pct}%"></div></div></div>`;
      })
      .join('');
    const mand = (g.mandatory || [])
      .map((m) => `<div class="gman ${m.done ? 'ok' : ''}">${m.done ? '✓' : '○'} ${esc(m.title)}</div>`)
      .join('');
    el('bGate').innerHTML = `
      <div class="gatehead" id="bGateHead">${gateOpen ? '▾' : '▸'} Graduation gate — Blind 75</div>
      ${g.met ? '<div class="gatemet">GATE MET — apprenticeship complete</div>' : ''}
      <div class="gatecount">${g.full} / ${g.target} fully banked${pend}</div>
      <div class="gatebody ${gateOpen ? '' : 'hidden'}">${cats}<div class="gmanhead">Mandatory</div>${mand}</div>`;
    el('bGateHead').addEventListener('click', () => {
      gateOpen = !gateOpen;
      renderGate(g);
    });
  }

  const collapsed = new Set<string>(); // collapsed category headers

  function renderList(): void {
    const items = visible();
    const banked = items.filter((p) => p.solved).length;
    let html = '';
    let lastCat = ' ';
    for (const p of items) {
      const cat = p.category || 'Other';
      if (cat !== lastCat) {
        lastCat = cat;
        const n = items.filter((x) => (x.category || 'Other') === cat).length;
        const done = items.filter((x) => (x.category || 'Other') === cat && x.solved).length;
        html += `<div class="pcat" data-cat="${esc(cat)}">${collapsed.has(cat) ? '▸' : '▾'} ${esc(cat)} <span class="ccount">${done}/${n}</span></div>`;
      }
      if (collapsed.has(cat)) continue;
      const prov = p.solved && !p.writeup; // passed but write-up pending
      html += `<div class="prow ${p.id === currentId ? 'active' : ''}" data-id="${p.id}">
        <span class="ptick ${p.solved ? (prov ? 'prov' : 'on') : ''}"${prov ? ` title="passed — write-up pending" data-wid="${p.id}"` : ''}>${p.solved ? (prov ? '◐' : '✓') : ''}</span>
        <span class="ptitle">${esc(p.title)}</span>
        ${prov ? `<span class="pwrite" data-wid="${p.id}" title="finish the write-up (+1 ⬡)">✍</span>` : ''}
        <span class="pdiff d-${esc(p.difficulty)}">${esc(p.difficulty)}</span></div>`;
    }
    el('bList').innerHTML = html;
    el('bList')
      .querySelectorAll<HTMLElement>('.prow')
      .forEach((r) => r.addEventListener('click', () => void open(r.dataset.id as string)));
    el('bList')
      .querySelectorAll<HTMLElement>('[data-wid]')
      .forEach((b) =>
        b.addEventListener('click', (ev) => {
          ev.stopPropagation();
          void openWriteup(b.dataset.wid as string);
        }),
      );
    el('bList')
      .querySelectorAll<HTMLElement>('.pcat')
      .forEach((h) =>
        h.addEventListener('click', () => {
          const c = h.dataset.cat as string;
          if (collapsed.has(c)) collapsed.delete(c);
          else collapsed.add(c);
          renderList();
        }),
      );
    el('bCount').textContent = `${items.length} problems · ${banked} banked`;
  }

  async function open(id: string): Promise<void> {
    currentId = id;
    const d = await GetProblem(id, lang);
    if (disposed || !d.found) return;
    currentCategory = d.category;
    closeDrawer();
    closeHintPanel(); // hints are per-problem; drop the stale ones
    el('bTitle').textContent = d.title;
    const diff = el('bDiff');
    diff.textContent = d.difficulty;
    diff.className = `tag d-${d.difficulty}`;
    el('bCat').textContent = d.category;
    el('bPrompt').textContent = d.prompt;
    editor.setValue(d.starter);
    const model = editor.getModel();
    if (model) monaco.editor.setModelLanguage(model, monacoLang(lang));
    if (!graded) {
      el('bVerdict').innerHTML = refNote;
    } else if (toolchainOK) {
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
    closeHintPanel(); // the hint panel shares the right edge
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

  // ── Centered modal (write-up gate, mentor preview) ──────────
  function closeModal(): void {
    el('bModal').classList.add('hidden');
    el('bModal').innerHTML = '';
  }

  function showModal(html: string): HTMLElement {
    const m = el('bModal');
    m.innerHTML = `<div class="modalcard">${html}</div>`;
    m.classList.remove('hidden');
    m.querySelectorAll<HTMLElement>('[data-close]').forEach((b) => b.addEventListener('click', closeModal));
    const f =
      m.querySelector<HTMLElement>('[data-autofocus]') ||
      m.querySelector<HTMLElement>('textarea, input, button');
    f?.focus();
    return m;
  }

  // ── Hint economy (Track A2) ──────────────────────────────────
  let wallet: guiapi.WalletView | null = null;
  let walletAt = 0; // Date.now() at fetch — anchors the recharge countdown
  let walletFetching = false;
  let hintTimer: number | undefined;
  let hintBusy = false;

  function fmtClock(sec: number): string {
    return `${Math.floor(sec / 60)}:${String(sec % 60).padStart(2, '0')}`;
  }

  function renderWalletLine(): void {
    const w = root.querySelector('#hWallet') as HTMLElement | null;
    if (!w || !wallet) return;
    let line = `⬡ ${wallet.tokens} token${wallet.tokens === 1 ? '' : 's'} · ◉ ${wallet.nudgeCharges}/${wallet.nudgeMax} nudges`;
    if (wallet.nudgeCharges < wallet.nudgeMax) {
      const left = wallet.nextRechargeSec - Math.floor((Date.now() - walletAt) / 1000);
      if (left <= 0) void refreshWallet(); // a charge just recharged
      else line += ` · next in ${fmtClock(left)}`;
    }
    w.textContent = line;
  }

  async function refreshWallet(): Promise<void> {
    if (walletFetching) return;
    walletFetching = true;
    try {
      const w = await GetWallet(lang);
      if (disposed) return;
      wallet = w;
      walletAt = Date.now();
      renderWalletLine();
    } finally {
      walletFetching = false;
    }
  }

  function stopHintTick(): void {
    if (hintTimer !== undefined) {
      clearInterval(hintTimer);
      hintTimer = undefined;
    }
  }

  function closeHintPanel(): void {
    stopHintTick();
    el('bHintPanel').classList.add('hidden');
  }

  function setHintButtons(disabled: boolean): void {
    el('bHintPanel')
      .querySelectorAll<HTMLButtonElement>('[data-tier]')
      .forEach((b) => (b.disabled = disabled));
  }

  // setHintOut writes into the panel's output slot, re-queried fresh — the
  // panel may have been re-rendered while a mentor call was in flight.
  function setHintOut(cls: string, html: string): void {
    const out = root.querySelector('#hOut') as HTMLElement | null;
    if (!out) return;
    out.className = cls;
    out.innerHTML = html;
  }

  const hintCache = new Map<string, string>(); // problem ID → last paid-for hint HTML

  async function askHint(tier: number): Promise<void> {
    if (hintBusy || !currentId) return;
    const reqId = currentId; // pin the problem this hint belongs to
    hintBusy = true;
    setHintButtons(true);
    setHintOut('hout vmuted', tier === 1 ? 'thinking…' : 'mentor is thinking… (can take up to 45s)');
    try {
      const res = await RequestHint(lang, reqId, tier, editor.getValue());
      if (disposed) return;
      if (res.err) {
        void refreshWallet(); // refused requests come back with a zero-valued wallet
        if (currentId === reqId) setHintOut('hout herr', esc(res.err));
        return;
      }
      wallet = res.wallet;
      walletAt = Date.now();
      renderWalletLine();
      const notice =
        (res.pity ? '<div class="hnotice">free hint — you\'ve earned it for persistence</div>' : '') +
        (res.refunded
          ? '<div class="hnotice">AI unavailable — answered from templates, token refunded</div>'
          : '');
      const html = `${notice}<div class="htext">${esc(res.text)}</div><div class="hsrc">${esc(res.source)} · tier ${res.tier}</div>`;
      hintCache.set(reqId, html); // survives a panel toggle / problem round-trip
      if (currentId === reqId) setHintOut('hout', html);
    } finally {
      if (!disposed) {
        hintBusy = false;
        setHintButtons(false);
      }
    }
  }

  async function openPreview(): Promise<void> {
    if (!currentId) return;
    const txt = await GetMentorPreview(lang, currentId, 2, editor.getValue());
    if (disposed) return;
    showModal(`<div class="mhead">What gets sent</div>
      <div class="vmuted" style="margin-bottom:8px">The exact prompt a paid hint sends to your mentor — nothing else leaves your machine.</div>
      <pre class="vpre mpre">${esc(txt)}</pre>
      <div class="workbar"><button class="btn small" data-close>Close</button></div>`);
  }

  function openHintPanel(): void {
    if (!graded || !currentId) return;
    closeDrawer(); // the Learn drawer shares the right edge
    const d = el('bHintPanel');
    d.innerHTML = `
      <div class="drawhead"><span>✦ Hints</span><button class="btn small" id="bHintX">✕</button></div>
      <div class="drawbody">
        <div id="hWallet" class="hwallet">…</div>
        <div class="hrow"><div><div class="hname">Nudge</div><div class="hcost">free · uses a ◉ charge</div></div><button class="btn small" data-tier="1">Ask</button></div>
        <div class="hrow"><div><div class="hname">Strategy</div><div class="hcost">1 ⬡ · the approach, not the code</div></div><button class="btn small" data-tier="2">Ask</button></div>
        <div class="hrow"><div><div class="hname">Walkthrough</div><div class="hcost">3 ⬡ · step-by-step plan</div></div><button class="btn small" data-tier="3">Ask</button></div>
        <div id="hOut" class="hout vmuted">Stuck? A nudge costs nothing but a charge.</div>
        <a id="hPreview" class="hlink">what gets sent</a>
      </div>`;
    d.classList.remove('hidden');
    el('bHintX').addEventListener('click', closeHintPanel);
    el('hPreview').addEventListener('click', () => void openPreview());
    d.querySelectorAll<HTMLButtonElement>('[data-tier]').forEach((b) =>
      b.addEventListener('click', () => {
        const tier = Number(b.dataset.tier);
        if (tier === 1) {
          void askHint(1);
          return;
        }
        // paid tiers: two-click armed confirm (✕ → "sure?" pattern)
        if (!b.dataset.armed) {
          b.dataset.armed = '1';
          b.textContent = 'sure?';
          b.classList.add('armed');
          setTimeout(() => {
            if (disposed || !b.dataset.armed) return;
            delete b.dataset.armed;
            b.textContent = 'Ask';
            b.classList.remove('armed');
          }, 3000);
          return;
        }
        delete b.dataset.armed;
        b.textContent = 'Ask';
        b.classList.remove('armed');
        void askHint(tier);
      }),
    );
    if (hintBusy) {
      // a mentor call is still in flight from before the panel was toggled
      setHintButtons(true);
      setHintOut('hout vmuted', 'mentor is thinking… (can take up to 45s)');
    } else {
      const cached = hintCache.get(currentId);
      if (cached) setHintOut('hout', cached);
    }
    void refreshWallet();
    stopHintTick();
    hintTimer = window.setInterval(renderWalletLine, 1000);
  }

  // ── Write-up gate (Track A1) ─────────────────────────────────
  async function openWriteup(id: string): Promise<void> {
    const w = await GetWriteup(lang, id);
    if (disposed) return;
    if (w.done) {
      // the row was stale (write-up already banked) — sync the list quietly
      const item = problems.find((p) => p.id === id);
      if (item && !item.writeup) {
        item.writeup = true;
        renderList();
      }
      return;
    }
    if (!w.solved) return;
    const mcq = w.hasMcq
      ? `<div class="wq">${esc(w.question)}</div>
        <div class="wopts" id="wOpts">${(w.options || [])
          .map((o, i) => `<label class="wopt"><input type="radio" name="wmcq" value="${i}"><span>${esc(o)}</span></label>`)
          .join('')}</div>`
      : '';
    const m = showModal(`
      <div class="mhead">✍ ${esc(w.title)}</div>
      <div class="vmuted">Explain your solve to bank it fully (+1 ⬡).</div>
      ${mcq}
      <textarea id="wText" class="specinput wtext" data-autofocus placeholder="How did you approach it? (min ${w.minLen} chars)"></textarea>
      <div id="wMsg" class="herr hidden"></div>
      <div class="workbar">
        <button id="wSubmit" class="btn small primary">Submit</button>
        <button class="btn small" data-close title="Keeps the solve provisional — the ✍ mark in the list brings you back">Later</button>
      </div>`);
    const btn = m.querySelector('#wSubmit') as HTMLButtonElement;
    const msg = m.querySelector('#wMsg') as HTMLElement;
    const opts = m.querySelector('#wOpts') as HTMLElement | null;
    // picking a different answer clears the rejected state for a clean retry
    m.querySelectorAll<HTMLInputElement>('input[name="wmcq"]').forEach((r) =>
      r.addEventListener('change', () => {
        msg.classList.add('hidden');
        opts?.classList.remove('bad', 'shake');
      }),
    );
    btn.addEventListener('click', async () => {
      const sel = m.querySelector('input[name="wmcq"]:checked') as HTMLInputElement | null;
      if (w.hasMcq && !sel) {
        msg.textContent = 'pick an answer first';
        msg.classList.remove('hidden');
        return;
      }
      const text = (m.querySelector('#wText') as HTMLTextAreaElement).value;
      // disable Submit AND Later through the await — no double-submit, no
      // closing the modal with a grade in flight
      const buttons = Array.from(m.querySelectorAll<HTMLButtonElement>('button'));
      buttons.forEach((b) => (b.disabled = true));
      let res: guiapi.WriteupResult;
      try {
        res = await SubmitWriteup(lang, id, sel ? Number(sel.value) : 0, text);
      } finally {
        if (!disposed) buttons.forEach((b) => (b.disabled = false));
      }
      if (disposed) return;
      if (!res.accepted) {
        if (w.hasMcq && !res.mcqCorrect && opts) {
          opts.classList.remove('shake'); // restart the animation on repeat misses
          void opts.offsetWidth;
          opts.classList.add('shake', 'bad');
        }
        msg.textContent =
          res.err || (!res.mcqCorrect ? 'not quite — check the complexity again' : 'not accepted');
        msg.classList.remove('hidden');
        return;
      }
      const item = problems.find((p) => p.id === id);
      if (item) {
        item.solved = true;
        item.writeup = true;
      }
      renderList();
      void refreshScore();
      wallet = res.wallet;
      walletAt = Date.now();
      renderWalletLine();
      const card = m.querySelector('.modalcard') as HTMLElement | null;
      if (!card) return; // modal was closed (Esc) mid-flight — state is already synced
      const follow = res.followup
        ? `<div class="wfollow">Mentor asks: ${esc(res.followup)}</div><div class="vmuted" style="margin-top:6px">Something to mull over — no answer needed.</div>`
        : '';
      card.innerHTML = `
        <div class="mhead vok">✓ Banked fully — +${res.tokensAwarded} ⬡</div>
        ${follow}
        <div class="workbar"><button class="btn small primary" data-close>Close</button></div>`;
      card.querySelectorAll<HTMLElement>('[data-close]').forEach((b) => b.addEventListener('click', closeModal));
      (card.querySelector('[data-close]') as HTMLElement | null)?.focus();
    });
  }

  async function grade(): Promise<void> {
    if (!graded || !currentId || !toolchainOK) return;
    const btn = el('bGrade') as HTMLButtonElement;
    btn.disabled = true;
    btn.textContent = 'Grading…';
    el('bVerdict').innerHTML = '<div class="vmuted">Compiling + running hidden tests…</div>';
    try {
      const res = await Grade(lang, currentId, editor.getValue());
      if (disposed) return;
      if (res.passed) {
        const item = problems.find((p) => p.id === currentId);
        if (item) {
          item.solved = true;
          item.writeup = !res.writeupPending;
          renderList();
        }
        void refreshScore();
        if (res.tokensAwarded > 0) void refreshWallet(); // clean-solve payout → hint wallet is stale
        const bankLine = res.newlyBanked
          ? '<div class="vhead vok">✓ Banked</div>'
          : '<div class="vhead vok">✓ Already banked</div>';
        const tokLine =
          res.tokensAwarded > 0 ? `<div class="vtok">+${res.tokensAwarded} ⬡ clean solve</div>` : '';
        const saveWarn = res.saveErr
          ? `<div class="cdetail">⚠ progress not saved: ${esc(res.saveErr)}</div>`
          : '';
        const writeBtn = res.writeupPending
          ? '<button id="bWriteup" class="btn primary">✍ Explain it to bank fully (+1 ⬡)</button>'
          : '';
        el('bVerdict').innerHTML =
          bankLine + tokLine + saveWarn + renderVerdict(res) +
          `<div class="workbar">${writeBtn}<button id="bNext" class="btn ${res.writeupPending ? '' : 'primary'}">Next problem →</button></div>`;
        if (res.writeupPending) {
          const wid = currentId;
          (el('bWriteup') as HTMLButtonElement).addEventListener('click', () => void openWriteup(wid));
        }
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

  // Esc closes the topmost layer: modal → hint panel → Learn drawer.
  function onKeydown(ev: KeyboardEvent): void {
    if (ev.key !== 'Escape') return;
    if (!el('bModal').classList.contains('hidden')) closeModal();
    else if (!el('bHintPanel').classList.contains('hidden')) closeHintPanel();
    else if (!el('bDrawer').classList.contains('hidden')) closeDrawer();
  }
  document.addEventListener('keydown', onKeydown);

  el('bGrade').addEventListener('click', () => void grade());
  el('bReset').addEventListener('click', () => {
    if (currentId) void open(currentId);
  });
  el('bLearn').addEventListener('click', () => void openDrawer());
  el('bHintBtn').addEventListener('click', () => {
    if (el('bHintPanel').classList.contains('hidden')) openHintPanel();
    else closeHintPanel();
  });
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
  // authoritative toolchain check (cached after the first run per session) —
  // pointless for a reference-only language, whose Grade stays disabled anyway
  if (graded) {
    void (async () => {
      const st = await CheckLang(lang);
      if (!disposed) applyGate(st);
    })();
  }
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
    stopHintTick();
    closeModal(); // drop modal content now; the router wipes the root next
    document.removeEventListener('keydown', onKeydown);
    editor.dispose();
  };
}
