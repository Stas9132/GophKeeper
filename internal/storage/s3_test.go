package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/stas9132/GophKeeper/internal/logger"
	"reflect"
	"testing"
)

func TestNewS3(t *testing.T) {
	type args struct {
		l logger.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    *S3
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewS3(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewS3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewS3() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3_Login(t *testing.T) {
	type fields struct {
		Client *minio.Client
		Logger logger.Logger
	}
	type args struct {
		ctx      context.Context
		user     string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &S3{
				Client: tt.fields.Client,
				Logger: tt.fields.Logger,
			}
			got, err := s.Login(tt.args.ctx, tt.args.user, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3_Logout(t *testing.T) {
	type fields struct {
		Client *minio.Client
		Logger logger.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &S3{
				Client: tt.fields.Client,
				Logger: tt.fields.Logger,
			}
			got, err := s.Logout(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Logout() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS3_Register(t *testing.T) {
	type fields struct {
		Client *minio.Client
		Logger logger.Logger
	}
	type args struct {
		ctx      context.Context
		user     string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &S3{
				Client: tt.fields.Client,
				Logger: tt.fields.Logger,
			}
			got, err := s.Register(tt.args.ctx, tt.args.user, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}
