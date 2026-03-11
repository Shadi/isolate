package sandbox

import (
	"testing"
)

func TestBuildBwrapArgs(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want []string // subset of args that must appear in order
	}{
		{
			name: "basic command",
			cfg: Config{
				Command: "/bin/echo",
				Args:    []string{"hello"},
			},
			want: []string{
				"--ro-bind", "/usr", "/usr",
				"--proc", "/proc",
				"--dev", "/dev",
				"--die-with-parent",
				"--", "/bin/echo", "hello",
			},
		},
		{
			name: "with volume",
			cfg: Config{
				Command: "/bin/ls",
				Volumes: []Volume{
					{Source: "/home/user/data", Target: "/data"},
				},
			},
			want: []string{"--bind", "/home/user/data", "/data"},
		},
		{
			name: "with read-only volume",
			cfg: Config{
				Command: "/bin/ls",
				Volumes: []Volume{
					{Source: "/home/user/data", Target: "/data", ReadOnly: true},
				},
			},
			want: []string{"--ro-bind", "/home/user/data", "/data"},
		},
		{
			name: "no network",
			cfg: Config{
				Command:   "/bin/echo",
				NoNetwork: true,
			},
			want: []string{"--unshare-net"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildBwrapArgs(&tt.cfg)
			if !containsSubsequence(got, tt.want) {
				t.Errorf("BuildBwrapArgs() = %v, want subsequence %v", got, tt.want)
			}
		})
	}
}

func TestBuildBwrapArgs_commandAtEnd(t *testing.T) {
	cfg := Config{
		Command: "/usr/bin/myapp",
		Args:    []string{"--flag", "value"},
	}
	args := BuildBwrapArgs(&cfg)

	separatorIdx := -1
	for i, a := range args {
		if a == "--" {
			separatorIdx = i
			break
		}
	}
	if separatorIdx == -1 {
		t.Fatal("no -- separator found")
	}
	tail := args[separatorIdx+1:]
	if len(tail) < 3 || tail[0] != "/usr/bin/myapp" || tail[1] != "--flag" || tail[2] != "value" {
		t.Errorf("command not at end: got tail %v", tail)
	}
}

func TestBuildBwrapArgs_bare(t *testing.T) {
	cfg := Config{
		Command: "/myapp",
		Bare:    true,
	}
	args := BuildBwrapArgs(&cfg)

	for _, banned := range []string{"/usr", "/etc", "/proc", "/dev", "/tmp"} {
		if containsElement(args, banned) {
			t.Errorf("bare mode should not mount %s, got args %v", banned, args)
		}
	}

	if !containsSubsequence(args, []string{"--die-with-parent", "--", "/myapp"}) {
		t.Errorf("bare mode missing --die-with-parent or command, got %v", args)
	}
}

func TestBuildBwrapArgs_defaultMountsPresent(t *testing.T) {
	cfg := Config{Command: "/bin/echo"}
	args := BuildBwrapArgs(&cfg)

	for _, expected := range []string{"/usr", "/proc", "/dev", "/tmp"} {
		if !containsElement(args, expected) {
			t.Errorf("default mode should mount %s, got args %v", expected, args)
		}
	}
}

func TestBuildBwrapArgs_defaultSafeEtc(t *testing.T) {
	cfg := Config{Command: "/bin/echo"}
	args := BuildBwrapArgs(&cfg)

	// Safe /etc paths should be present
	safePaths := []string{
		"/etc/resolv.conf",
		"/etc/ssl",
		"/etc/ld.so.cache",
		"/etc/ld.so.conf",
		"/etc/ld.so.conf.d",
		"/etc/nsswitch.conf",
		"/etc/localtime",
	}
	for _, p := range safePaths {
		if !containsElement(args, p) {
			t.Errorf("default mode should mount safe path %s, got args %v", p, args)
		}
	}

	// Full /etc should NOT be mounted
	// Check that "/etc" doesn't appear as a direct mount target (it would appear
	// right after --ro-bind as both source and target)
	for i, a := range args {
		if a == "/etc" && i > 0 && (args[i-1] == "--ro-bind" || args[i-1] == "--bind") {
			t.Errorf("default mode should NOT mount full /etc, got args %v", args)
			break
		}
	}
}

func TestBuildBwrapArgs_fullEtc(t *testing.T) {
	cfg := Config{
		Command: "/bin/echo",
		FullEtc: true,
	}
	args := BuildBwrapArgs(&cfg)

	// Full /etc should be mounted
	if !containsSubsequence(args, []string{"--ro-bind", "/etc", "/etc"}) {
		t.Errorf("full-etc mode should mount /etc, got args %v", args)
	}
}

func TestBuildBwrapArgs_bareIgnoresFullEtc(t *testing.T) {
	cfg := Config{
		Command: "/bin/echo",
		Bare:    true,
		FullEtc: true,
	}
	args := BuildBwrapArgs(&cfg)

	if containsElement(args, "/etc") {
		t.Errorf("bare mode should override full-etc, got args %v", args)
	}
}
