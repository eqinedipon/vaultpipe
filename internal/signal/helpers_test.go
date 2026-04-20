package signal_test

import (
	"os/exec"
	"runtime"
)

// testCmd returns a long-running no-op command suitable for use as a
// target process in signal forwarding tests.
func testCmd() *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", "ping -n 100 127.0.0.1 >nul")
	}
	return exec.Command("sleep", "30")
}
