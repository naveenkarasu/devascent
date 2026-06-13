import { monaco, makeEditor } from './monaco-setup';
import {
  StartOrientation,
  SubmitOrientationCode,
  SubmitOrientationChoice,
  SubmitOrientationSpec,
  AdvanceOrientation,
  CheckLang,
  RecheckLang,
  GetInstallGuide,
} from '../wailsjs/go/main/App';
import { guiapi } from '../wailsjs/go/models';
import { esc, renderVerdict, renderInstallGuide } from './util';

function nav(view: string): void {
  window.dispatchEvent(new CustomEvent('devascent-nav', { detail: view }));
}

// mountOrientation renders the entrance test for lang and returns a disposer.
export function mountOrientation(root: HTMLElement, lang: string): () => void {
  let editor: monaco.editor.IStandaloneCodeEditor | null = null;
  let disposed = false; // set by the disposer; async continuations must bail
  const disposeEditor = () => {
    if (editor) {
      editor.dispose();
      editor = null;
    }
  };

  function intro(): void {
    disposeEditor();
    root.innerHTML = `
      <div class="centerpane">
        <div class="card-lg">
          <h2 class="h2">Entrance Test — ${esc(lang)}</h2>
          <p class="lead">A short adaptive check — write a little code, read an error, explain a spec in your own words — to place you at the right starting point. Be honest; it only routes you.</p>
          <p class="ask">How much have you coded before?</p>
          <div class="levelrow">
            <button class="btn" data-lvl="never">Never / barely</button>
            <button class="btn" data-lvl="a-little">A little</button>
            <button class="btn primary" data-lvl="regularly">Regularly</button>
          </div>
        </div>
      </div>`;
    root
      .querySelectorAll<HTMLButtonElement>('[data-lvl]')
      .forEach((b) => b.addEventListener('click', () => void begin(b.dataset.lvl as string)));
    // The test grades real code in lang's toolchain — gate entry if it's absent
    // (missing langs report instantly; the canary cost only applies when present).
    void (async () => {
      const st = await CheckLang(lang);
      if (!disposed && st.status !== 'available') void renderGate(st);
    })();
  }

  async function renderGate(st: guiapi.LangStatus): Promise<void> {
    const g = await GetInstallGuide(lang);
    if (disposed) return;
    const reason = st.reason || `${lang} toolchain ${st.status}`;
    root.innerHTML = `
      <div class="centerpane">
        <div class="card-lg">
          <h2 class="h2">Entrance Test — ${esc(lang)}</h2>
          <p class="lead">The entrance test grades real ${esc(lang)} code in your own toolchain, and it isn't installed yet. Install it, then re-check — or browse Tutorial Island and the primers meanwhile (reading needs no toolchain).</p>
          ${renderInstallGuide(g, reason)}
          <div class="levelrow">
            <button class="btn primary" id="oRecheck">&#8634; Re-check</button>
          </div>
          <div id="oGateMsg" class="vmuted"></div>
        </div>
      </div>`;
    (root.querySelector('#oRecheck') as HTMLButtonElement).addEventListener('click', async () => {
      const msg = root.querySelector('#oGateMsg') as HTMLElement;
      msg.textContent = 'Re-checking (compiles + runs a canary)…';
      const st2 = await RecheckLang(lang);
      if (disposed) return;
      if (st2.status === 'available') intro();
      else msg.textContent = `Still ${st2.status}: ${st2.reason || 'not detected'}`;
    });
  }

  async function begin(level: string): Promise<void> {
    const step = await StartOrientation(lang, level);
    if (disposed) return;
    renderStep(step);
  }

  function renderStep(step: guiapi.OrientationStep): void {
    disposeEditor();
    if (step.done) {
      renderResults(step);
      return;
    }
    let body = '';
    if (step.kind === 'code') {
      body = `<div id="oEditor" class="editor tall"></div>`;
    } else if (step.kind === 'choice') {
      body = `<div class="choices">${step.choices
        .map(
          (c, i) =>
            `<label class="choice"><input type="radio" name="oc" value="${i}" /><span>${esc(c)}</span></label>`,
        )
        .join('')}</div>`;
    } else {
      body = `<textarea id="oSpec" class="specinput" placeholder="Explain in your own words…"></textarea>`;
    }
    root.innerHTML = `
      <div class="orient">
        <div class="oprogress">Item ${step.index} of ${step.total} · <span class="mut">${esc(step.measures)}</span></div>
        <div class="probhead"><span class="ptitle-h">${step.kind === 'code' ? 'Write the function' : 'Question'}</span></div>
        <div class="prompt-text">${esc(step.prompt)}</div>
        ${body}
        <div class="workbar">
          <button id="oSubmit" class="btn primary">${step.kind === 'code' ? '&#9654; Run' : 'Submit'}</button>
          ${step.kind === 'code' ? '<button id="oAdvance" class="btn hidden"></button>' : ''}
        </div>
        <div id="oOutcome" class="verdict"><div class="vmuted">${step.kind === 'code' ? 'Write the function, then Run. Fix and re-run as many times as you need.' : 'Answer, then submit.'}</div></div>
      </div>`;
    if (step.kind === 'code') {
      editor = makeEditor(root.querySelector('#oEditor') as HTMLElement, lang, step.starter);
    }
    (root.querySelector('#oSubmit') as HTMLButtonElement).addEventListener('click', () => void submit(step));
  }

  async function submit(step: guiapi.OrientationStep): Promise<void> {
    const out = root.querySelector('#oOutcome') as HTMLElement;
    const btn = root.querySelector('#oSubmit') as HTMLButtonElement;
    btn.disabled = true;
    // CODE items grade WITHOUT advancing — a fail keeps you on the item to edit
    // and run again; only Continue (pass) or Skip commits the result.
    if (step.kind === 'code') {
      if (!editor) {
        btn.disabled = false;
        return;
      }
      out.innerHTML = '<div class="vmuted">Grading…</div>';
      btn.textContent = 'Running…';
      const outcome = await SubmitOrientationCode(editor.getValue());
      if (disposed) return;
      setCodeActions(outcome.passed);
      renderCodeVerdict(outcome, out);
      return;
    }
    // CHOICE / SPEC are one-shot: the backend advances, the player moves on.
    let outcome: guiapi.DiagOutcome;
    if (step.kind === 'choice') {
      const sel = root.querySelector('input[name="oc"]:checked') as HTMLInputElement | null;
      if (!sel) {
        out.innerHTML = '<div class="vmuted">Pick an option first.</div>';
        btn.disabled = false;
        return;
      }
      outcome = await SubmitOrientationChoice(parseInt(sel.value, 10));
    } else {
      outcome = await SubmitOrientationSpec((root.querySelector('#oSpec') as HTMLTextAreaElement).value);
    }
    if (disposed) return;
    btn.textContent = 'Submitted';
    renderOutcome(outcome, out);
  }

  // The code item is retry-able, and the next step lives in the ACTION BAR
  // (next to Run) rather than buried inside the results box: a pass turns Run
  // into a quiet "✓ Passed" and lights up "Continue →"; a fail keeps Run live as
  // "▶ Run again" with a quieter "Skip this one →" beside it. Both advance.
  function setCodeActions(passed: boolean): void {
    const run = root.querySelector('#oSubmit') as HTMLButtonElement | null;
    const adv = root.querySelector('#oAdvance') as HTMLButtonElement | null;
    if (!run || !adv) return;
    adv.classList.remove('hidden');
    adv.onclick = () => void advance();
    if (passed) {
      run.disabled = true;
      run.textContent = '✓ Passed';
      run.classList.remove('primary');
      adv.textContent = 'Continue →';
      adv.classList.add('primary');
    } else {
      run.disabled = false;
      run.textContent = '▶ Run again';
      run.classList.add('primary');
      adv.textContent = 'Skip this one →';
      adv.classList.remove('primary');
    }
  }

  // renderCodeVerdict fills the results box with the status + hidden-test detail
  // ONLY — the next-step buttons live in the action bar (see setCodeActions).
  function renderCodeVerdict(outcome: guiapi.DiagOutcome, out: HTMLElement): void {
    if (outcome.verdict && outcome.verdict.results && outcome.verdict.results.length) {
      out.innerHTML = renderVerdict(outcome.verdict);
    } else if (outcome.feedback) {
      const head = outcome.passed
        ? '<div class="vhead vok">✓ Passed</div>'
        : '<div class="vhead vfail">✖ Not quite</div>';
      out.innerHTML = head + `<div class="cdetail">${esc(outcome.feedback)}</div>`;
    } else {
      out.innerHTML = outcome.passed
        ? '<div class="vhead vok">✓ Passed</div>'
        : '<div class="vhead vfail">✖ Tests failed — edit your code and Run again</div>';
    }
  }

  async function advance(): Promise<void> {
    const outcome = await AdvanceOrientation();
    if (disposed) return;
    renderStep(outcome.next);
  }

  function renderOutcome(outcome: guiapi.DiagOutcome, out: HTMLElement): void {
    const head = outcome.passed
      ? '<div class="vhead vok">✓ Correct</div>'
      : '<div class="vhead vfail">✖ Not quite</div>';
    let detail = '';
    if (outcome.verdict && outcome.verdict.results && outcome.verdict.results.length) {
      detail = renderVerdict(outcome.verdict);
    } else if (outcome.feedback) {
      detail = `<div class="cdetail">${esc(outcome.feedback)}</div>`;
    }
    out.innerHTML = head + detail + `<div class="workbar"><button id="oNext" class="btn primary">Next →</button></div>`;
    (root.querySelector('#oNext') as HTMLButtonElement).addEventListener('click', () => renderStep(outcome.next));
  }

  function renderResults(step: guiapi.OrientationStep): void {
    // Placement-aware next actions — don't send an acer to Tutorial Island.
    const routes: Record<string, { msg: string; buttons: string }> = {
      'test-out': {
        msg: 'You aced it — you can jump straight to the bench.',
        buttons: `<button class="btn primary" data-nav="bench">Go to the Bench →</button>`,
      },
      'dev-literacy': {
        msg: 'Solid coding. A short dev-literacy brush-up (terminal, errors, files) is recommended before the bench.',
        buttons: `<button class="btn primary" data-nav="devlit">Dev-Literacy brush-up →</button>
                  <button class="btn" data-nav="bench">Skip to the Bench →</button>`,
      },
      'tutorial-full': {
        msg: 'Great start — Tutorial Island will get you bench-ready, step by step.',
        buttons: `<button class="btn primary" data-nav="tutorial">Go to Tutorial Island →</button>`,
      },
    };
    const r = routes[step.placement] ?? routes['tutorial-full'];
    const sig = (label: string, ok: number, total: number) =>
      total > 0 ? `<span class="sigchip">${label} ${ok}/${total}</span>` : '';
    root.innerHTML = `
      <div class="centerpane">
        <div class="card-lg">
          <h2 class="h2">Placement</h2>
          <p class="bigscore">${step.score} / ${step.total}</p>
          <div class="sigrow">
            ${sig('coding', step.codingOK, step.codingTotal)}
            ${sig('machine', step.machineOK, step.machineTotal)}
            ${sig('spec', step.specOK, step.specTotal)}
          </div>
          <p class="lead">${esc(r.msg)}</p>
          <div class="levelrow">
            <button class="btn" id="oRedo">↺ Retake</button>
            ${r.buttons}
          </div>
        </div>
      </div>`;
    (root.querySelector('#oRedo') as HTMLButtonElement).addEventListener('click', () => intro());
    root
      .querySelectorAll<HTMLButtonElement>('[data-nav]')
      .forEach((b) => b.addEventListener('click', () => nav(b.dataset.nav as string)));
  }

  intro();
  return () => {
    disposed = true;
    disposeEditor();
  };
}
