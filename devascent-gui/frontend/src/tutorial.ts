import { monaco, makeEditor } from './monaco-setup';
import {
  StartTutorial,
  GetLesson,
  GradeLessonStage,
  ResumeTutorial,
  AdvanceTutorial,
} from '../wailsjs/go/main/App';
import { guiapi } from '../wailsjs/go/models';
import { esc, renderVerdict } from './util';

function nav(view: string): void {
  window.dispatchEvent(new CustomEvent('devascent-nav', { detail: view }));
}

const KIND_LABEL: Record<string, string> = { i_do: 'Watch', we_do: 'Together', you_do: 'Your turn' };

// mountTutorial renders Tutorial Island as a stepper — one stage at a time,
// pass-gated forward movement, frontier persisted per language — and returns a
// disposer.
export function mountTutorial(root: HTMLElement, lang: string): () => void {
  let editor: monaco.editor.IStandaloneCodeEditor | null = null;
  let disposed = false;
  let total = 0;
  let lesson: guiapi.LessonView | null = null;
  let li = 0; // shown position
  let si = 0;
  let fLi = 0; // persisted frontier (furthest reached)
  let fSi = 0;
  const passed = new Set<string>(); // "li:si" stages passed this session

  const disposeEditor = () => {
    if (editor) {
      editor.dispose();
      editor = null;
    }
  };
  const cmp = (al: number, as_: number, bl: number, bs: number) => al - bl || as_ - bs;

  async function start(): Promise<void> {
    total = await StartTutorial(lang);
    if (disposed) return;
    if (total === 0) {
      root.innerHTML = `<div class="centerpane"><div class="card-lg"><h2 class="h2">No lessons for ${esc(lang)} yet</h2></div></div>`;
      return;
    }
    const pos = await ResumeTutorial();
    if (disposed) return;
    fLi = pos.lesson;
    fSi = pos.stage;
    if (pos.done) {
      renderComplete();
      return;
    }
    void show(fLi, fSi);
  }

  // satisfied: the shown stage permits moving past it — no task, passed this
  // session, or it sits behind the frontier (completed in an earlier session).
  function satisfied(st: guiapi.LessonStageView): boolean {
    return !st.hasTask || passed.has(`${li}:${si}`) || cmp(li, si, fLi, fSi) < 0;
  }

  function nextPos(): [number, number] | null {
    if (!lesson) return null;
    if (si + 1 < lesson.stages.length) return [li, si + 1];
    return [li + 1, 0]; // li+1 === total → tutorial complete
  }

  function updateNext(): void {
    const btn = root.querySelector('#tNext') as HTMLButtonElement | null;
    if (!btn || !lesson) return;
    const np = nextPos();
    if (!np) return;
    const unlocked = cmp(np[0], np[1], fLi, fSi) <= 0 || satisfied(lesson.stages[si]);
    btn.disabled = !unlocked;
    btn.title = unlocked ? '' : 'Pass this stage to continue';
  }

  async function onNext(): Promise<void> {
    const np = nextPos();
    if (!np) return;
    if (cmp(np[0], np[1], fLi, fSi) > 0) {
      const p = await AdvanceTutorial(np[0], np[1]);
      if (disposed) return;
      fLi = p.lesson;
      fSi = p.stage;
      if (p.done) {
        renderComplete();
        return;
      }
    }
    if (np[0] >= total) {
      renderComplete();
      return;
    }
    void show(np[0], np[1]);
  }

  function onPrev(): void {
    if (si > 0) void show(li, si - 1);
    else if (li > 0) void show(li - 1, 0);
  }

  async function show(l: number, s: number): Promise<void> {
    disposeEditor();
    li = l;
    lesson = await GetLesson(l);
    if (disposed || !lesson.found) return;
    si = Math.min(s, lesson.stages.length - 1);
    const st = lesson.stages[si];

    const chips = lesson.stages
      .map((x, i) => {
        const done = cmp(li, i, fLi, fSi) < 0 || passed.has(`${li}:${i}`);
        const open = cmp(li, i, fLi, fSi) <= 0;
        const cls = i === si ? 'cur' : open ? 'open' : 'locked';
        return `<button class="stepchip ${cls}" data-si="${i}" ${open ? '' : 'disabled'}>
          ${done ? '✓ ' : ''}${esc(KIND_LABEL[x.kind] ?? x.kind)}</button>`;
      })
      .join('<span class="steparrow">→</span>');

    const task = st.hasTask
      ? `${st.prompt ? `<div class="prompt-text">${esc(st.prompt)}</div>` : ''}
         <div id="tEd" class="editor tall"></div>
         <div class="workbar"><button id="tRun" class="btn primary">&#9654; Run</button></div>
         <div id="tV" class="verdict"><div class="vmuted">Write the code, then Run. A pass unlocks the next stage.</div></div>`
      : '';

    root.innerHTML = `
      <div class="tutorial">
        <div class="tprogress">Lesson ${lesson.index} of ${lesson.total} · stage ${si + 1} of ${lesson.stages.length}</div>
        <h2 class="h2">${esc(lesson.title)}</h2>
        <div class="stepper">${chips}</div>
        <div class="stage">
          <div class="stagehead s-${esc(st.kind)}">${esc(KIND_LABEL[st.kind] ?? st.kind)} · ${esc(st.title)}</div>
          ${st.body ? `<div class="stagebody">${esc(st.body)}</div>` : ''}
          ${task}
        </div>
        <div class="workbar tnav">
          <button id="tPrev" class="btn" ${li === 0 && si === 0 ? 'disabled' : ''}>← Back</button>
          <button id="tNext" class="btn primary">${si + 1 < lesson.stages.length ? 'Next stage →' : li + 1 < total ? 'Next lesson →' : 'Finish tutorial ✓'}</button>
        </div>
      </div>`;

    if (st.hasTask) {
      editor = makeEditor(root.querySelector('#tEd') as HTMLElement, lang, st.starter);
      (root.querySelector('#tRun') as HTMLButtonElement).addEventListener('click', () => void run());
    }
    root.querySelectorAll<HTMLButtonElement>('.stepchip').forEach((c) =>
      c.addEventListener('click', () => void show(li, parseInt(c.dataset.si as string, 10))),
    );
    (root.querySelector('#tPrev') as HTMLButtonElement).addEventListener('click', onPrev);
    (root.querySelector('#tNext') as HTMLButtonElement).addEventListener('click', () => void onNext());
    updateNext();
  }

  async function run(): Promise<void> {
    if (!editor) return;
    const btn = root.querySelector('#tRun') as HTMLButtonElement;
    const v = root.querySelector('#tV') as HTMLElement;
    btn.disabled = true;
    btn.textContent = 'Running…';
    v.innerHTML = '<div class="vmuted">Running…</div>';
    try {
      const res = await GradeLessonStage(li, si, editor.getValue());
      if (disposed) return;
      v.innerHTML = renderVerdict(res);
      if (res.passed) {
        passed.add(`${li}:${si}`);
        void show(li, si); // re-render: chip turns ✓, Next unlocks
        return;
      }
    } finally {
      if (!disposed) {
        const b = root.querySelector('#tRun') as HTMLButtonElement | null;
        if (b) {
          b.disabled = false;
          b.textContent = '▶ Run';
        }
      }
    }
    updateNext();
  }

  function renderComplete(): void {
    disposeEditor();
    root.innerHTML = `
      <div class="centerpane">
        <div class="card-lg">
          <h2 class="h2">Tutorial Island — complete</h2>
          <p class="lead">You've shipped real ${esc(lang)} code through every lesson. The bench is where it counts now — real problems, hidden tests, your toolchain.</p>
          <div class="levelrow">
            <button class="btn" id="tReview">↺ Review lessons</button>
            <button class="btn primary" id="tToBench">Go to the Bench →</button>
          </div>
        </div>
      </div>`;
    (root.querySelector('#tReview') as HTMLButtonElement).addEventListener('click', () => void show(0, 0));
    (root.querySelector('#tToBench') as HTMLButtonElement).addEventListener('click', () => nav('bench'));
  }

  void start();
  return () => {
    disposed = true;
    disposeEditor();
  };
}
