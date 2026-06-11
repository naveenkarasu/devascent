package guiapi

// ── Capability gating ─────────────────────────────────────────────────────────
// Two independent axes, never conflated (runtime-detection design spec): (A)
// runtime availability — is the toolchain installed — handled here via the
// shared Detector; (B) grading maturity — handled by GradedLanguages. Presence
// is the cheap sweep for picker marks; Capability is the authoritative
// compile+run canary (cached, so a session pays it once per language).

import (
	"context"
	"runtime"

	"devascent/internal/toolchain"
)

// LangStatus is one language's detected toolchain state.
type LangStatus struct {
	Lang     string `json:"lang"`
	Status   string `json:"status"`   // available | missing | broken | unknown
	Verified bool   `json:"verified"` // capability-verified (vs presence-only)
	Version  string `json:"version"`  // best-effort ("3.13.1"); may be empty
	Reason   string `json:"reason"`   // why missing/broken — shown on the install panel
}

func probeView(p toolchain.Probe) LangStatus {
	return LangStatus{
		Lang:     p.Lang,
		Status:   p.Status.String(),
		Verified: p.Depth == toolchain.DepthCapability,
		Version:  p.Version,
		Reason:   p.Reason,
	}
}

// LangPresence runs the fast presence sweep (LookPath + version) over the
// graded languages — cheap enough for the selector's marks at startup.
func (e *Engine) LangPresence() []LangStatus {
	out := make([]LangStatus, 0, len(GradedLanguages()))
	for _, lang := range GradedLanguages() {
		out = append(out, probeView(e.det.Presence(lang)))
	}
	return out
}

// CheckLang runs the authoritative capability canary (real compile+run; cached
// in the shared detector, so only the first check per session pays the cost).
func (e *Engine) CheckLang(lang string) LangStatus {
	return probeView(e.det.Capability(context.Background(), lang))
}

// RecheckLang drops the cached probe and re-runs the canary — the "I just
// installed it, check again" button.
func (e *Engine) RecheckLang(lang string) LangStatus {
	e.det.Invalidate(lang)
	return probeView(e.det.Capability(context.Background(), lang))
}

// InstallGuideView is one language's install guide resolved for this OS.
type InstallGuideView struct {
	Found  bool     `json:"found"`
	Lang   string   `json:"lang"`
	Label  string   `json:"label"`
	Notes  string   `json:"notes"`
	OS     string   `json:"os"` // windows | macos | linux
	Link   string   `json:"link"`
	Steps  []string `json:"steps"`
	Verify string   `json:"verify"`
}

func osKey() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "darwin":
		return "macos"
	default:
		return "linux"
	}
}

// InstallGuideFor returns lang's install guide for the running OS.
func (e *Engine) InstallGuideFor(lang string) InstallGuideView {
	g, ok := e.cat.InstallGuideForLang(lang)
	if !ok {
		return InstallGuideView{Found: false, Lang: lang, OS: osKey()}
	}
	v := InstallGuideView{Found: true, Lang: lang, Label: g.Label, Notes: g.Notes, OS: osKey()}
	if st, ok := g.OS[v.OS]; ok {
		v.Link = st.Link
		v.Steps = st.Steps
		v.Verify = st.Verify
	}
	return v
}
