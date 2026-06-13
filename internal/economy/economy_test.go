package economy

import (
	"testing"
	"time"

	"devascent/internal/save"
)

var t0 = time.Date(2026, 6, 12, 12, 0, 0, 0, time.UTC)

func TestLoadGrantsStartingStashOnce(t *testing.T) {
	st := &save.State{}
	w := Load(st, t0)
	if !w.Init || w.Tokens != StartTokens || w.NudgeCharges != NudgeMax {
		t.Fatalf("first load: got %+v", w)
	}
	w.Spend(2)
	w.SpendNudge(t0)
	w.Store(st)

	w2 := Load(st, t0.Add(time.Minute))
	if w2.Tokens != StartTokens-2 || w2.NudgeCharges != NudgeMax-1 {
		t.Fatalf("second load re-granted: got %+v", w2)
	}
}

func TestNudgeRechargeAccrual(t *testing.T) {
	st := &save.State{}
	w := Load(st, t0)
	for i := 0; i < NudgeMax; i++ {
		if !w.SpendNudge(t0) {
			t.Fatalf("spend %d failed", i)
		}
	}
	if w.SpendNudge(t0) {
		t.Fatal("spend on empty pool succeeded")
	}
	// 2 full periods + change → exactly 2 charges back.
	w.Recharge(t0.Add(2*NudgeRecharge + time.Minute))
	if w.NudgeCharges != 2 {
		t.Fatalf("charges after 2 periods = %d, want 2", w.NudgeCharges)
	}
	// A very long absence caps at NudgeMax.
	w.Recharge(t0.Add(48 * time.Hour))
	if w.NudgeCharges != NudgeMax {
		t.Fatalf("charges after 48h = %d, want %d", w.NudgeCharges, NudgeMax)
	}
}

func TestNudgeClockStartsOnFirstSpendBelowCap(t *testing.T) {
	st := &save.State{}
	w := Load(st, t0)
	// At cap for an hour: no banked time — spending then waiting one period
	// yields exactly one charge back.
	later := t0.Add(time.Hour)
	w.Recharge(later)
	if !w.SpendNudge(later) {
		t.Fatal("spend at cap failed")
	}
	if got := w.NextRecharge(later); got != NudgeRecharge {
		t.Fatalf("next recharge = %v, want full period", got)
	}
	w.Recharge(later.Add(NudgeRecharge))
	if w.NudgeCharges != NudgeMax {
		t.Fatalf("charges = %d, want back at cap", w.NudgeCharges)
	}
}

func TestSpendRefundAward(t *testing.T) {
	w := Wallet{Tokens: 1, Init: true}
	if w.Spend(WalkthroughCost) {
		t.Fatal("overdraft allowed")
	}
	if w.Tokens != 1 {
		t.Fatal("failed spend mutated balance")
	}
	if !w.Spend(StrategyCost) || w.Tokens != 0 {
		t.Fatalf("spend: %+v", w)
	}
	w.Refund(StrategyCost)
	w.Award(SolveAward("hard"))
	if w.Tokens != 1+HardSolveAward {
		t.Fatalf("tokens = %d", w.Tokens)
	}
}

func TestPityEligible(t *testing.T) {
	first := t0.Format(time.RFC3339)
	now := t0.Add(PityMinElapsed)
	cases := []struct {
		name string
		rec  save.SolveRecord
		now  time.Time
		want bool
	}{
		{"fails + time", save.SolveRecord{FailedRuns: PityMinFails, FirstTryAt: first}, now, true},
		{"too few fails (under solo time)", save.SolveRecord{FailedRuns: PityMinFails - 1, FirstTryAt: first}, now, false},
		{"too soon", save.SolveRecord{FailedRuns: PityMinFails, FirstTryAt: first}, t0.Add(time.Minute), false},
		{"already used", save.SolveRecord{FailedRuns: PityMinFails, FirstTryAt: first, PityUsed: true}, now, false},
		{"no first try", save.SolveRecord{FailedRuns: PityMinFails}, now, false},
		// solo-time path: enough time alone unlocks it even with 0 distinct fails…
		{"solo time, no fails", save.SolveRecord{FailedRuns: 0, FirstTryAt: first}, t0.Add(PitySoloElapsed), true},
		// …but below the solo threshold, 0 fails stays locked.
		{"below solo, no fails", save.SolveRecord{FailedRuns: 0, FirstTryAt: first}, t0.Add(PitySoloElapsed - time.Minute), false},
	}
	for _, c := range cases {
		if got := PityEligible(c.rec, c.now); got != c.want {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}

func TestFailHash_DistinctOnly(t *testing.T) {
	a := FailHash("def f():\n    return 1")
	aSpaced := FailHash("  def f():\n    return 1  ") // edge whitespace ignored
	b := FailHash("def f():\n    return 2")           // a real edit
	if a != aSpaced {
		t.Fatal("trailing/leading whitespace should not change the fingerprint")
	}
	if a == b {
		t.Fatal("a genuine code change must change the fingerprint")
	}
}

func TestMasteryWeights(t *testing.T) {
	if MasteryWeight(TierNone) != 1.0 || MasteryWeight(TierNudge) != 1.0 {
		t.Fatal("free tiers must not discount")
	}
	if MasteryWeight(TierStrategy) != 0.7 || MasteryWeight(TierWalkthrough) != 0.4 {
		t.Fatal("paid tier weights wrong")
	}
}
