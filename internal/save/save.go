package save

// Package save persists run progress as JSON (atomic writes) under the per-OS
// user config dir. v3: ONE SLOT PER LANGUAGE (save-<lang>.json) so a player
// trains multiple languages in parallel; the legacy single save.json is
// migrated into its language's slot on first access. Shared by the TUI and
// the GUI — neither forks the format.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const SchemaVersion = 4 // 3 = per-language slot files; 4 = +wallet/solve records (additive)

type State struct {
	SchemaVersion int    `json:"schema_version"`
	Language      string `json:"language"`
	Editor        string `json:"editor"` // chosen editor command, e.g. "code -w"
	Level         string `json:"level"`  // self-report band: never | a-little | regularly
	Stage         string `json:"stage"`  // "intake" | "devliteracy" | "tutorial" | "done"
	DiagIdx       int    `json:"diag_idx"`
	DiagAced      int    `json:"diag_aced"` // deprecated (v1); routing now uses the signal counters
	LessonIdx     int    `json:"lesson_idx"`
	StageIdx      int    `json:"stage_idx"`

	// Step -1 diagnostic outcome (v2): 3-signal routing + placement for Step 0 hints.
	Placement    string `json:"placement"` // "test-out" | "dev-literacy" | "tutorial-full"
	CodingOK     int    `json:"coding_ok"`
	CodingTotal  int    `json:"coding_total"`
	MachineOK    int    `json:"machine_ok"`
	MachineTotal int    `json:"machine_total"`
	SpecOK       int    `json:"spec_ok"`
	SpecTotal    int    `json:"spec_total"`
	PassedAdd    bool   `json:"passed_add"`
	IntakePassed int    `json:"intake_passed"` // questions passed (for the results screen)
	DevIdx       int    `json:"dev_idx"`

	// chosen per-run item sets (so a resumed run replays the SAME variants)
	DiagIDs []string `json:"diag_ids"`
	DevIDs  []string `json:"dev_ids"`

	// Step 0 bench progress
	BenchIDs    []string `json:"bench_ids"`
	BenchIdx    int      `json:"bench_idx"`
	BenchSolved int      `json:"bench_solved"`
	SolvedIDs   []string `json:"solved_ids"` // distinct problems banked across the whole bench
	Step0Done   bool     `json:"step0_done"` // completion milestone reached

	// Track A (v4): hint economy wallet + per-problem solve records. Additive —
	// a v3 save loads with zero values and WalletInit=false triggers the
	// starting grant on first bench entry.
	Tokens            int                    `json:"tokens"`
	WalletInit        bool                   `json:"wallet_init"`
	NudgeCharges      int                    `json:"nudge_charges"`
	NudgeRechargeAt   string                 `json:"nudge_recharge_at"` // RFC3339 accrual anchor
	SolveRecords      map[string]SolveRecord `json:"solve_records,omitempty"`
	MilestonesAwarded []string               `json:"milestones_awarded,omitempty"` // gate categories already paid out

	UpdatedAt string `json:"updated_at"`
}

// SolveRecord is per-problem A1/A2 state, keyed by problem ID alongside
// SolvedIDs. Created on first hint use or first grading attempt. A problem in
// SolvedIDs without WriteupDone is a PROVISIONAL solve (passed the tests, not
// yet explained); WriteupDone makes it fully banked for the graduation gate.
type SolveRecord struct {
	HintTier     int    `json:"hint_tier"`                // max paid tier used: 0 none, 2 strategy, 3 walkthrough (nudges never recorded)
	PityUsed     bool   `json:"pity_used,omitempty"`      // the one-time free strategy hint
	FailedRuns   int    `json:"failed_runs,omitempty"`    // DISTINCT failed attempts (pity eligibility)
	LastFailHash string `json:"last_fail_hash,omitempty"` // fingerprint of the last failing code, so re-running the same code isn't a new failure
	FirstTryAt   string `json:"first_try_at,omitempty"`   // RFC3339 of first grade attempt (pity eligibility)
	WriteupDone  bool   `json:"writeup_done,omitempty"`
	WriteupText  string `json:"writeup_text,omitempty"` // shown to the A4 mentor; never content-graded
	MCQCorrect   bool   `json:"mcq_correct,omitempty"`
}

