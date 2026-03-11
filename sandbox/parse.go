package sandbox

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
)

// ParseVolume parses a volume string in the format "target:source" or "target:source:ro".
func ParseVolume(s string) (Volume, error) {
	if s == "" {
		return Volume{}, fmt.Errorf("empty volume string")
	}

	parts := strings.Split(s, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return Volume{}, fmt.Errorf("invalid volume format %q: expected target:source[:ro]", s)
	}

	v := Volume{
		Target: parts[0],
		Source: parts[1],
	}

	if len(parts) == 3 {
		if parts[2] == "ro" {
			v.ReadOnly = true
		} else {
			return Volume{}, fmt.Errorf("invalid volume option %q: expected 'ro'", parts[2])
		}
	}

	return v, nil
}

func parseQuantity(s, label string) (resource.Quantity, error) {
	if s == "" {
		return resource.Quantity{}, fmt.Errorf("empty %s string", label)
	}

	q, err := resource.ParseQuantity(s)
	if err != nil {
		return resource.Quantity{}, fmt.Errorf("invalid %s value %q: %w", label, s, err)
	}

	return q, nil
}

// ParseMemory parses a memory string using Kubernetes resource.Quantity format.
// Supports: "256Mi", "1Gi", "512Ki", "256M", "1G", or plain bytes "1048576".
func ParseMemory(s string) (uint64, error) {
	q, err := parseQuantity(s, "memory")
	if err != nil {
		return 0, err
	}

	val := q.Value()
	if val <= 0 {
		return 0, fmt.Errorf("memory must be positive, got %d", val)
	}

	return uint64(val), nil
}

// ParseCPU parses a CPU limit string using Kubernetes resource.Quantity format.
// Supports: "1", "0.5", "500m" (millicores), "2".
func ParseCPU(s string) (float64, error) {
	q, err := parseQuantity(s, "cpu")
	if err != nil {
		return 0, err
	}

	val := q.AsApproximateFloat64()
	if val <= 0 {
		return 0, fmt.Errorf("cpu must be positive, got %f", val)
	}

	return val, nil
}
