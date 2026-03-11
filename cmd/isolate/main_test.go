package main

import "testing"

func TestParseFlags(t *testing.T) {
	args := []string{
		"--max-memory", "256Mi",
		"--max-cpu", "1",
		"--volume", "/app:/home/user/mydir",
		"--volume", "/data:/mnt/data:ro",
		"--no-network",
		"--workdir", "/app",
		"myapp", "--flag", "value",
	}

	cfg, err := parseFlags(args)
	if err != nil {
		t.Fatalf("parseFlags() error: %v", err)
	}

	if cfg.MaxMemory != "256Mi" {
		t.Errorf("MaxMemory = %q, want 256Mi", cfg.MaxMemory)
	}
	if cfg.MaxCPU != "1" {
		t.Errorf("MaxCPU = %q, want 1", cfg.MaxCPU)
	}
	if cfg.Command != "myapp" {
		t.Errorf("Command = %q, want myapp", cfg.Command)
	}
	if len(cfg.Args) != 2 || cfg.Args[0] != "--flag" || cfg.Args[1] != "value" {
		t.Errorf("Args = %v, want [--flag value]", cfg.Args)
	}
	if !cfg.NoNetwork {
		t.Error("NoNetwork = false, want true")
	}
	if cfg.WorkingDir != "/app" {
		t.Errorf("WorkingDir = %q, want /app", cfg.WorkingDir)
	}
	if len(cfg.Volumes) != 2 {
		t.Fatalf("Volumes len = %d, want 2", len(cfg.Volumes))
	}
	if cfg.Volumes[0].Source != "/home/user/mydir" || cfg.Volumes[0].Target != "/app" {
		t.Errorf("Volumes[0] = %+v", cfg.Volumes[0])
	}
	if cfg.Volumes[1].Source != "/mnt/data" || cfg.Volumes[1].Target != "/data" || !cfg.Volumes[1].ReadOnly {
		t.Errorf("Volumes[1] = %+v", cfg.Volumes[1])
	}
}

func TestParseFlags_shortVolume(t *testing.T) {
	args := []string{
		"-v", "/app:/home/user/mydir",
		"myapp",
	}

	cfg, err := parseFlags(args)
	if err != nil {
		t.Fatalf("parseFlags() error: %v", err)
	}

	if len(cfg.Volumes) != 1 {
		t.Fatalf("Volumes len = %d, want 1", len(cfg.Volumes))
	}
	if cfg.Volumes[0].Source != "/home/user/mydir" || cfg.Volumes[0].Target != "/app" {
		t.Errorf("Volumes[0] = %+v", cfg.Volumes[0])
	}
}

func TestParseFlags_noCommand(t *testing.T) {
	_, err := parseFlags([]string{"--max-memory", "256Mi"})
	if err == nil {
		t.Error("expected error for missing command")
	}
}

func TestParseFlags_minimal(t *testing.T) {
	cfg, err := parseFlags([]string{"/bin/echo", "hello"})
	if err != nil {
		t.Fatalf("parseFlags() error: %v", err)
	}
	if cfg.Command != "/bin/echo" {
		t.Errorf("Command = %q, want /bin/echo", cfg.Command)
	}
	if len(cfg.Args) != 1 || cfg.Args[0] != "hello" {
		t.Errorf("Args = %v, want [hello]", cfg.Args)
	}
}
