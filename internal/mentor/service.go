package mentor

// Service orchestrates the seam: backend registry, two-phase detection
// (Present → canary Probe, cached), and the Hint call that enforces the
// guardrails. The game never blocks on AI: any failure quietly answers from
// templates and reports FellBack so the caller can refund the token.

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// Status is one row of the mentor picker.
type Status struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Present  bool   `json:"present"`
	Info     string `json:"info"`
	Selected bool   `json:"selected"`
	Probed   bool   `json:"probed"`   // a probe has run this session
	ProbeOK  bool   `json:"probeOk"`  // …and succeeded
	ProbeErr string `json:"probeErr"` // …or failed with this
}

type Service struct {
	mu       sync.Mutex
	cfg      Config
	backends map[string]Backend
	order    []string
	probeOK  map[string]bool
	probeErr map[string]string
}

// NewService builds the registry from the stored config.
func NewService(cfg Config) *Service {
	s := &Service{
		cfg:      cfg,
		backends: map[string]Backend{},
		probeOK:  map[string]bool{},
		probeErr: map[string]string{},
	}
	for _, b := range []Backend{
		newOllama(cfg), newOpenAI(cfg), newClaude(cfg), newCodex(cfg), newCopilot(cfg),
	} {
		s.backends[b.ID()] = b
		s.order = append(s.order, b.ID())
	}
	return s
}

// Config returns the current mentor configuration.
func (s *Service) Config() Config {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cfg
}

// Statuses lists templates + every backend for the picker (presence only —
// probing is explicit because it costs a real AI call).
func (s *Service) Statuses() []Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Status{{
		ID: "template", Name: "Built-in templates (offline)", Present: true,
		Info: "always available", Selected: s.cfg.Backend == "", Probed: true, ProbeOK: true,
	}}
	for _, id := range s.order {
		b := s.backends[id]
		present, info := b.Present()
		st := Status{ID: id, Name: b.Name(), Present: present, Info: info, Selected: s.cfg.Backend == id}
		if ok, probed := s.probeOK[id], s.probeErr[id]; probed != "" || ok {
			st.Probed = true
			st.ProbeOK = ok
			st.ProbeErr = probed
		}
		out = append(out, st)
	}
	return out
}

// Probe runs the canary round-trip against one backend and caches the result.
func (s *Service) Probe(ctx context.Context, id string) error {
	if id == "" || id == "template" {
		return nil
	}
	s.mu.Lock()
	b := s.backends[id]
	s.mu.Unlock()
	if b == nil {
		return fmt.Errorf("unknown mentor backend %q", id)
	}
	if present, _ := b.Present(); !present {
		err := fmt.Errorf("%s is not installed (or not reachable)", b.Name())
		s.recordProbe(id, err)
		return err
	}
	pctx, cancel := context.WithTimeout(ctx, ProbeTimeout)
	defer cancel()
	out, err := b.Ask(pctx, "Reply with exactly "+Canary+" and nothing else.")
	if err == nil && !strings.Contains(out, Canary) {
		err = fmt.Errorf("%s answered but not with the expected sentinel", b.Name())
	}
	s.recordProbe(id, err)
	return err
}

func (s *Service) recordProbe(id string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err == nil {
		s.probeOK[id] = true
		s.probeErr[id] = ""
	} else {
		s.probeOK[id] = false
		s.probeErr[id] = err.Error()
	}
}

// Select probes a backend and, on success, persists it as the mentor.
// Selecting "" or "template" always succeeds (offline mode).
func (s *Service) Select(ctx context.Context, id string) error {
	if id == "template" {
		id = ""
	}
	if id != "" {
		if err := s.Probe(ctx, id); err != nil {
			return err
		}
	}
	s.mu.Lock()
	s.cfg.Backend = id
	cfg := s.cfg
	s.mu.Unlock()
	return SaveConfig(cfg)
}

// SetEndpoint updates the openai-compat connection details (re-probed on select).
func (s *Service) SetEndpoint(endpoint, model, apiKey string) error {
	s.mu.Lock()
	s.cfg.Endpoint = endpoint
	s.cfg.Model = model
	s.cfg.APIKey = apiKey
	s.backends["openai-compat"] = newOpenAI(s.cfg)
	s.backends["ollama"] = newOllama(s.cfg)
	delete(s.probeOK, "openai-compat")
	delete(s.probeOK, "ollama")
	cfg := s.cfg
	s.mu.Unlock()
	return SaveConfig(cfg)
}

// AIEnabled reports whether an AI backend (not templates) is selected — the
// "AI if connected, else skipped" guard for follow-ups and review notes.
func (s *Service) AIEnabled() bool {
	return s.selected() != nil
}

// selected returns the active AI backend, nil for templates-only.
func (s *Service) selected() Backend {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cfg.Backend == "" {
		return nil
	}
	return s.backends[s.cfg.Backend]
}

// Preview is the transparency view: exactly what Hint would send.
func (s *Service) Preview(req Request) string {
	return BuildPrompt(req)
}

// Hint answers a tier-2/3/followup/review request. Tier-1 nudges never come
// here — call Nudge directly. Never returns an error: AI failure of any kind
// (timeout, garbage, guardrail violation) falls back to templates with
// FellBack=true so the caller refunds the token.
func (s *Service) Hint(ctx context.Context, req Request) Response {
	b := s.selected()
	if b == nil {
		return Response{Text: templateAnswer(req), Source: "template"}
	}
	actx, cancel := context.WithTimeout(ctx, AskTimeout)
	defer cancel()
	out, err := b.Ask(actx, BuildPrompt(req))
	if err == nil {
		err = Validate(req.Kind, req.Lang, out)
	}
	if err != nil {
		return Response{Text: templateAnswer(req), Source: "template", FellBack: true}
	}
	return Response{Text: strings.TrimSpace(scrub(out)), Source: b.ID()}
}
