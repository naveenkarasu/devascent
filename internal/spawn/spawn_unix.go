//go:build !windows

package spawn

import "os/exec"

func hide(*exec.Cmd) {} // consoles aren't windows off Windows
