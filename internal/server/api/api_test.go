package api

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/internal/server/db"
	"github.com/stas9132/GophKeeper/keeper"
	"google.golang.org/grpc/metadata"
	"io"
	"reflect"
	"testing"
)

type dbStub struct {
	retMeta *db.Meta
	retErr  error
	retIsOk bool
}

func (d dbStub) Register(ctx context.Context, user, password string) (bool, error) {
	return d.retIsOk, d.retErr
}

func (d dbStub) Login(ctx context.Context, user, password string) (bool, error) {
	return d.retIsOk, d.retErr
}

func (d dbStub) Logout(ctx context.Context) (bool, error) {
	return d.retIsOk, d.retErr
}

func (d dbStub) PutMeta(ctx context.Context, s string, s2 string, i int) (*db.Meta, error) {
	return d.retMeta, d.retErr
}

func (d dbStub) GetMeta(context.Context, string, string) (*db.Meta, error) {
	return d.retMeta, d.retErr
}

type s3Stub struct {
	retUpdInfo minio.UploadInfo
	retErr     error
}

func (s s3Stub) PutObject(ctx context.Context, s2 string, s3 string, reader io.Reader, i int64, options minio.PutObjectOptions) (minio.UploadInfo, error) {
	return s.retUpdInfo, s.retErr
}

func (s s3Stub) MakeBucket(ctx context.Context, s2 string, options minio.MakeBucketOptions) error {
	return s.retErr
}

func TestAPI_Get(t *testing.T) {
	type fields struct {
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		s3                        S3
		db                        DB
	}
	type args struct {
		ctx context.Context
		in  *keeper.ObjMain
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keeper.ObjMain
		wantErr bool
	}{{
		name:    "Unkonown issuer",
		fields:  fields{},
		args:    args{ctx: context.Background()},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Error answer from db.GetMeta()",
		fields:  fields{db: &dbStub{retErr: errors.New("error")}},
		args:    args{ctx: context.WithValue(context.Background(), "iss", "user")},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Ok",
		fields:  fields{db: &dbStub{retMeta: &db.Meta{ObjId: "uuid", ObjType: 1}}},
		args:    args{ctx: context.WithValue(context.Background(), "iss", "user"), in: &keeper.ObjMain{Name: "name"}},
		want:    &keeper.ObjMain{Name: "name", S3Link: "uuid", Type: 1},
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.s3,
				db:                        tt.fields.db,
			}
			got, err := a.Get(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_Health(t *testing.T) {
	type fields struct {
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		s3                        S3
		db                        DB
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
	}{{
		name:    "Ok",
		fields:  fields{},
		args:    args{},
		want:    &keeper.HealthMain{Status: "ok", Version: config.Version, Message: "Service is healthy"},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.s3,
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
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		s3                        S3
		db                        DB
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
	}{{
		name:    "Error answer from db.Login()",
		fields:  fields{db: &dbStub{retErr: errors.New("error")}},
		args:    args{},
		want:    nil,
		wantErr: true,
	}, {
		name:    "NotOk answer from db.Login()",
		fields:  fields{db: &dbStub{retIsOk: false}},
		args:    args{},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Error answer from grpc.SetHeader()",
		fields:  fields{db: &dbStub{retIsOk: true}},
		args:    args{ctx: metadata.NewOutgoingContext(context.Background(), nil)},
		want:    nil,
		wantErr: true,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.s3,
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
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		s3                        S3
		db                        DB
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
	}{{
		name:    "Error answer from db.Logout()",
		fields:  fields{db: &dbStub{retErr: errors.New("error")}},
		args:    args{},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Error answer from grpc.SetHeader()",
		fields:  fields{db: &dbStub{}},
		args:    args{ctx: context.Background()},
		want:    nil,
		wantErr: true,
	},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.s3,
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

func TestAPI_Put(t *testing.T) {
	type fields struct {
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		s3                        S3
		db                        DB
	}
	type args struct {
		ctx context.Context
		in  *keeper.ObjMain
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keeper.Empty
		wantErr bool
	}{{
		name:    "Unkonown issuer",
		fields:  fields{},
		args:    args{ctx: context.Background()},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Error answer from db.PutMeta()",
		fields:  fields{db: &dbStub{retErr: errors.New("error")}},
		args:    args{ctx: context.WithValue(context.Background(), "iss", "user")},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Error answer from s3.PutObject()",
		fields:  fields{db: &dbStub{retMeta: &db.Meta{}}, s3: &s3Stub{retErr: errors.New("error")}},
		args:    args{ctx: context.WithValue(context.Background(), "iss", "user")},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Error answer from s3.PutObject()",
		fields:  fields{db: &dbStub{retMeta: &db.Meta{}}, s3: &s3Stub{}, Logger: logger.NewSlogLogger()},
		args:    args{ctx: context.WithValue(context.Background(), "iss", "user")},
		want:    &keeper.Empty{},
		wantErr: false,
	},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.s3,
				db:                        tt.fields.db,
			}
			got, err := a.Put(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Put() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_Register(t *testing.T) {
	type fields struct {
		Logger                    logger.Logger
		UnimplementedKeeperServer keeper.UnimplementedKeeperServer
		s3                        S3
		db                        DB
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
	}{{
		name:    "Error answer from s3.MakeBucket()",
		fields:  fields{s3: &s3Stub{retErr: errors.New("error")}, Logger: logger.NewSlogLogger()},
		args:    args{},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Error answer from sd.Register()",
		fields:  fields{s3: &s3Stub{}, Logger: logger.NewSlogLogger(), db: &dbStub{retErr: errors.New("error")}},
		args:    args{},
		want:    nil,
		wantErr: true,
	}, {
		name:    "Error answer from grpc.SetHeader()",
		fields:  fields{db: &dbStub{retIsOk: true}, s3: &s3Stub{}},
		args:    args{ctx: metadata.NewOutgoingContext(context.Background(), nil)},
		want:    nil,
		wantErr: true,
	},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &API{
				Logger:                    tt.fields.Logger,
				UnimplementedKeeperServer: tt.fields.UnimplementedKeeperServer,
				s3:                        tt.fields.s3,
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
