package client

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/keeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"reflect"
	"testing"
	"testing/fstest"
)

type kClStub struct {
	retObj    *keeper.ObjMain
	retHealth *keeper.HealthMain
	retErr    error
}

func (s *kClStub) Health(context.Context, *keeper.Empty, ...grpc.CallOption) (*keeper.HealthMain, error) {
	return s.retHealth, s.retErr
}

func (s *kClStub) Register(context.Context, *keeper.AuthMain, ...grpc.CallOption) (*keeper.Empty, error) {
	return &keeper.Empty{}, s.retErr
}

func (s *kClStub) Login(ctx context.Context, in *keeper.AuthMain, opts ...grpc.CallOption) (*keeper.Empty, error) {
	return &keeper.Empty{}, s.retErr
}

func (s *kClStub) Logout(context.Context, *keeper.Empty, ...grpc.CallOption) (*keeper.Empty, error) {
	return &keeper.Empty{}, s.retErr
}

func (s *kClStub) Put(ctx context.Context, opts ...grpc.CallOption) (keeper.Keeper_PutClient, error) {
	return &kPutClStub{}, s.retErr
}

func (s *kClStub) Get(context.Context, *keeper.ObjMain, ...grpc.CallOption) (*keeper.ObjMain, error) {
	return s.retObj, s.retErr
}

type kPutClStub struct {
}

func (k kPutClStub) Send(main *keeper.ObjMain) error {
	return nil
}

func (k kPutClStub) CloseAndRecv() (*keeper.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (k kPutClStub) Header() (metadata.MD, error) {
	//TODO implement me
	panic("implement me")
}

func (k kPutClStub) Trailer() metadata.MD {
	//TODO implement me
	panic("implement me")
}

func (k kPutClStub) CloseSend() error {
	return nil
}

func (k kPutClStub) Context() context.Context {
	//TODO implement me
	panic("implement me")
}

func (k kPutClStub) SendMsg(m any) error {
	//TODO implement me
	panic("implement me")
}

func (k kPutClStub) RecvMsg(m any) error {
	//TODO implement me
	panic("implement me")
}

type s3Stub struct {
	retObj *minio.Object
	retErr error
}

func (s *s3Stub) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	return s.retObj, s.retErr
}

type minioObjStub struct {
	retN   int
	retErr error
}

func (s *minioObjStub) Read([]byte) (int, error) {
	return s.retN, s.retErr
}

type loggerStub struct {
}

func (l loggerStub) Info(msg string, args ...any) {
}

