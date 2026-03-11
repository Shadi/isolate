package sandbox

import "fmt"

// BuildSystemdRunPrefix builds the systemd-run command prefix for cgroup resource limits.
// Returns nil if no limits are configured.
func BuildSystemdRunPrefix(cfg *Config) ([]string, error) {
	if cfg.MaxMemory == "" && cfg.MaxCPU == "" {
		return nil, nil
	}

	args := []string{"systemd-run", "--user", "--scope"}

	if cfg.MaxMemory != "" {
		bytes, err := ParseMemory(cfg.MaxMemory)
		if err != nil {
			return nil, fmt.Errorf("invalid memory limit: %w", err)
		}
		args = append(args, "-p", fmt.Sprintf("MemoryMax=%d", bytes))
	}

	if cfg.MaxCPU != "" {
		cores, err := ParseCPU(cfg.MaxCPU)
		if err != nil {
			return nil, fmt.Errorf("invalid cpu limit: %w", err)
		}
		percent := int(cores * 100)
		args = append(args, "-p", fmt.Sprintf("CPUQuota=%d%%", percent))
	}

	// Separator: everything after "--" is the actual command
	args = append(args, "--")

	return args, nil
}
