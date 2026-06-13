package tui

import "devascent/internal/content"

// tui-local view of content. Converted from the loaded content.Catalog so the
// TUI keeps simple value types while content lives in data (YAML).

type lessonStage struct {
	kind  string
	title string
	body  string
	task  *codeTask
}

type lesson struct {
	id, title string
	stages    []lessonStage
}

func diagTask(di content.Diagnostic) codeTask {
	t := codeTask{prompt: di.Prompt}
	if di.Task != nil {
		t.funcName = di.Task.FuncName
		t.code = di.Task.Starter
		t.tests = di.Task.Tests
	}
	return t
}

func toLesson(l content.Lesson) lesson {
	out := lesson{id: l.ID, title: l.Title}
	for _, s := range l.Stages {
		st := lessonStage{kind: s.Kind, title: s.Title, body: s.Body}
		if s.Task != nil {
			st.task = &codeTask{
				prompt:   s.Task.Prompt,
				funcName: s.Task.FuncName,
				code:     s.Task.Starter,
				tests:    s.Task.Tests,
			}
		}
		out.stages = append(out.stages, st)
	}
	return out
}
