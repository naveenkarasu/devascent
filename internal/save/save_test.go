package save

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadLang_NoFile(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	got, err := LoadLang("python")
	if err != nil {
		t.Fatalf("LoadLang() error = %v; want nil", err)
	}
	if got != nil {
		t.Fatalf("LoadLang() = %v; want nil", got)
	}
}

func TestRoundTripPerLang(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())

	if err := SaveLang("python", State{Stage: "tutorial", LessonIdx: 3, StageIdx: 1}); err != nil {
		t.Fatalf("SaveLang(python) error = %v", err)
	}
	if err := SaveLang("go", State{Stage: "bench", SolvedIDs: []string{"nc-two-sum"}}); err != nil {
		t.Fatalf("SaveLang(go) error = %v", err)
	}

	py, err := LoadLang("python")
	if err != nil || py == nil {
		t.Fatalf("LoadLang(python) = %v, %v", py, err)
	}
	if py.Language != "python" || py.Stage != "tutorial" || py.LessonIdx != 3 {
		t.Errorf("python slot wrong: %+v", py)
	}
	if py.SchemaVersion != SchemaVersion || py.UpdatedAt == "" {
		t.Errorf("version/timestamp not stamped: %+v", py)
	}

	goSt, err := LoadLang("go")
	if err != nil || goSt == nil {
		t.Fatalf("LoadLang(go) = %v, %v", goSt, err)
	}
	// Slots are independent: go's solve never leaks into python's slot.
	if goSt.Stage != "bench" || len(goSt.SolvedIDs) != 1 || len(py.SolvedIDs) != 0 {
		t.Errorf("slot isolation broken: go=%+v py=%+v", goSt, py)
	}
}

func TestSaveLangForcesSlotLanguage(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	// A state claiming another language still lands in (and is stamped with)
	// the slot it was saved to.
	if err := SaveLang("rust", State{Language: "python", Stage: "intake"}); err != nil {
		t.Fatal(err)
	}
	s, err := LoadLang("rust")
	if err != nil || s == nil {
		t.Fatalf("LoadLang(rust) = %v, %v", s, err)
	}
	if s.Language != "rust" {
		t.Errorf("Language = %q; want rust", s.Language)
	}
}

func TestMigrateLegacySingleSave(t *testing.T) {
	dirp := t.TempDir()
	t.Setenv("DEVASCENT_SAVE_DIR", dirp)

	// Write a v2-era single save.json by hand.
	legacy := State{SchemaVersion: 2, Language: "csharp", Stage: "done", Placement: "tutorial-full",
		SolvedIDs: []string{"a", "b"}, UpdatedAt: time.Now().UTC().Format(time.RFC3339)}
	data, _ := json.Marshal(legacy)
	if err := os.WriteFile(filepath.Join(dirp, "save.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}

	// First access migrates it into the csharp slot.
	s, err := LoadLang("csharp")
	if err != nil || s == nil {
		t.Fatalf("LoadLang(csharp) after migration = %v, %v", s, err)
	}
	if s.Placement != "tutorial-full" || len(s.SolvedIDs) != 2 {
		t.Errorf("migrated state wrong: %+v", s)
	}
	if _, err := os.Stat(filepath.Join(dirp, "save.json")); !os.IsNotExist(err) {
		t.Error("legacy save.json still present after migration")
	}
	// Other slots stay empty.
	if other, _ := LoadLang("python"); other != nil {
		t.Errorf("python slot unexpectedly exists: %+v", other)
	}
}

func TestProfilesAndLatest(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if ps, err := Profiles(); err != nil || len(ps) != 0 {
		t.Fatalf("empty dir Profiles = %v, %v", ps, err)
	}
	if s, err := LoadLatest(); err != nil || s != nil {
		t.Fatalf("empty dir LoadLatest = %v, %v", s, err)
	}

	if err := SaveLang("python", State{Stage: "done", SolvedIDs: []string{"x"}}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1100 * time.Millisecond) // RFC3339 second resolution: force a later stamp
	if err := SaveLang("go", State{Stage: "bench"}); err != nil {
		t.Fatal(err)
	}

	ps, err := Profiles()
	if err != nil || len(ps) != 2 {
		t.Fatalf("Profiles = %v, %v; want 2", ps, err)
	}
	if ps[0].Lang != "go" || ps[1].Lang != "python" {
		t.Errorf("profiles not most-recent-first: %+v", ps)
	}
	if ps[1].Banked != 1 {
		t.Errorf("python banked = %d; want 1", ps[1].Banked)
	}

	latest, err := LoadLatest()
	if err != nil || latest == nil || latest.Language != "go" {
		t.Fatalf("LoadLatest = %+v, %v; want the go slot", latest, err)
	}
}

func TestDeleteLang(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())
	if err := SaveLang("python", State{Stage: "intake"}); err != nil {
		t.Fatal(err)
	}
	if err := DeleteLang("python"); err != nil {
		t.Fatalf("DeleteLang error = %v", err)
	}
	if got, err := LoadLang("python"); err != nil || got != nil {
		t.Fatalf("LoadLang after delete = %v, %v", got, err)
	}
	if err := DeleteLang("python"); err != nil {
		t.Fatalf("DeleteLang on missing slot error = %v", err)
	}
}
