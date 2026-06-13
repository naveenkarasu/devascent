package guiapi

import (
	"strings"
	"testing"

	"devascent/internal/economy"
	"devascent/internal/engine"
)

// trackAEngine builds an engine on a temp save dir with one solved problem.
func trackAEngine(t *testing.T) (*Engine, string) {
	t.Helper()
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	e, err := New()
	if err != nil {
		t.Fatalf("engine: %v", err)
	}
	return e, "python"
}

func TestWallet_StartingGrantOnceAndPersisted(t *testing.T) {
	e, lang := trackAEngine(t)
	w := e.Wallet(lang)
	if w.Tokens != economy.StartTokens || w.NudgeCharges != economy.NudgeMax {
		t.Fatalf("first wallet: %+v", w)
	}
	// Fresh engine over the same save dir: no re-grant.
	e2, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if w2 := e2.Wallet(lang); w2.Tokens != economy.StartTokens {
		t.Fatalf("re-grant detected: %+v", w2)
	}
}

func TestRequestHint_NudgeIsFreeOfTokensAndEscalates(t *testing.T) {
	e, lang := trackAEngine(t)
	id := e.cat.Problems[0].ID
	h1 := e.RequestHint(lang, id, economy.TierNudge, "")
	h2 := e.RequestHint(lang, id, economy.TierNudge, "")
	if h1.Err != "" || h2.Err != "" {
		t.Fatalf("nudges errored: %q %q", h1.Err, h2.Err)
	}
	if h1.Text == h2.Text {
		t.Fatal("nudges did not escalate")
	}
	if h2.Wallet.Tokens != economy.StartTokens {
		t.Fatalf("nudge cost tokens: %+v", h2.Wallet)
	}
	if h2.Wallet.NudgeCharges != economy.NudgeMax-2 {
		t.Fatalf("nudge charges: %+v", h2.Wallet)
	}
	// Third spends the last charge; fourth must refuse with a countdown.
	e.RequestHint(lang, id, economy.TierNudge, "")
	h4 := e.RequestHint(lang, id, economy.TierNudge, "")
	if h4.Err == "" || !strings.Contains(h4.Err, "recharges") {
		t.Fatalf("empty pool not refused: %+v", h4)
	}
}

func TestRequestHint_PaidTiersDebitRecordAndRepeatFree(t *testing.T) {
	e, lang := trackAEngine(t)
	p := e.cat.Problems[0]
	e.Wallet(lang) // grant

	h := e.RequestHint(lang, p.ID, economy.TierStrategy, "code")
	if h.Err != "" || h.Source != "template" {
		t.Fatalf("strategy hint: %+v", h)
	}
	if h.Wallet.Tokens != economy.StartTokens-economy.StrategyCost {
		t.Fatalf("tokens after strategy: %+v", h.Wallet)
	}
	// Same tier again: free (already paid).
	h2 := e.RequestHint(lang, p.ID, economy.TierStrategy, "code")
	if h2.Wallet.Tokens != h.Wallet.Tokens {
		t.Fatalf("repeat hint charged again: %+v", h2.Wallet)
	}
	// The tier is recorded for the mastery discount.
	e.mu.Lock()
	rec := e.getSlot(lang).record(p.ID)
	e.mu.Unlock()
	if rec.HintTier != economy.TierStrategy {
		t.Fatalf("hint tier not recorded: %+v", rec)
	}
	// Walkthrough costs 3: with 2 left this must refuse.
	h3 := e.RequestHint(lang, p.ID, economy.TierWalkthrough, "code")
	if h3.Err == "" {
		t.Fatalf("overdraft allowed: %+v", h3)
	}
}

func TestGradePath_CleanBankAwardsAndOpensWriteup(t *testing.T) {
	e, lang := trackAEngine(t)
	p := e.cat.Problems[0]
	e.Wallet(lang)

	// Simulate the bench Grade success path without a real toolchain.
	newly, saveErr := e.bank(lang, p.ID)
	if !newly || saveErr != "" {
		t.Fatalf("bank: %v %q", newly, saveErr)
	}
	awarded, pending := e.trackABank(lang, p)
	if awarded != economy.SolveAward(p.Difficulty) || !pending {
		t.Fatalf("clean bank: awarded=%d pending=%v", awarded, pending)
	}

	// A hinted problem banks with no award.
	p2 := e.cat.Problems[1]
	e.RequestHint(lang, p2.ID, economy.TierStrategy, "")
	e.bank(lang, p2.ID)
	if a2, _ := e.trackABank(lang, p2); a2 != 0 {
		t.Fatalf("hinted bank still awarded %d", a2)
	}
}

