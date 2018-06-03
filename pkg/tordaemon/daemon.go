package tordaemon

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Tor struct {
	cmd *exec.Cmd
	ctx context.Context
}

func (t *Tor) Start(ctx context.Context) {
	fmt.Println("starting tor...")

	t.ctx = ctx

	t.cmd = exec.CommandContext(ctx, "tor", "-f", "/run/tor/torfile")
	t.cmd.Stdout = os.Stdout
	t.cmd.Stderr = os.Stderr

	err := t.cmd.Start()
	if err != nil {
		fmt.Print(err)
		return
	}
}

func (t *Tor) Reload() {
	fmt.Println("reloading tor...")

	t.cmd.Process.Signal(syscall.SIGHUP)
}
