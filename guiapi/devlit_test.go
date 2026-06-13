package guiapi

import "testing"

// The dev-literacy session: a wrong answer stays on the same task with the
// hint; correct answers (built from each task's own accept forms) walk the
// whole set to a Done step with a full score.
func TestDevLiteracySession(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir()) // never touch the real save
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	d := e.StartDevLiteracy()
	first := d.Step()
	if first.Done || first.Total == 0 {
		t.Fatalf("empty dev-literacy session: %+v", first)
	}

	// Wrong answer: no advance, hint surfaced.
	out := d.Submit("definitely-not-a-command --nope")
	if out.Passed {
		t.Fatal("nonsense command unexpectedly passed")
	}
	if out.Next.Done || out.Next.Index != first.Index {
		t.Fatalf("wrong answer advanced the session: %+v", out.Next)
	}

	// Correct answers from each task's own accepted forms.
	for i := 0; i < len(d.tasks); i++ {
		task := d.tasks[d.idx]
		ans := ""
		if len(task.Accept) > 0 {
			ans = task.Accept[0]
		} else if len(task.Commands) > 0 {
			ans = task.Commands[0]
			for _, fl := range task.Flags {
				ans += " " + fl
			}
		}
		out = d.Submit(ans)
		if !out.Passed {
			t.Fatalf("task %s rejected its own accepted answer %q", task.ID, ans)
		}
	}
	final := d.Step()
	if !final.Done || final.Passed != final.Total {
		t.Fatalf("expected full-score done, got %+v", final)
	}
}