func TestWriteupFlow_MCQGateAndAward(t *testing.T) {
	e, lang := trackAEngine(t)
	p := e.cat.Problems[0]
	e.Wallet(lang)
	e.bank(lang, p.ID)
	e.trackABank(lang, p)

	v := e.Writeup(lang, p.ID)
	if !v.Solved || v.Done || !v.HasMCQ || len(v.Options) != 4 {
		t.Fatalf("writeup view: %+v", v)
	}

	q, _ := engine.ComplexityMCQ(p)
	wrong := (q.Correct + 1) % 4
	text := "I used the standard approach: scan once and keep what matters."

	if r := e.SubmitWriteup(lang, p.ID, wrong, text); r.Accepted || r.MCQCorrect {
		t.Fatalf("wrong MCQ accepted: %+v", r)
	}
	if r := e.SubmitWriteup(lang, p.ID, q.Correct, "short"); r.Accepted || r.Err == "" {
		t.Fatalf("short text accepted: %+v", r)
	}

	before := e.Wallet(lang).Tokens
	r := e.SubmitWriteup(lang, p.ID, q.Correct, text)
	if !r.Accepted || r.TokensAwarded < economy.WriteupAward {
		t.Fatalf("good writeup rejected: %+v", r)
	}
	if r.Wallet.Tokens != before+r.TokensAwarded {
		t.Fatalf("award not applied: before=%d %+v", before, r)
	}
	if rr := e.SubmitWriteup(lang, p.ID, q.Correct, text); rr.Err == "" {
		t.Fatal("double submission allowed")
	}
	if !e.Writeup(lang, p.ID).Done {
		t.Fatal("writeup not marked done")
	}
}

func TestGate_ProvisionalVsFull(t *testing.T) {
	e, lang := trackAEngine(t)
	e.Wallet(lang)

	// Find a blind75 problem and bank it without a write-up.
	var b75 string
	for _, p := range e.cat.Problems {
		for _, l := range p.Lists {
			if l == "blind75" {
				b75 = p.ID
				break
			}
		}
		if b75 != "" {
			break
		}
	}
	if b75 == "" {
		t.Fatal("no blind75 problem in catalog")
	}
	e.bank(lang, b75)

	g := e.Gate(lang)
	if g.Full != 0 || g.Provisional != 1 {
		t.Fatalf("gate before writeup: full=%d prov=%d", g.Full, g.Provisional)
	}
	q, _ := engine.ComplexityMCQ(e.probByID[b75])
	r := e.SubmitWriteup(lang, b75, q.Correct, "Standard technique for this category, explained briefly here.")
	if !r.Accepted {
		t.Fatalf("writeup: %+v", r)
	}
	g = e.Gate(lang)
	if g.Full != 1 || g.Provisional != 0 {
		t.Fatalf("gate after writeup: full=%d prov=%d", g.Full, g.Provisional)
	}
	if g.Target != engine.Blind75Target || len(g.Categories) == 0 || len(g.Mandatory) != 5 {
		t.Fatalf("gate shape: %+v", g)
	}
}

func TestMentorFacade_TemplatesAndPreview(t *testing.T) {
	e, lang := trackAEngine(t)
	p := e.cat.Problems[0]

	rows := e.MentorBackends()
	if len(rows) < 6 || rows[0].ID != "template" || !rows[0].Present {
		t.Fatalf("backends: %+v", rows)
	}
	if msg := e.SelectMentor("template"); msg != "" {
		t.Fatalf("template select: %q", msg)
	}
	prev := e.MentorPreview(lang, p.ID, economy.TierStrategy, "def f(): pass")
	if !strings.Contains(prev, p.Title) || !strings.Contains(prev, "def f()") {
		t.Fatalf("preview missing context: %q", prev)
	}
	if strings.Contains(prev, p.Solution) && p.Solution != "" {
		t.Fatal("preview leaked the canonical solution")
	}
}
