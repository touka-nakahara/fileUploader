package server

import "testing"

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Just Run"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewServer()
		})
	}
}
