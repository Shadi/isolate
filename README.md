# isolate

A Go CLI and library that runs commands in a sandboxed environment using [bubblewrap](https://github.com/containers/bubblewrap), and resource limits using cgroups.

I wanted to have docker-like isolation when running some binaries that I don't trust, or when I wanted to make sure that they don't have network access,
bwrap directly seemed too low-level, and also I couldn't set resource limits on the app, this package use bwrap for sandboxing, and cgroups(via systemd-run) for resource limits.

## Install

```bash
go install github.com/shadi/isolate/cmd/isolate@latest
```

Requires `bwrap` and `systemd-run` on the host.

## Usage

```
isolate [flags] command [args...]
```

## Flags

| Flag | Description |
|------|-------------|
| `--max-memory` | Memory limit using k8s format (`256Mi`, `1Gi`, `512M`) — enforced via cgroups |
| `--max-cpu` | CPU limit in cores or millicores (`1`, `0.5`, `500m`) — enforced via cgroups |
| `-v`, `--volume` | Bind mount as `target:source[:ro]` (repeatable) |
| `--no-network` | Disable network access |
| `--full-etc` | Mount all of `/etc` (default only mounts safe paths like `resolv.conf`, `ssl/`, `hosts`) |
| `--bare` | Skip all default mounts — you control everything via `--volume` |
| `--workdir` | Working directory inside the sandbox |

## Examples

Run a command with memory and CPU limits:
```bash
isolate --max-memory 256Mi --max-cpu 2 -v /app:/home/user/myproject /usr/bin/python3 train.py
```

Mount a directory and restrict network:
```bash
isolate --no-network -v /app:/home/user/myproject /bin/ls /app
```

Full `/etc` access (exposes secrets — use only when needed):
```bash
isolate --full-etc /usr/bin/curl https://example.com
```

Bare mode with manual mounts:
```bash
isolate --bare -v /usr:/usr:ro -v /app:/home/user/code /usr/bin/myapp
```

## Default mounts

Unless `--bare` is set, the sandbox includes:

- `/usr` (read-only), `/lib`, `/bin`, `/sbin` (symlinks into `/usr`)
- `/proc`, `/dev`, `/tmp`
- Safe `/etc` subset: `resolv.conf`, `ssl/`, `ca-certificates/`, `ld.so.*`, `nsswitch.conf`, `localtime`, `hosts`
- systemd-resolved socket (for DNS)

## Library usage

```go
import "github.com/shadi/isolate/sandbox"

s, _ := sandbox.New(&sandbox.Config{
    Command:   "/usr/bin/python3",
    Args:      []string{"train.py"},
    MaxMemory: "256Mi",
    MaxCPU:    "500m",
    NoNetwork: true,
    Volumes:   []sandbox.Volume{{Source: "/home/user/ml", Target: "/workspace"}},
})
s.Run(context.Background())
```

## How it works

| Resource | Mechanism |
|----------|-----------|
| Filesystem isolation | bubblewrap namespaces |
| Memory limit | `systemd-run -p MemoryMax=...` (cgroups v2) |
| CPU limit | `systemd-run -p CPUQuota=...` (cgroups v2) |
| Network isolation | `bwrap --unshare-net` |
