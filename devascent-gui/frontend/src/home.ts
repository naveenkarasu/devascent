import { GetProgress, GetProfiles } from '../wailsjs/go/main/App';
import { guiapi } from '../wailsjs/go/models';
import { esc } from './util';

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
              (p) => `<button class="profchip ${p.lang === lang ? 'on' : ''}" data-lang="${esc(p.lang)}">
                <span class="b">${esc(p.lang)}</span><span class="mut"> · ${p.banked} banked${p.placement ? ' · ' + esc(p.placement) : ''}</span>
              </button>`,
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
          <div class="brandrow"><div class="brandmark"></div><div class="appname">DevAscent</div></div>
          <div class="subtitle">Grow from junior to senior by doing the real work — graded by a real compiler.</div>
          <div class="narrative">A run flows through <b>Orientation</b> (an adaptive entrance test → tutorial or test-out) into <b>The Apprenticeship</b> — a curated DSA bench where you write real functions, graded by real hidden tests in your own language toolchain.</div>
          <div class="savechips">
            ${placeChip}
            <span class="chip" style="color:var(--green)">&#10003; banked ${pr.banked}</span>
            <span class="chip">&#955; ${esc(lang)}</span>
          </div>
          ${profileRow}
          <div class="ver muted mono">graded by your local toolchain</div>
        </div>
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
      </div>`;
    root
      .querySelectorAll<HTMLButtonElement>('[data-nav]')
      .forEach((b) => b.addEventListener('click', () => nav(b.dataset.nav as string)));
    root
      .querySelectorAll<HTMLButtonElement>('.profchip')
      .forEach((b) => b.addEventListener('click', () => switchLang(b.dataset.lang as string)));
  })();

  return () => {
    disposed = true;
  };
}
