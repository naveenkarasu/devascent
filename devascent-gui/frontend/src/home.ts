import {
  GetProgress,
  GetProfiles,
  DeleteProfile,
  GetMentorBackends,
  SelectMentor,
  SetMentorEndpoint,
} from '../wailsjs/go/main/App';
import { guiapi } from '../wailsjs/go/models';
import { esc } from './util';
import { langIcon } from './langicons';
import appIcon from './assets/appicon.png'; // the APP logo (channel logo is splash-only)

function nav(view: string): void {
  window.dispatchEvent(new CustomEvent('devascent-nav', { detail: view }));
}

function switchLang(lang: string): void {
  window.dispatchEvent(new CustomEvent('devascent-lang', { detail: lang }));
}

// continueTarget smart-routes the Continue button from the saved run state.
function continueTarget(pr: guiapi.Progress): { view: string; label: string } {
  if (!pr.placement) return { view: 'orientation', label: 'Continue — take the Entrance Test' };
  if (pr.step0Met || pr.placement === 'test-out')
    return { view: 'bench', label: 'Continue — back to the Bench' };
  if (pr.placement === 'dev-literacy')
    return { view: 'devlit', label: 'Continue — Dev-Literacy brush-up' };
  return { view: 'tutorial', label: 'Continue — Tutorial Island' };
}