// Profile is the per-language slot summary shown by profile pickers.
type Profile struct {
	Lang      string `json:"lang"`
	Stage     string `json:"stage"`
	Placement string `json:"placement"`
	Level     string `json:"level"`
	Banked    int    `json:"banked"` // distinct problems banked
	UpdatedAt string `json:"updatedAt"`
}

// Dir returns the directory that holds the save files (exported for sibling
// machine-level config like the mentor selection, which lives next to the
// slots and honors DEVASCENT_SAVE_DIR in tests).
func Dir() (string, error) { return dir() }

// dir returns the directory that holds the save files.
// If DEVASCENT_SAVE_DIR is set, that value is used directly;
// otherwise it falls back to os.UserConfigDir()/DevAscent.
func dir() (string, error) {
	if v := os.Getenv("DEVASCENT_SAVE_DIR"); v != "" {
		return v, nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "DevAscent"), nil
}

func normLang(lang string) string {
	if lang == "" {
		return "python"
	}
	return lang
}

// slotPath returns the save file for one language's slot.
func slotPath(lang string) (string, error) {
	d, err := dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "save-"+normLang(lang)+".json"), nil
}

// migrate moves a legacy single save.json (v2 and earlier) into its
// language's slot. Idempotent; the slot wins if it already exists (the legacy
// file is parked as save.json.bak so nothing is lost and migration stops
// re-triggering).
func migrate() error {
	d, err := dir()
	if err != nil {
		return err
	}
	legacy := filepath.Join(d, "save.json")
	data, err := os.ReadFile(legacy)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		// Corrupt legacy file: park it; slots start fresh.
		return os.Rename(legacy, legacy+".bak")
	}
	target, err := slotPath(s.Language)
	if err != nil {
		return err
	}
	if _, err := os.Stat(target); err == nil {
		return os.Rename(legacy, legacy+".bak")
	}
	return os.Rename(legacy, target)
}

// readState unmarshals one slot file; (nil, nil) when absent.
func readState(p string) (*State, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// LoadLang reads the save slot for one language (migrating any legacy single
// save first). Returns (nil, nil) when the slot doesn't exist.
func LoadLang(lang string) (*State, error) {
	if err := migrate(); err != nil {
		return nil, err
	}
	p, err := slotPath(lang)
	if err != nil {
		return nil, err
	}
	return readState(p)
}

// SaveLang persists s into lang's slot using an atomic write. s.Language is
// forced to the slot's language so a state can never land in the wrong slot.
func SaveLang(lang string, s State) error {
	lang = normLang(lang)
	s.Language = lang
	s.SchemaVersion = SchemaVersion
	s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	d, err := dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	tmp := filepath.Join(d, "save-"+lang+".json.tmp")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	final, err := slotPath(lang)
	if err != nil {
		return err
	}
	return os.Rename(tmp, final)
}

// Profiles lists every language slot, most recently played first.
func Profiles() ([]Profile, error) {
	if err := migrate(); err != nil {
		return nil, err
	}
	d, err := dir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(d)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Profile
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "save-") || !strings.HasSuffix(name, ".json") {
			continue
		}
		s, err := readState(filepath.Join(d, name))
		if err != nil || s == nil {
			continue // a corrupt slot shouldn't hide the others
		}
		lang := strings.TrimSuffix(strings.TrimPrefix(name, "save-"), ".json")
		out = append(out, Profile{
			Lang: lang, Stage: s.Stage, Placement: s.Placement, Level: s.Level,
			Banked: len(s.SolvedIDs), UpdatedAt: s.UpdatedAt,
		})
	}
	// most recent first (RFC3339 sorts lexicographically)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].UpdatedAt > out[j-1].UpdatedAt; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
}

// LoadLatest returns the most recently played slot (nil when no slot exists).
// This is the TUI's resume entry: a single-language player gets exactly the
// old behavior; a multi-language player resumes the last language they played.
func LoadLatest() (*State, error) {
	ps, err := Profiles()
	if err != nil {
		return nil, err
	}
	if len(ps) == 0 {
		return nil, nil
	}
	return LoadLang(ps[0].Lang)
}

// DeleteLang removes one language's slot. Missing file returns nil.
func DeleteLang(lang string) error {
	p, err := slotPath(lang)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
