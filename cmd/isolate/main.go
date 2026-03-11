package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"

	"github.com/shadi/isolate/sandbox"
)

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	cfg, err := parseFlags(os.Args[1:])
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse flags")
	}

	log.Info().
		Str("command", cfg.Command).
		Str("memory", cfg.MaxMemory).
		Str("cpu", cfg.MaxCPU).
		Int("volumes", len(cfg.Volumes)).
		Bool("no_network", cfg.NoNetwork).
		Msg("starting sandbox")

	s, err := sandbox.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create sandbox")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := s.Run(ctx); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		log.Fatal().Err(err).Msg("sandbox failed")
	}
}

func parseFlags(args []string) (*sandbox.Config, error) {
	fs := flag.NewFlagSet("isolate", flag.ContinueOnError)

	maxMemory := fs.String("max-memory", "", "Memory limit (e.g., 256Mi, 1Gi)")
	maxCPU := fs.String("max-cpu", "", "CPU limit (e.g., 1, 0.5, 500m)")
	volumes := fs.StringArrayP("volume", "v", nil, "Bind mount target:source[:ro] (repeatable)")
	noNetwork := fs.Bool("no-network", false, "Disable network access")
	bare := fs.Bool("bare", false, "Skip all default system mounts")
	fullEtc := fs.Bool("full-etc", false, "Mount all of /etc (includes secrets; default mounts only safe paths)")
	workDir := fs.String("workdir", "", "Working directory inside sandbox")

	fs.SetInterspersed(false)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	remaining := fs.Args()
	if len(remaining) == 0 {
		return nil, errors.New("usage: isolate [flags] command [args...]")
	}

	cfg := &sandbox.Config{
		Command:    remaining[0],
		Args:       remaining[1:],
		MaxMemory:  *maxMemory,
		MaxCPU:     *maxCPU,
		NoNetwork:  *noNetwork,
		Bare:       *bare,
		FullEtc:    *fullEtc,
		WorkingDir: *workDir,
	}

	for _, v := range *volumes {
		vol, err := sandbox.ParseVolume(v)
		if err != nil {
			return nil, err
		}
		cfg.Volumes = append(cfg.Volumes, vol)
	}

	return cfg, nil
}
