package economy

// Package economy implements the Track-A hint economy: a free, slowly
// recharging pool of NUDGE charges plus earned MENTOR TOKENS that buy the two
// deeper hint tiers. Spending a paid tier discounts the solve's mastery
// weight, which the write-up credit ladder (A1) and the graduation gate (A3)
// consume. Pure logic — callers own persistence (save.State) and pass the
// clock in, so everything is deterministic under test.

import (
	"time"

	"devascent/internal/save"
)

// Tuning values — playtest starting points (decided 2026-06-12).
const (
	StartTokens     = 3 // granted once, when the wallet first initializes
	CleanSolveAward = 1 // banking a problem with no paid hints
	HardSolveAward  = 2 // same, when the problem is hard
	WriteupAward    = 1 // completing a write-up (A1)
	MilestoneAward  = 3 // completing a gate category minimum (A3)

	StrategyCost    = 1
	WalkthroughCost = 3

	NudgeMax      = 3
	NudgeRecharge = 20 * time.Minute

	PityMinFails   = 5                // failed grade attempts before pity unlocks…
	PityMinElapsed = 20 * time.Minute // …and minimum time on the problem
)

// Hint tiers. Nudges are free and never recorded against a solve; only the
// paid tiers discount mastery.
const (
	TierNone        = 0
	TierNudge       = 1
	TierStrategy    = 2
	TierWalkthrough = 3
)

// MasteryWeight is how much a solve counts toward mastery after hints.
func MasteryWeight(tier int) float64 {
	switch tier {
	case TierWalkthrough:
		return 0.4
	case TierStrategy:
		return 0.7
	default:
		return 1.0
	}
}

// HintCost returns the token price of a tier (nudges cost a charge, not tokens).
func HintCost(tier int) int {
	switch tier {
	case TierStrategy:
		return StrategyCost
	case TierWalkthrough:
		return WalkthroughCost
	default:
		return 0
	}
}

// SolveAward is the token payout for banking a problem with no paid hints.
func SolveAward(difficulty string) int {
	if difficulty == "hard" {
		return HardSolveAward
	}
	return CleanSolveAward
}

// Wallet is the in-memory view of the save's wallet fields.
type Wallet struct {
	Tokens       int
	Init         bool
	NudgeCharges int
	RechargeAt   time.Time // accrual anchor; zero until initialized
}

// Load builds a Wallet from a save state, granting the starting stash on
// first use and accruing any nudge recharges earned since the last play.
func Load(st *save.State, now time.Time) Wallet {
	w := Wallet{Tokens: st.Tokens, Init: st.WalletInit, NudgeCharges: st.NudgeCharges}
	if st.NudgeRechargeAt != "" {
		if t, err := time.Parse(time.RFC3339, st.NudgeRechargeAt); err == nil {
			w.RechargeAt = t
		}
	}
	if !w.Init {
		w.Init = true
		w.Tokens = StartTokens
		w.NudgeCharges = NudgeMax
		w.RechargeAt = now
	}
	w.Recharge(now)
	return w
}

// Store writes the wallet back onto a save state.
func (w Wallet) Store(st *save.State) {
	st.Tokens = w.Tokens
	st.WalletInit = w.Init
	st.NudgeCharges = w.NudgeCharges
	if !w.RechargeAt.IsZero() {
		st.NudgeRechargeAt = w.RechargeAt.UTC().Format(time.RFC3339)
	}
}

// Recharge accrues nudge charges: one per full NudgeRecharge period since the
// anchor, capped at NudgeMax. At cap the anchor tracks now so the next spend
// starts a fresh period instead of getting banked time for free.
func (w *Wallet) Recharge(now time.Time) {
	if w.RechargeAt.IsZero() || w.RechargeAt.After(now) {
		w.RechargeAt = now
	}
	for w.NudgeCharges < NudgeMax && now.Sub(w.RechargeAt) >= NudgeRecharge {
		w.NudgeCharges++
		w.RechargeAt = w.RechargeAt.Add(NudgeRecharge)
	}
	if w.NudgeCharges >= NudgeMax {
		w.NudgeCharges = NudgeMax
		w.RechargeAt = now
	}
}

// NextRecharge reports how long until the next nudge charge accrues
// (zero when already at cap).
func (w *Wallet) NextRecharge(now time.Time) time.Duration {
	w.Recharge(now)
	if w.NudgeCharges >= NudgeMax {
		return 0
	}
	d := NudgeRecharge - now.Sub(w.RechargeAt)
	if d < 0 {
		return 0
	}
	return d
}

// SpendNudge consumes one nudge charge; false when empty.
func (w *Wallet) SpendNudge(now time.Time) bool {
	w.Recharge(now)
	if w.NudgeCharges <= 0 {
		return false
	}
	if w.NudgeCharges == NudgeMax {
		w.RechargeAt = now // first spend below cap starts the refill clock
	}
	w.NudgeCharges--
	return true
}

// Spend debits tokens; false (and no change) when the balance is short.
func (w *Wallet) Spend(cost int) bool {
	if cost < 0 || w.Tokens < cost {
		return false
	}
	w.Tokens -= cost
	return true
}

// Refund returns tokens (mentor call failed after debit).
func (w *Wallet) Refund(cost int) {
	if cost > 0 {
		w.Tokens += cost
	}
}

// Award credits tokens.
func (w *Wallet) Award(n int) {
	if n > 0 {
		w.Tokens += n
	}
}

// PityEligible reports whether the one-time free strategy hint unlocks for a
// problem: enough failed attempts, enough time invested, not already used.
func PityEligible(rec save.SolveRecord, now time.Time) bool {
	if rec.PityUsed || rec.FailedRuns < PityMinFails || rec.FirstTryAt == "" {
		return false
	}
	first, err := time.Parse(time.RFC3339, rec.FirstTryAt)
	if err != nil {
		return false
	}
	return now.Sub(first) >= PityMinElapsed
}