func (l loggerStub) Warn(msg string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l loggerStub) Error(msg string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func TestClient_Get(t *testing.T) {
	type fields struct {
		KeeperClient keeper.KeeperClient
		Logger       logger.Logger
		user         string
		token        string
		s3           S3
	}
	type args struct {
		flds []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		f       func()
	}{{
		name:    "Wrong arg flds",
		fields:  fields{},
		args:    args{flds: []string{"1"}},
		want:    "",
		wantErr: true,
	}, {
		name:    "Error answer from KeeperClient.Get()",
		fields:  fields{KeeperClient: &kClStub{retErr: errors.New("error")}},
		args:    args{flds: []string{"1", "2"}},
		want:    "",
		wantErr: true,
	}, {
		name:    "Error answer from S3.Get()",
		fields:  fields{KeeperClient: &kClStub{retObj: &keeper.ObjMain{}}, s3: &s3Stub{retErr: errors.New("error")}},
		args:    args{flds: []string{"1", "2"}},
		want:    "",
		wantErr: true,
	},
	//{
	//	name: "Error answer from aes.NewCipher()",
	//	fields: fields{KeeperClient: &kClStub{retObj: &keeper.ObjMain{}},
	//		s3: &s3Stub{retObj: &minio.Object{}}},
	//	args:    args{flds: []string{"1", "2"}},
	//	want:    "",
	//	wantErr: true,
	//	f: func() {
	//		config.AESKey = []byte("wrong")
	//	},
	//},
	}

	// TODO: Add test cases.

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.f != nil {
				tt.f()
			}
			c := &Client{
				KeeperClient: tt.fields.KeeperClient,
				Logger:       tt.fields.Logger,
				user:         tt.fields.user,
				token:        tt.fields.token,
				s3:           tt.fields.s3,
			}
			got, err := c.Get(tt.args.flds)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Health(t *testing.T) {
	type fields struct {
		KeeperClient keeper.KeeperClient
		Logger       logger.Logger
		user         string
		token        string
		s3           S3
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{{
		name: "Need mock for verify len(c.tocken) > 0",
		fields: fields{token: "tocken",
			KeeperClient: &kClStub{retHealth: &keeper.HealthMain{}},
			Logger:       &loggerStub{}},
		wantErr: false,
	}, {
		name: "Error answer from KeeperClient.Health()",
		fields: fields{token: "tocken",
			KeeperClient: &kClStub{retErr: errors.New("error")}},
		wantErr: true,
	}, {
		name: "Ok",
		fields: fields{token: "tocken",
			KeeperClient: &kClStub{retHealth: &keeper.HealthMain{}},
			Logger:       &loggerStub{}},
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				KeeperClient: tt.fields.KeeperClient,
				Logger:       tt.fields.Logger,
				user:         tt.fields.user,
				token:        tt.fields.token,
				s3:           tt.fields.s3,
			}
			if err := c.Health(); (err != nil) != tt.wantErr {
				t.Errorf("Health() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Login(t *testing.T) {
	type fields struct {
		KeeperClient keeper.KeeperClient
		Logger       logger.Logger
		user         string
		token        string
		s3           S3
	}
	type args struct {
		flds []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{{
		name:    "Wrong arg flds",
		fields:  fields{},
		args:    args{flds: []string{"1"}},
		wantErr: true,
	}, {
		name:    "Error answer from KeeperClient.Login",
		fields:  fields{KeeperClient: &kClStub{retErr: errors.New("error")}},
		args:    args{flds: []string{"1", "2", "3"}},
		wantErr: true,
	},
	//{
	//	name:    "Ok, reflect canset() return false",
	//	fields:  fields{KeeperClient: &kClStub{}},
	//	args:    args{flds: []string{"1", "2", "3"}},
	//	wantErr: false,
	//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				KeeperClient: tt.fields.KeeperClient,
				Logger:       tt.fields.Logger,
				user:         tt.fields.user,
				token:        tt.fields.token,
				s3:           tt.fields.s3,
			}
			if err := c.Login(tt.args.flds); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Logout(t *testing.T) {
	type fields struct {
		KeeperClient keeper.KeeperClient
		Logger       logger.Logger
		user         string
		token        string
		s3           S3
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{{
		name:    "Error answer from KeeperClient.Logout()",
		fields:  fields{KeeperClient: &kClStub{retErr: errors.New("error")}},
		wantErr: true,
	}, {
		name:    "Ok",
		fields:  fields{KeeperClient: &kClStub{}},
		wantErr: false,
	},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				KeeperClient: tt.fields.KeeperClient,
				Logger:       tt.fields.Logger,
				user:         tt.fields.user,
				token:        tt.fields.token,
				s3:           tt.fields.s3,
			}
			if err := c.Logout(); (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Put(t *testing.T) {
	type fields struct {
		KeeperClient keeper.KeeperClient
		Logger       logger.Logger
		user         string
		token        string
		s3           S3
	}
	type args struct {
		flds []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		f       func()
	}{
		{
			name:    "Wrong arg flds",
			fields:  fields{},
			args:    args{flds: []string{"1"}},
			wantErr: true,
		},
		{
			name:    "Key login/password & error on KeeperClient.Put()",
			fields:  fields{KeeperClient: &kClStub{retErr: errors.New("error")}},
			args:    args{flds: []string{"1", "2", "1", "4"}},
			wantErr: true,
		},
		{
			name:    "Key text & error on KeeperClient.Put()",
			fields:  fields{KeeperClient: &kClStub{retErr: errors.New("error")}},
			args:    args{flds: []string{"1", "2", "2", "4"}},
			wantErr: true,
		},
		{
			name:    "Key file & error on Open()",
			fields:  fields{KeeperClient: &kClStub{retErr: errors.New("error")}},
			args:    args{flds: []string{"1", "2", "3", "4"}},
			wantErr: true,
		},
		{
			name:    "Key file & error on Stat()",
			fields:  fields{KeeperClient: &kClStub{retErr: errors.New("error")}},
			args:    args{flds: []string{"1", "2", "3", "4"}},
			wantErr: true,
			f: func() {
				dirFS = fstest.MapFS{"4": &fstest.MapFile{}}
			},
		},
		{
			name:    "Key card & error on KeeperClient.Put()",
			fields:  fields{KeeperClient: &kClStub{retErr: errors.New("error")}},
			args:    args{flds: []string{"1", "2", "4", "4"}},
			wantErr: true,
		},
		{
			name:    "Unknown key name",
			fields:  fields{},
			args:    args{flds: []string{"1", "2", "", "4"}},
			wantErr: true,
		},
		{
			name:    "Ok",
			fields:  fields{KeeperClient: &kClStub{}, Logger: &loggerStub{}},
			args:    args{flds: []string{"1", "2", "2", "4"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.f != nil {
				tt.f()
			}
			c := &Client{
				KeeperClient: tt.fields.KeeperClient,
				Logger:       tt.fields.Logger,
				user:         tt.fields.user,
				token:        tt.fields.token,
				s3:           tt.fields.s3,
			}
			if err := c.Put(tt.args.flds); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Register(t *testing.T) {
	type fields struct {
		KeeperClient keeper.KeeperClient
		Logger       logger.Logger
		user         string
		token        string
		s3           S3
	}
	type args struct {
		flds []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{{
		name:    "Wrong arg flds",
		fields:  fields{},
		args:    args{flds: []string{"1"}},
		wantErr: true,
	}, {
		name:    "Wrong arg flds",
		fields:  fields{KeeperClient: &kClStub{retErr: errors.New("error")}},
		args:    args{flds: []string{"1", "2", "3"}},
		wantErr: true,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				KeeperClient: tt.fields.KeeperClient,
				Logger:       tt.fields.Logger,
				user:         tt.fields.user,
				token:        tt.fields.token,
				s3:           tt.fields.s3,
			}
			if err := c.Register(tt.args.flds); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		l       logger.Logger
		tlsCred credentials.TransportCredentials
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.l, tt.args.tlsCred)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient1(t *testing.T) {
	type args struct {
		l       logger.Logger
		tlsCred credentials.TransportCredentials
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{{name: "ok", args: struct {
		l       logger.Logger
		tlsCred credentials.TransportCredentials
	}{l: slog.Default(), tlsCred: credentials.NewTLS(nil)}, want: &Client{}, wantErr: false},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.l, tt.args.tlsCred)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}
