package content

import "testing"

func TestLoad(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Diagnostics) == 0 {
		t.Fatal("no diagnostics loaded")
	}
	if len(c.Lessons) < 10 {
		t.Fatalf("expected at least 10 lessons, got %d", len(c.Lessons))
	}
	// Lessons must be ordered.
	for i := 1; i < len(c.Lessons); i++ {
		if c.Lessons[i].Order < c.Lessons[i-1].Order {
			t.Fatalf("lessons not ordered at %d", i)
		}
	}
	// The counted intake opens with a coding-floor item (lowest slot order).
	if c.Diagnostics[0].Slot != "coding-floor" {
		t.Fatalf("expected a coding-floor item first by order, got slot %q (%q)", c.Diagnostics[0].Slot, c.Diagnostics[0].ID)
	}
	// The (uncounted) warm-up intro is loaded separately.
	if len(c.Intro) == 0 {
		t.Fatal("expected at least one intro warm-up item")
	}
	if len(c.DevTasks) < 5 {
		t.Fatalf("expected at least 5 dev-literacy tasks, got %d", len(c.DevTasks))
	}
	l, ok := c.LessonByID("functions")
	if !ok {
		t.Fatal("functions lesson not found")
	}
	if len(l.Stages) != 3 {
		t.Fatalf("expected 3 stages, got %d", len(l.Stages))
	}
	if l.Stages[1].Task == nil || len(l.Stages[1].Task.Tests) == 0 {
		t.Fatal("we_do stage missing task/tests")
	}
}
