package upgrade

import (
	"fmt"
	"os"
	"os/exec"
)

func Run() error {
	brewPath, err := exec.LookPath("brew")
	if err != nil {
		return fmt.Errorf("brew not found. Install Homebrew from https://brew.sh")
	}

	cmd := exec.Command(brewPath, "upgrade", "CircleCI-Public/circleci/chunk")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("upgrade failed: %w", err)
	}

	return nil
}