// mountHome renders the front door: brand + narrative left, the live run card
// right (the approved start-screen layout, now data-driven from the real save).
export function mountHome(root: HTMLElement, lang: string): () => void {
  let disposed = false;

  void (async () => {
    const [pr, profiles] = await Promise.all([GetProgress(lang), GetProfiles()]);
    if (disposed) return;
    const cont = continueTarget(pr);
    // Profile picker: one chip per existing language slot; the active language
    // is highlighted, clicking another switches the whole app to that slot.
    const profileRow = (profiles || []).length
      ? `<div class="profrow">
          <span class="proflabel">profiles</span>
          ${(profiles || [])
            .map(
              (p) => `<span class="profpair"><button class="profchip ${p.lang === lang ? 'on' : ''}" data-lang="${esc(p.lang)}">
                <span class="b">${esc(p.lang)}</span><span class="mut"> · ${p.banked} banked${p.placement ? ' · ' + esc(p.placement) : ''}</span>
              </button><button class="profx" data-del="${esc(p.lang)}" title="Delete this profile">✕</button></span>`,
            )
            .join('')}
        </div>`
      : '';
    const placeChip = pr.placement
      ? `<span class="chip">&#9670; placement · ${esc(pr.placement)}</span>`
      : `<span class="chip">&#9670; entrance test not taken</span>`;
    const milestone = pr.step0Met
      ? `<div class="g">&#10003; Step 0 milestone met &nbsp;<span class="mut">(${esc(pr.track)})</span></div>`
      : `<div class="mut">step 0 &middot; bank ${pr.bankTarget} &middot; ${pr.catTarget} categories &middot; ${pr.hardTarget} hard</div>`;
    root.innerHTML = `
      <div class="main startbody">
        <div class="info">
          <div class="brandrow"><img class="brandlogo" src="${appIcon}" alt="" draggable="false"><div class="appname">DevAscent</div></div>
          <div class="subtitle">Grow from junior to senior by doing the real work — graded by a real compiler.</div>
          <div class="narrative">A run flows through <b>Orientation</b> (an adaptive entrance test → tutorial or test-out) into <b>The Apprenticeship</b> — a curated DSA bench where you write real functions, graded by real hidden tests in your own language toolchain.</div>
          <div class="savechips">
            ${placeChip}
            <span class="chip" style="color:var(--green)">&#10003; banked ${pr.banked}</span>
            <span class="chip">${langIcon(lang)} ${esc(lang)}</span>
          </div>
          ${profileRow}
          <div class="ver muted mono">graded by your local toolchain</div>
        </div>
        <div class="rightcol">
          <div class="card">
            <div class="cardhead"><div style="font-weight:600">Your run</div>
              <span class="chip"><span class="dot" style="background:var(--green)"></span> apprentice</span></div>
            <div class="term">
              <div><span class="mut">~/bench ❯</span> <span class="b">progress</span></div>
              <div class="g">&#10003; banked ${pr.banked}/${pr.bankTarget} &nbsp;<span class="mut">cats ${pr.cats}/${pr.catTarget} · hard ${pr.hard}/${pr.hardTarget}</span></div>
              ${milestone}
              <div style="margin-top:8px"><span class="mut">~/orientation ❯</span> status</div>
              <div><span class="y">&#9670;</span> ${pr.placement ? `placed: ${esc(pr.placement)} <span class="mut">(${esc(pr.level)})</span>` : 'not taken yet'}</div>
              <div style="margin-top:8px"><span class="mut">❯</span> <span class="caret"></span></div>
            </div>
            <div class="actions">
              <button class="btn primary" data-nav="${cont.view}">&#9654;&nbsp; ${esc(cont.label)}</button>
              <button class="btn" data-nav="orientation">&#9670;&nbsp; Entrance Test</button>
              <button class="btn" data-nav="bench">&#9646;&nbsp; Browse the Bench</button>
            </div>
          </div>
          <div class="card mentorcard">
            <div class="cardhead"><div style="font-weight:600">Mentor</div>
              <span class="chip" id="mState">…</span></div>
            <div class="mentornote">Paid hints can be answered by an AI you already have — Claude Code, Codex, Copilot, or a local model via Ollama/LM&nbsp;Studio. The game ships none, stores no keys, and falls back to the built-in playbook offline.</div>
            <div id="mRows"></div>
          </div>
        </div>
      </div>`;
    root
      .querySelectorAll<HTMLButtonElement>('[data-nav]')
      .forEach((b) => b.addEventListener('click', () => nav(b.dataset.nav as string)));
    root
      .querySelectorAll<HTMLButtonElement>('.profchip')
      .forEach((b) => b.addEventListener('click', () => switchLang(b.dataset.lang as string)));
    // delete profile: two-click armed confirm (✕ → "sure?" → gone); disarms after 3s
    root.querySelectorAll<HTMLButtonElement>('.profx').forEach((b) =>
      b.addEventListener('click', async () => {
        if (!b.dataset.armed) {
          b.dataset.armed = '1';
          b.textContent = 'sure?';
          b.classList.add('armed');
          setTimeout(() => {
            if (disposed) return;
            delete b.dataset.armed;
            b.textContent = '✕';
            b.classList.remove('armed');
          }, 3000);
          return;
        }
        const err = await DeleteProfile(b.dataset.del as string);
        if (disposed) return;
        if (err) {
          b.textContent = '!';
          b.title = err;
          return;
        }
        nav('home'); // re-render the view with the slot gone
      }),
    );
    void renderMentor();
  })();

  // renderMentor (re)draws the BYO-AI backend picker rows; re-run after any
  // selection so the "in use" mark and probe states stay live.
  async function renderMentor(): Promise<void> {
    const rows = await GetMentorBackends();
    if (disposed) return;
    const host = root.querySelector('#mRows') as HTMLElement | null;
    if (!host) return;
    const st = root.querySelector('#mState');
    if (st) st.textContent = (rows || []).find((r) => r.selected)?.name || 'Built-in playbook';
    host.innerHTML = (rows || [])
      .map((s) => {
        const probe = s.probed
          ? s.probeOk
            ? '<span class="mprobe ok">probe ✓</span>'
            : '<span class="mprobe bad">probe ✖</span>'
          : '';
        const action = s.selected
          ? '<span class="msel">✓ in use</span>'
          : s.present
            ? `<button class="btn small muse" data-id="${esc(s.id)}">Use</button>`
            : '';
        const cfg =
          s.id === 'openai-compat'
            ? '<button class="btn small" id="mCfgBtn" title="Configure endpoint, model and key">⚙</button>'
            : '';
        const form =
          s.id === 'openai-compat'
            ? `<div class="mform hidden" id="mForm">
                <input id="mEndpoint" class="minput" placeholder="endpoint URL — e.g. http://localhost:11434/v1">
                <input id="mModel" class="minput" placeholder="model — e.g. llama3.1">
                <input id="mKey" class="minput" type="password" placeholder="API key (optional)">
                <div><button class="btn small primary" id="mSave">Save &amp; use</button></div>
              </div>`
            : '';
        return `<div class="mrow">
            <span class="mdot ${s.present ? 'on' : ''}"></span>
            <span class="mname">${esc(s.name)}</span>
            <span class="minfo">${esc(s.info)}</span>
            ${probe}${cfg}${action}
          </div>
          ${form}
          <div class="merr ${s.probeErr ? '' : 'hidden'}" data-merr="${esc(s.id)}">${esc(s.probeErr)}</div>`;
      })
      .join('');
    const showErr = (id: string, err: string): void => {
      const e = host.querySelector(`[data-merr="${id}"]`) as HTMLElement | null;
      if (!e) return;
      e.textContent = err;
      e.classList.remove('hidden');
    };
    // a probe can take a while — lock every probe-triggering button while one runs
    const setProbing = (on: boolean): void =>
      host
        .querySelectorAll<HTMLButtonElement>('.muse, #mSave')
        .forEach((x) => (x.disabled = on));
    host.querySelectorAll<HTMLButtonElement>('.muse').forEach((b) =>
      b.addEventListener('click', async () => {
        setProbing(true);
        b.textContent = 'probing…';
        const err = await SelectMentor(b.dataset.id as string);
        if (disposed) return;
        if (err) {
          setProbing(false);
          b.textContent = 'Use';
          showErr(b.dataset.id as string, err);
          return;
        }
        void renderMentor();
      }),
    );
    host
      .querySelector('#mCfgBtn')
      ?.addEventListener('click', () => host.querySelector('#mForm')?.classList.toggle('hidden'));
    const save = host.querySelector('#mSave') as HTMLButtonElement | null;
    save?.addEventListener('click', async () => {
      const val = (id: string) => (host.querySelector('#' + id) as HTMLInputElement).value.trim();
      setProbing(true);
      save.textContent = 'probing…';
      let err = await SetMentorEndpoint(val('mEndpoint'), val('mModel'), val('mKey'));
      if (!disposed && !err) err = await SelectMentor('openai-compat');
      if (disposed) return;
      setProbing(false);
      save.textContent = 'Save & use';
      if (err) {
        showErr('openai-compat', err);
        return;
      }
      void renderMentor();
    });
  }

  return () => {
    disposed = true;
  };
}
