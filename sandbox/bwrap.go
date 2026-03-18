package sandbox

import (
	"fmt"
	"os"
)

// mount represents a single bwrap mount directive.
type mount struct {
	flag   string // e.g. "--ro-bind", "--symlink", "--proc", "--tmpfs", "--ro-bind-try"
	source string
	target string // empty for single-arg flags like --tmpfs
}

// baseMounts are always applied unless --bare is set.
var baseMounts = []mount{
	{"--ro-bind", "/usr", "/usr"},
	{"--symlink", "usr/lib", "/lib"},
	{"--symlink", "usr/lib64", "/lib64"},
	{"--symlink", "usr/bin", "/bin"},
	{"--symlink", "usr/sbin", "/sbin"},
	{"--proc", "/proc", ""},
	{"--dev", "/dev", ""},
	{"--tmpfs", "/tmp", ""},
}

// safeEtcMounts expose only the /etc paths needed for programs to function.
// These contain no secrets or sensitive system configuration.
var safeEtcMounts = []mount{
	{"--ro-bind-try", "/etc/resolv.conf", "/etc/resolv.conf"},
	{"--ro-bind-try", "/etc/ssl", "/etc/ssl"},
	{"--ro-bind-try", "/etc/ca-certificates", "/etc/ca-certificates"},
	{"--ro-bind-try", "/etc/ld.so.cache", "/etc/ld.so.cache"},
	{"--ro-bind-try", "/etc/ld.so.conf", "/etc/ld.so.conf"},
	{"--ro-bind-try", "/etc/ld.so.conf.d", "/etc/ld.so.conf.d"},
	{"--ro-bind-try", "/etc/nsswitch.conf", "/etc/nsswitch.conf"},
	{"--ro-bind-try", "/etc/localtime", "/etc/localtime"},
	{"--ro-bind-try", "/etc/hosts", "/etc/hosts"},
	{"--ro-bind-try", "/run/systemd/resolve", "/run/systemd/resolve"},
}

// fullEtcMounts expose all of /etc (including secrets like shadow, ssh keys, etc).
var fullEtcMounts = []mount{
	{"--ro-bind", "/etc", "/etc"},
	{"--ro-bind-try", "/run/systemd/resolve", "/run/systemd/resolve"},
}

func (m mount) args() []string {
	if m.target == "" {
		return []string{m.flag, m.source}
	}
	return []string{m.flag, m.source, m.target}
}

func appendMounts(args []string, mounts []mount) []string {
	for _, m := range mounts {
		args = append(args, m.args()...)
	}
	return args
}

// BuildBwrapArgs constructs the argument list for bwrap.
func BuildBwrapArgs(cfg *Config) ([]string, error) {
	var args []string

	if !cfg.Bare {
		args = appendMounts(args, baseMounts)

		if cfg.FullEtc {
			args = appendMounts(args, fullEtcMounts)
		} else {
			args = appendMounts(args, safeEtcMounts)
		}
	}

	args = append(args, "--die-with-parent")

	if !cfg.Bare {
		args = append(args, "--unshare-pid")
	}

	if cfg.NoNetwork {
		args = append(args, "--unshare-net")
	}

	if cfg.MountCwd {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("getting current directory: %w", err)
		}
		args = append(args, "--bind", cwd, cwd)
		if cfg.WorkingDir == "" {
			args = append(args, "--chdir", cwd)
		}
	}

	for _, v := range cfg.Volumes {
		if v.ReadOnly {
			args = append(args, "--ro-bind", v.Source, v.Target)
		} else {
			args = append(args, "--bind", v.Source, v.Target)
		}
	}

	if cfg.WorkingDir != "" {
		args = append(args, "--chdir", cfg.WorkingDir)
	}

	args = append(args, "--")
	args = append(args, cfg.Command)
	args = append(args, cfg.Args...)

	return args, nil
}
