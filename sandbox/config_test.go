package sandbox

import "testing"

func TestParseVolume(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Volume
		wantErr bool
	}{
		{
			name:  "simple bind",
			input: "/app:/home/user/mydir",
			want:  Volume{Source: "/home/user/mydir", Target: "/app", ReadOnly: false},
		},
		{
			name:  "read-only bind",
			input: "/app:/home/user/mydir:ro",
			want:  Volume{Source: "/home/user/mydir", Target: "/app", ReadOnly: true},
		},
		{
			name:    "invalid format",
			input:   "/only-one-part",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVolume(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVolume(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseVolume(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseMemory(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		{name: "mebibytes", input: "256Mi", want: 256 * 1024 * 1024},
		{name: "gibibytes", input: "2Gi", want: 2 * 1024 * 1024 * 1024},
		{name: "kibibytes", input: "512Ki", want: 512 * 1024},
		{name: "decimal mega", input: "256M", want: 256 * 1000 * 1000},
		{name: "decimal giga", input: "2G", want: 2 * 1000 * 1000 * 1000},
		{name: "plain bytes", input: "1048576", want: 1048576},
		{name: "invalid", input: "abc", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMemory(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMemory(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseMemory(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseCPU(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{name: "one core", input: "1", want: 1.0},
		{name: "half core", input: "0.5", want: 0.5},
		{name: "two cores", input: "2", want: 2.0},
		{name: "quarter", input: "0.25", want: 0.25},
		{name: "millicores 500m", input: "500m", want: 0.5},
		{name: "millicores 100m", input: "100m", want: 0.1},
		{name: "millicores 2000m", input: "2000m", want: 2.0},
		{name: "zero", input: "0", wantErr: true},
		{name: "negative", input: "-1", wantErr: true},
		{name: "invalid", input: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCPU(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCPU(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseCPU(%q) = %f, want %f", tt.input, got, tt.want)
			}
		})
	}
}
