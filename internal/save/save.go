package save

// Package save persists Step -1 progress as JSON (atomic writes) under the
// per-OS user config dir, so a long run resumes across sessions.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const SchemaVersion = 2

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

	UpdatedAt string `json:"updated_at"`
}

// dir returns the directory that holds the save file.
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

// Path returns the full path to the save file.
func Path() (string, error) {
	d, err := dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "save.json"), nil
}

// Load reads and unmarshals the save file.
// If the file does not exist, it returns (nil, nil).
func Load() (*State, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
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

// Save persists s to disk using an atomic write (tmp file + rename).
func Save(s State) error {
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

	tmp := filepath.Join(d, "save.json.tmp")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}

	final := filepath.Join(d, "save.json")
	return os.Rename(tmp, final)
}

// Delete removes the save file. If the file does not exist, nil is returned.
func Delete() error {
	p, err := Path()
	if err != nil {
		return err
	}
	err = os.Remove(p)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}
