import { monaco, makeEditor } from './monaco-setup';
import { GetAdvancedTopics, GetAdvancedTopic, GradeAdvanced } from '../wailsjs/go/main/App';
import { guiapi } from '../wailsjs/go/models';
import { esc, renderVerdict } from './util';

const KIND_LABEL: Record<string, string> = {
  'fix-it': 'Fix it',
  'spot-the-bug': 'Spot the bug',
  implement: 'Implement',
  'predict-output': 'Predict the output',
  'explain-and-implement': 'Explain & implement',
};

// renderAdvVerdict adapts check-style grading (compiles / stdout) where there
// are no per-case rows — renderVerdict's "n/n hidden tests" would read oddly.
function renderAdvVerdict(res: guiapi.GradeResult, check: string): string {
  if (res.err) {
    return `<div class="vhead vfail">✖ Not fixed yet</div><pre class="vpre">${esc(res.err)}</pre>`;
  }
  if (res.casesTotal > 0) return renderVerdict(res);
  if (res.passed) {
    const what = check === 'stdout' ? 'output matches' : 'compiles cleanly';
    return `<div class="vhead vok">✓ Fixed — ${what}</div>`;
  }
  return `<div class="vhead vfail">✖ Not fixed yet${check === 'stdout' ? ' — output differs' : ''}</div>`;
}

// mountAdvanced renders the Advanced Topics view (deep dives + fix-it
// exercises over broken code, graded by the real toolchain) and returns a
// disposer.
export function mountAdvanced(root: HTMLElement, lang: string): () => void {
  const editors: monaco.editor.IStandaloneCodeEditor[] = [];
  let disposed = false;
  let topics: guiapi.AdvTopicSummary[] = [];
  let currentIdx = -1;
  const disposeEditors = () => {
    editors.forEach((e) => e.dispose());
    editors.length = 0;
  };

  root.innerHTML = `
    <div class="bench">
      <aside class="sidebar">
        <div class="pcount" style="padding-top:12px">Advanced topics · ${esc(lang)}</div>
        <div id="aList" class="plist"></div>
      </aside>
      <main class="work">
        <div id="aBody" class="advbody"><div class="vmuted" style="padding:16px">Pick a topic.</div></div>
      </main>
    </div>`;
  const el = (id: string) => root.querySelector('#' + id) as HTMLElement;

  function renderList(): void {
    let lastGroup = '';
    let html = '';
    for (const t of topics) {
      if (t.group !== lastGroup) {
        lastGroup = t.group;
        html += `<div class="pcat">${esc(t.group)}</div>`;
      }
      html += `<div class="prow ${t.index === currentIdx ? 'active' : ''}" data-idx="${t.index}">
        <span class="ptitle">${esc(t.title.split(' — ')[0])}</span>
        <span class="pdiff">${t.gradeable}/${t.exercises} ✎</span></div>`;
    }
    el('aList').innerHTML = html;
    el('aList')
      .querySelectorAll<HTMLElement>('.prow')
      .forEach((r) => r.addEventListener('click', () => void open(parseInt(r.dataset.idx as string, 10))));
  }

  async function open(idx: number): Promise<void> {
    currentIdx = idx;
    disposeEditors();
    const d = await GetAdvancedTopic(lang, idx);
    if (disposed || !d.found) return;
    const sections = (d.sections || [])
      .map(
        (s) => `<div class="psec"><div class="psectitle">${esc(s.title)}</div>
        ${(s.ops || [])
          .map((op) => `<div class="pop"><div class="poplabel">${esc(op.label)}</div><pre class="popcode">${esc(op.code.replace(/\s+$/, ''))}</pre></div>`)
          .join('')}</div>`,
      )
      .join('');
    const exercises = (d.exercises || [])
      .map((ex) => {
        const head = `<div class="stagehead s-you_do">${esc(KIND_LABEL[ex.kind] ?? ex.kind)}${ex.gradeable ? '' : ' · reading'}</div>`;
        const prompt = ex.prompt ? `<div class="stagebody">${esc(ex.prompt)}</div>` : '';
        if (!ex.gradeable) {
          // reveal-only: show the broken code + explanation inline
          return `<div class="stage">${head}${prompt}
            ${ex.brokenCode ? `<pre class="popcode">${esc(ex.brokenCode.replace(/\s+$/, ''))}</pre>` : ''}
            ${ex.bug ? `<div class="stagebody" style="margin-top:8px">${esc(ex.bug)}</div>` : ''}
          </div>`;
        }
        return `<div class="stage">${head}${prompt}
          <div id="aEd-${ex.index}" class="editor tall"></div>
          <div class="workbar">
            <button id="aRun-${ex.index}" class="btn primary">&#9654; Check fix</button>
            <button id="aRev-${ex.index}" class="btn">&#128065; Reveal answer</button>
          </div>
          <div id="aV-${ex.index}" class="verdict"><div class="vmuted">Fix the broken code, then check.</div></div>
          <div id="aSol-${ex.index}" class="hidden">
            ${ex.bug ? `<div class="stagebody" style="margin:8px 0">${esc(ex.bug)}</div>` : ''}
            <pre class="popcode">${esc(ex.fixedCode.replace(/\s+$/, ''))}</pre>
          </div>
        </div>`;
      })
      .join('');
    el('aBody').innerHTML = `
      <div class="advhead">
        <span class="ptitle-h">${esc(d.title)}</span>
        ${d.tag ? `<span class="tag cat">${esc(d.tag)}</span>` : ''}
      </div>
      ${d.summary ? `<div class="psummary" style="padding:0 16px">${esc(d.summary)}</div>` : ''}
      <div style="padding:0 16px">${sections}</div>
      <div class="stages" style="padding:0 16px 20px">${exercises}</div>`;
    (d.exercises || []).forEach((ex) => {
      if (!ex.gradeable) return;
      const ed = makeEditor(root.querySelector(`#aEd-${ex.index}`) as HTMLElement, lang, ex.brokenCode);
      editors.push(ed);
      (root.querySelector(`#aRun-${ex.index}`) as HTMLButtonElement).addEventListener('click', () =>
        void check(ex, ed),
      );
      (root.querySelector(`#aRev-${ex.index}`) as HTMLButtonElement).addEventListener('click', () => {
        (root.querySelector(`#aSol-${ex.index}`) as HTMLElement).classList.toggle('hidden');
      });
    });
    renderList();
  }

  async function check(ex: guiapi.AdvExerciseView, ed: monaco.editor.IStandaloneCodeEditor): Promise<void> {
    const v = root.querySelector(`#aV-${ex.index}`) as HTMLElement;
    const btn = root.querySelector(`#aRun-${ex.index}`) as HTMLButtonElement;
    btn.disabled = true;
    btn.textContent = 'Checking…';
    v.innerHTML = '<div class="vmuted">Compiling + checking in your toolchain…</div>';
    try {
      const res = await GradeAdvanced(lang, currentIdx, ex.index, ed.getValue());
      if (disposed) return;
      v.innerHTML = renderAdvVerdict(res, ex.check);
    } finally {
      if (!disposed) {
        btn.disabled = false;
        btn.textContent = '▶ Check fix';
      }
    }
  }

  void (async () => {
    topics = await GetAdvancedTopics(lang);
    if (disposed) return;
    if (!topics || !topics.length) {
      el('aBody').innerHTML = `<div class="vmuted" style="padding:16px">No advanced topics for ${esc(lang)} yet.</div>`;
      topics = [];
      return;
    }
    renderList();
    void open(0);
  })();

  return () => {
    disposed = true;
    disposeEditors();
  };
}
