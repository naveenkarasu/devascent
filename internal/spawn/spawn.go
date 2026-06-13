package spawn

// Package spawn marks child processes as windowless. On Windows, a
// GUI-subsystem parent (the Wails app) gets a visible console window flashed
// for every console child it starts — toolchain probes at launch and on
// language switch, compiler/runner invocations on every grade, mentor CLIs.
// Hide suppresses that. Never use it on interactive children that need a
// console of their own (the TUI's external-editor launch).

import "os/exec"

// Hide marks cmd to run without a console window (no-op off Windows).
func Hide(cmd *exec.Cmd) { hide(cmd) }
