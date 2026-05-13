package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/CircleCI-Public/chunk-cli/internal/iostream"
	"github.com/CircleCI-Public/chunk-cli/internal/validate"
)

func newHooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Manage chunk hook execution",
	}
	cmd.AddCommand(newHooksDisableCmd())
	cmd.AddCommand(newHooksEnableCmd())
	cmd.AddCommand(newHooksStatusCmd())
	return cmd
}

func hooksProjectRoot() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err == nil {
		if s := strings.TrimSpace(string(out)); s != "" {
			return s
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}

func hooksSentinelPath() string {
	return filepath.Join(hooksProjectRoot(), ".chunk", "hooks-disabled")
}

func newHooksDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "disable",
		Short:        "Disable chunk validate hooks",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := hooksSentinelPath()
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				return fmt.Errorf("create .chunk directory: %w", err)
			}
			if err := os.WriteFile(p, []byte{}, 0o644); err != nil {
				return fmt.Errorf("create hooks-disabled sentinel: %w", err)
			}
			streams := iostream.FromCmd(cmd)
			streams.ErrPrintln("Hooks disabled. Run 'chunk hooks enable' to re-enable.")
			return nil
		},
	}
}

func newHooksEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "enable",
		Short:        "Re-enable chunk validate hooks",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := hooksSentinelPath()
			if err := os.Remove(p); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("remove hooks-disabled sentinel: %w", err)
			}
			streams := iostream.FromCmd(cmd)
			streams.ErrPrintln("Hooks enabled.")
			return nil
		},
	}
}

func newHooksStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "status",
		Short:        "Show whether hooks are enabled or disabled",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			streams := iostream.FromCmd(cmd)
			root := hooksProjectRoot()
			if validate.HooksDisabled(root) {
				streams.ErrPrintln("disabled")
			} else {
				streams.ErrPrintln("enabled")
			}
			return nil
		},
	}
}
