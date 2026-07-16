package parser

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Command
		wantErr bool
	}{
		{
			name:  "SET without TTL",
			input: "SET key1 value1",
			want: &Command{
				Name:  SET,
				Key:   "key1",
				Value: "value1",
				TTL:   0,
			},
			wantErr: false,
		},
		{
			name:  "SET with TTL",
			input: "SET key1 value1 EX 10",
			want: &Command{
				Name:  SET,
				Key:   "key1",
				Value: "value1",
				TTL:   10 * time.Second,
			},
			wantErr: false,
		},
		{
			name:  "GET existing key",
			input: "GET mykey",
			want: &Command{
				Name:  GET,
				Key:   "mykey",
				Value: "",
				TTL:   0,
			},
			wantErr: false,
		},
		{
			name:  "DEL existing key",
			input: "DEL mykey",
			want: &Command{
				Name:  DEL,
				Key:   "mykey",
				Value: "",
				TTL:   0,
			},
			wantErr: false,
		},
		{
			name:  "KEYS command",
			input: "KEYS",
			want: &Command{
				Name:  KEYS,
				Key:   "",
				Value: "",
				TTL:   0,
			},
			wantErr: false,
		},
		{
			name:  "PING command",
			input: "PING",
			want: &Command{
				Name:  PING,
				Key:   "",
				Value: "",
				TTL:   0,
			},
			wantErr: false,
		},
		{
			name:    "empty line",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unknown command",
			input:   "UNKNOWN arg1 arg2",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "SET without key and value",
			input:   "SET",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "SET without value",
			input:   "SET key1",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "GET without key",
			input:   "GET",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "DEL without key",
			input:   "DEL",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "SET with invalid TTL",
			input:   "SET key1 value1 EX abc",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("Parse() = nil, want non-nil Command")
					return
				}

				if got.Name != tt.want.Name {
					t.Errorf("Parse().Name = %q, want %q", got.Name, tt.want.Name)
				}

				if got.Key != tt.want.Key {
					t.Errorf("Parse().Key = %q, want %q", got.Key, tt.want.Key)
				}

				if got.Value != tt.want.Value {
					t.Errorf("Parse().Value = %q, want %q", got.Value, tt.want.Value)
				}

				if got.TTL != tt.want.TTL {
					t.Errorf("Parse().TTL = %v, want %v", got.TTL, tt.want.TTL)
				}
			}
		})
	}
}
