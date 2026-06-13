import { guiapi } from '../wailsjs/go/models';

// esc HTML-escapes a string for safe innerHTML interpolation.
export function esc(s: string): string {
  return s.replace(
    /[&<>"]/g,
    (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' })[c] as string,
  );
}

// renderInstallGuide turns an install guide into HTML (the caller appends its
// own Re-check button and wires it).
export function renderInstallGuide(g: guiapi.InstallGuideView, reason: string): string {
  const head = `<div class="psectitle">Install ${esc(g.label || g.lang)}</div>`;
  const why = reason ? `<div class="poplabel" style="margin-bottom:8px">${esc(reason)}</div>` : '';
  if (!g.found) {
    return head + why + `<div class="vmuted">No install guide for ${esc(g.lang)} yet — see INSTALL.md in the repo.</div>`;
  }
  const notes = g.notes ? `<div class="psummary">${esc(g.notes)}</div>` : '';
  const link = g.link
    ? `<div class="pop"><div class="poplabel">Download</div><pre class="popcode">${esc(g.link)}</pre></div>`
    : '';
  const steps = (g.steps || [])
    .map((s, i) => `<div class="pop"><div class="poplabel">Step ${i + 1}</div><pre class="popcode">${esc(s)}</pre></div>`)
    .join('');
  const verify = g.verify
    ? `<div class="pop"><div class="poplabel">Verify it worked</div><pre class="popcode">${esc(g.verify)}</pre></div>`
    : '';
  return head + why + notes + link + steps + verify;
}

// renderOutput shows the player's own print/log output (for debugging),
// separate from the pass/fail verdict. Empty when the program printed nothing.
function renderOutput(stdout: string): string {
  if (!stdout) return '';
  return `<div class="voutlabel">Your output (print / log)</div><pre class="vout">${esc(stdout)}</pre>`;
}

// renderVerdict turns a GradeResult into HTML (shared by every graded view).
export function renderVerdict(res: guiapi.GradeResult): string {
  if (res.err) {
    return `<div class="vhead vfail">✖ Did not run</div><pre class="vpre">${esc(res.err)}</pre>${renderOutput(res.stdout)}`;
  }
  const passedCount = res.casesTotal - res.casesFailed;
  const head = res.passed
    ? `<div class="vhead vok">✓ Passed — ${res.casesTotal}/${res.casesTotal} hidden tests</div>`
    : `<div class="vhead vfail">✖ Failed — ${passedCount}/${res.casesTotal} hidden tests</div>`;
  const rows = (res.results || [])
    .map(
      (r) => `<div class="crow ${r.passed ? 'cok' : 'cfail'}">
        <span class="cmark">${r.passed ? '✓' : '✖'}</span>
        <span class="cname">${esc(r.name)}</span>
        ${
          r.passed
            ? ''
            : `<span class="cdetail">got <code>${esc(r.got)}</code> · want <code>${esc(r.expected)}</code>${r.err ? ' · ' + esc(r.err) : ''}</span>`
        }</div>`,
    )
    .join('');
  return head + `<div class="clist">${rows}</div>` + renderOutput(res.stdout);
}
