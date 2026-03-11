package sandbox

import (
	"context"
	"errors"
	"os/exec"
	"testing"
)

func TestSandbox_buildFullCommand(t *testing.T) {
	cfg := Config{
		Command:   "/bin/echo",
		Args:      []string{"hello"},
		MaxMemory: "256Mi",
		MaxCPU:    "1",
		NoNetwork: true,
		Volumes: []Volume{
			{Source: "/home/user/app", Target: "/app"},
		},
	}

	s, err := New(&cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	cmd, err := s.buildCommand(context.Background())
	if err != nil {
		t.Fatalf("buildCommand() error: %v", err)
	}

	if cmd.Args[0] != "systemd-run" {
		t.Errorf("expected systemd-run prefix, got %v", cmd.Args[0])
	}

	if !containsElement(cmd.Args, "bwrap") {
		t.Error("expected bwrap in command args")
	}
}

func TestSandbox_buildFullCommand_noCgroupLimits(t *testing.T) {
	cfg := Config{
		Command: "/bin/echo",
		Args:    []string{"hi"},
	}

	s, err := New(&cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	cmd, err := s.buildCommand(context.Background())
	if err != nil {
		t.Fatalf("buildCommand() error: %v", err)
	}

	if cmd.Args[0] != "bwrap" {
		t.Errorf("expected bwrap as first arg, got %v", cmd.Args[0])
	}
}

func TestSandbox_Run_propagatesExitCode(t *testing.T) {
	cfg := Config{Command: "/usr/bin/false"}

	s, err := New(&cfg)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	runErr := s.Run(context.Background())
	if runErr == nil {
		t.Fatal("expected error from Run, got nil")
	}

	var exitErr *exec.ExitError
	if !errors.As(runErr, &exitErr) {
		t.Fatalf("expected *exec.ExitError, got %T: %v", runErr, runErr)
	}

	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
}
