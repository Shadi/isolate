package sandbox

// Config holds all sandbox configuration.
type Config struct {
	Command    string
	Args       []string
	MaxMemory  string   // e.g. "256Mi", "1Gi"
	MaxCPU     string   // e.g. "1", "500m", "0.5"
	Volumes    []Volume // bind mounts
	NoNetwork  bool     // disable networking
	Bare       bool     // skip all default system mounts
	FullEtc    bool     // mount all of /etc instead of safe subset
	MountCwd   bool     // auto-mount current working directory
	WorkingDir string
	Env        []string
}

// Volume represents a bind mount.
type Volume struct {
	Source   string
	Target   string
	ReadOnly bool
}
