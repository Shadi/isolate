package sandbox

import (
	"testing"
)

func TestBuildSystemdRunArgs(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want []string
	}{
		{
			name: "memory limit only",
			cfg: Config{
				MaxMemory: "256Mi",
				Command:   "/bin/echo",
			},
			want: []string{
				"systemd-run", "--user", "--scope",
				"-p", "MemoryMax=268435456",
			},
		},
		{
			name: "cpu limit only",
			cfg: Config{
				MaxCPU:  "1",
				Command: "/bin/echo",
			},
			want: []string{
				"systemd-run", "--user", "--scope",
				"-p", "CPUQuota=100%",
			},
		},
		{
			name: "half cpu",
			cfg: Config{
				MaxCPU:  "500m",
				Command: "/bin/echo",
			},
			want: []string{
				"-p", "CPUQuota=50%",
			},
		},
		{
			name: "memory and cpu",
			cfg: Config{
				MaxMemory: "512Mi",
				MaxCPU:    "2",
				Command:   "/bin/echo",
			},
			want: []string{
				"-p", "MemoryMax=536870912",
				"-p", "CPUQuota=200%",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildSystemdRunPrefix(&tt.cfg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !containsSubsequence(got, tt.want) {
				t.Errorf("BuildSystemdRunPrefix() = %v, want subsequence %v", got, tt.want)
			}
		})
	}
}

func TestBuildSystemdRunArgs_noLimits(t *testing.T) {
	cfg := Config{Command: "/bin/echo"}
	got, err := BuildSystemdRunPrefix(&cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil prefix when no limits, got %v", got)
	}
}
