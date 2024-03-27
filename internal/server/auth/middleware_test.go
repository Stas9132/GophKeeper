package auth

import (
	"github.com/stas9132/GophKeeper/internal/logger"
	"google.golang.org/grpc"
	"reflect"
	"testing"
)

func TestUnaryServerInterceptor(t *testing.T) {
	type args struct {
		logger logger.Logger
	}
	tests := []struct {
		name string
		args args
		want grpc.UnaryServerInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnaryServerInterceptor(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnaryServerInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_interceptor(t *testing.T) {
	type args struct {
		logger logger.Logger
	}
	tests := []struct {
		name string
		args args
		want grpc.UnaryServerInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := interceptor(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("interceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}
