package main

import (
	"context"
	"fmt"
	"os/exec"
)

var format = "Startup finished in %fs (firmware) + %fs (loader) + %fs (kernel) + %fs (userspace)"

type bootTime struct {
	Firmware float64 `json:"-"`
	Loader   float64 `json:"-"`
	Kernel   float64 `json:"kernel"`
	User     float64 `json:"user"`
}

func (b bootTime) String() string {
	return fmt.Sprintf("kernel: %.4f, user: %.4f, total: %.4f", b.Kernel, b.User, b.Kernel+b.User)
}

// readBootTime runs systemd-analyze to get system boot times and returns them
// in a bootTime struct
func readBootTime(ctx context.Context) (bootTime, error) {
	cmd := exec.CommandContext(ctx, "systemd-analyze")
	output, err := cmd.Output()
	if err != nil {
		return bootTime{}, err
	}

	var bt bootTime
	n, err := fmt.Sscanf(string(output), format, &bt.Firmware, &bt.Loader, &bt.Kernel, &bt.User)
	if n != 4 {
		return bootTime{}, fmt.Errorf("failed to parse '%s'", string(output))
	}
	if err != nil {
		return bootTime{}, err
	}
	return bt, nil
}
