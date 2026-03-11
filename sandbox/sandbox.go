package sandbox

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// Sandbox manages an isolated process.
type Sandbox struct {
	cfg *Config
}

// New creates a new Sandbox from the given config.
func New(cfg *Config) (*Sandbox, error) {
	if cfg.Command == "" {
		return nil, fmt.Errorf("command is required")
	}

	return &Sandbox{cfg: cfg}, nil
}

// buildCommand assembles the full exec.Cmd.
func (s *Sandbox) buildCommand(ctx context.Context) (*exec.Cmd, error) {
	args, err := s.buildArgs()
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Env = s.buildEnv()

	return cmd, nil
}

// buildArgs assembles the full argument list: [systemd-run ...] bwrap ...
func (s *Sandbox) buildArgs() ([]string, error) {
	prefix, err := BuildSystemdRunPrefix(s.cfg)
	if err != nil {
		return nil, err
	}

	bwrapCmd := append([]string{"bwrap"}, BuildBwrapArgs(s.cfg)...)

	if prefix != nil {
		return append(prefix, bwrapCmd...), nil
	}
	return bwrapCmd, nil
}

// buildEnv merges the current environment with config env vars.
func (s *Sandbox) buildEnv() []string {
	env := os.Environ()
	env = append(env, s.cfg.Env...)
	return env
}

// Run executes the sandboxed command, connecting stdin/stdout/stderr.
func (s *Sandbox) Run(ctx context.Context) error {
	cmd, err := s.buildCommand(ctx)
	if err != nil {
		return err
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
