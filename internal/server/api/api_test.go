package api

import (
	"context"
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/keeper"
	"reflect"
	"sync"
	"testing"
)

func TestAPI_Health(t *testing.T) {
	type fields struct {
		Mutex                     sync.Mutex
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		storage                   S3
		db                        map[string][]Key
	}
	type args struct {
		ctx context.Context
		in  *keeper.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keeper.HealthMain
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Mutex:                     tt.fields.Mutex,
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.storage,
				db:                        tt.fields.db,
			}
			got, err := a.Health(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Health() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Health() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_Login(t *testing.T) {
	type fields struct {
		Mutex                     sync.Mutex
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		storage                   S3
		db                        map[string][]Key
	}
	type args struct {
		ctx context.Context
		in  *keeper.AuthMain
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keeper.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Mutex:                     tt.fields.Mutex,
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.storage,
				db:                        tt.fields.db,
			}
			got, err := a.Login(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_Logout(t *testing.T) {
	type fields struct {
		Mutex                     sync.Mutex
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		storage                   S3
		db                        map[string][]Key
	}
	type args struct {
		ctx context.Context
		in  *keeper.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keeper.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Mutex:                     tt.fields.Mutex,
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.storage,
				db:                        tt.fields.db,
			}
			got, err := a.Logout(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Logout() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_Register(t *testing.T) {
	type fields struct {
		Mutex                     sync.Mutex
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		storage                   S3
		db                        map[string][]Key
	}
	type args struct {
		ctx context.Context
		in  *keeper.AuthMain
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keeper.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Mutex:                     tt.fields.Mutex,
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.storage,
				db:                        tt.fields.db,
			}
			got, err := a.Register(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_Sync(t *testing.T) {
	type fields struct {
		Mutex                     sync.Mutex
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		storage                   S3
		db                        map[string][]Key
	}
	type args struct {
		ctx context.Context
		in  *keeper.SyncMain
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keeper.SyncMain
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Mutex:                     tt.fields.Mutex,
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.storage,
				db:                        tt.fields.db,
			}
			got, err := a.Sync(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sync() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAPI(t *testing.T) {
	type args struct {
		logger logger.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    *API
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAPI(tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPI() got = %v, want %v", got, tt.want)
			}
		})
	}
}
