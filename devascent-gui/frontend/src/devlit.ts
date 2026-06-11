import { StartDevLiteracy, SubmitDevCommand } from '../wailsjs/go/main/App';
import { guiapi } from '../wailsjs/go/models';
import { esc } from './util';

function nav(view: string): void {
  window.dispatchEvent(new CustomEvent('devascent-nav', { detail: view }));
}

// mountDevLit renders the dev-literacy track (a command checker styled as a
// terminal transcript — not a real shell) and returns a disposer.
export function mountDevLit(root: HTMLElement): () => void {
  let disposed = false;
  let transcript: string[] = [];

  async function start(): Promise<void> {
    const step = await StartDevLiteracy();
    if (disposed) return;
    transcript = [];
    renderStep(step);
  }

  function renderStep(step: guiapi.DevLitStep): void {
    if (step.done) {
      renderResults(step);
      return;
    }
    root.innerHTML = `
      <div class="orient">
        <div class="oprogress">Task ${step.index} of ${step.total} · <span class="mut">${esc(step.category)}</span></div>
        <div class="probhead"><span class="ptitle-h">${esc(step.title)}</span></div>
        <div class="prompt-text">${esc(step.prompt)}</div>
        <div class="term devterm">
          <div id="dLog"></div>
          <div class="dline"><span class="mut">❯</span> <input id="dInput" class="devinput" spellcheck="false"
            autocomplete="off" placeholder="type the command, then Enter…" /></div>
        </div>
        <div id="dOutcome" class="verdict"><div class="vmuted">This is a checker, not a real shell — nothing runs.</div></div>
      </div>`;
    renderLog();
    const input = root.querySelector('#dInput') as HTMLInputElement;
    input.focus();
    input.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' && input.value.trim()) void submit(input.value);
    });
  }

  function renderLog(): void {
    const log = root.querySelector('#dLog') as HTMLElement | null;
    if (log) log.innerHTML = transcript.join('');
  }

  async function submit(ans: string): Promise<void> {
    const out = await SubmitDevCommand(ans);
    if (disposed) return;
    transcript.push(`<div><span class="mut">❯</span> ${esc(ans)}</div>`);
    if (out.passed) {
      transcript = []; // next task starts with a fresh transcript
      renderStep(out.next);
      const o = root.querySelector('#dOutcome');
      if (o) o.innerHTML = `<div class="vhead vok">✓ ${esc(out.success)}</div>`;
    } else {
      transcript.push(`<div class="y">✖ not quite — hint: ${esc(out.hint)}</div>`);
      renderLog();
      const input = root.querySelector('#dInput') as HTMLInputElement;
      input.value = '';
      input.focus();
    }
  }

  function renderResults(step: guiapi.DevLitStep): void {
    root.innerHTML = `
      <div class="centerpane">
        <div class="card-lg">
          <h2 class="h2">Dev-Literacy — done</h2>
          <p class="bigscore">${step.passed} / ${step.total}</p>
          <p class="lead">You can read a directory, inspect files, and ask git what's going on — that's what the bench assumes. You can redo this track any time from the tab.</p>
          <div class="levelrow">
            <button class="btn" id="dRedo">↺ New set</button>
            <button class="btn primary" id="dBench">Go to the Bench →</button>
          </div>
        </div>
      </div>`;
    (root.querySelector('#dRedo') as HTMLButtonElement).addEventListener('click', () => void start());
    (root.querySelector('#dBench') as HTMLButtonElement).addEventListener('click', () => nav('bench'));
  }

  void start();
  return () => {
    disposed = true;
  };
}
