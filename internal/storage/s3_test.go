package storage

import (
	"github.com/stas9132/GophKeeper/internal/logger"
	"log/slog"
	"testing"
)

func TestNewS3(t *testing.T) {
	type args struct {
		l logger.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test", args: struct{ l logger.Logger }{l: slog.Default()}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewS3(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewS3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
