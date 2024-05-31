package main

import (
	"log/slog"
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "serverStart"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_slog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "serverStart"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
			slog.Info("hello slog!", slog.String("user", "goher"))
		})
	}
}
