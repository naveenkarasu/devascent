package save

import (
	"testing"
)

func TestLoad_NoFile(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v; want nil", err)
	}
	if got != nil {
		t.Fatalf("Load() = %v; want nil", got)
	}
}

func TestRoundTrip(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())

	in := State{
		Language:  "python",
		Stage:     "tutorial",
		LessonIdx: 3,
		StageIdx:  1,
		DiagAced:  2,
	}
	if err := Save(in); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	out, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if out == nil {
		t.Fatal("Load() returned nil; want non-nil State")
	}

	if out.Language != in.Language {
		t.Errorf("Language = %q; want %q", out.Language, in.Language)
	}
	if out.Stage != in.Stage {
		t.Errorf("Stage = %q; want %q", out.Stage, in.Stage)
	}
	if out.LessonIdx != in.LessonIdx {
		t.Errorf("LessonIdx = %d; want %d", out.LessonIdx, in.LessonIdx)
	}
	if out.StageIdx != in.StageIdx {
		t.Errorf("StageIdx = %d; want %d", out.StageIdx, in.StageIdx)
	}
	if out.DiagAced != in.DiagAced {
		t.Errorf("DiagAced = %d; want %d", out.DiagAced, in.DiagAced)
	}
	if out.SchemaVersion != SchemaVersion {
		t.Errorf("SchemaVersion = %d; want %d", out.SchemaVersion, SchemaVersion)
	}
	if out.UpdatedAt == "" {
		t.Error("UpdatedAt is empty; want non-empty")
	}
}

func TestDelete(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir())

	// Save a file so there is something to delete.
	if err := Save(State{Language: "python", Stage: "intake"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Delete it.
	if err := Delete(); err != nil {
		t.Fatalf("Delete() error = %v; want nil", err)
	}

	// Subsequent Load should return nil, nil.
	got, err := Load()
	if err != nil {
		t.Fatalf("Load() after Delete error = %v; want nil", err)
	}
	if got != nil {
		t.Fatalf("Load() after Delete = %v; want nil", got)
	}

	// Delete on a missing file should return nil.
	if err := Delete(); err != nil {
		t.Fatalf("Delete() on missing file error = %v; want nil", err)
	}
}
