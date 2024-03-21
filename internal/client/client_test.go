package client

import (
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/keeper"
	"google.golang.org/grpc/credentials"
	"reflect"
	"sync"
	"testing"
)

func TestClient_Get(t *testing.T) {
	type fields struct {
		KeeperClient keeper.KeeperClient
		Logger       logger.Logger
		user         string
		token        string
		s3           Storage
		storedKeys   []Keys
		Mutex        sync.Mutex
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
	}{
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
				storedKeys:   tt.fields.storedKeys,
				Mutex:        tt.fields.Mutex,
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
		s3           Storage
		storedKeys   []Keys
		Mutex        sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
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
				storedKeys:   tt.fields.storedKeys,
				Mutex:        tt.fields.Mutex,
			}
			if err := c.Health(); (err != nil) != tt.wantErr {
				t.Errorf("Health() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_List(t *testing.T) {
	type fields struct {
		KeeperClient keeper.KeeperClient
		Logger       logger.Logger
		user         string
		token        string
		s3           Storage
		storedKeys   []Keys
		Mutex        sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Keys
		wantErr bool
	}{
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
				storedKeys:   tt.fields.storedKeys,
				Mutex:        tt.fields.Mutex,
			}
			got, err := c.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
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
		s3           Storage
		storedKeys   []Keys
		Mutex        sync.Mutex
	}
	type args struct {
		flds []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
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
				storedKeys:   tt.fields.storedKeys,
				Mutex:        tt.fields.Mutex,
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
		s3           Storage
		storedKeys   []Keys
		Mutex        sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
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
				storedKeys:   tt.fields.storedKeys,
				Mutex:        tt.fields.Mutex,
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
		s3           Storage
		storedKeys   []Keys
		Mutex        sync.Mutex
	}
	type args struct {
		flds []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
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
				storedKeys:   tt.fields.storedKeys,
				Mutex:        tt.fields.Mutex,
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
		s3           Storage
		storedKeys   []Keys
		Mutex        sync.Mutex
	}
	type args struct {
		flds []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
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
				storedKeys:   tt.fields.storedKeys,
				Mutex:        tt.fields.Mutex,
			}
			if err := c.Register(tt.args.flds); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SyncList(t *testing.T) {
	type fields struct {
		KeeperClient keeper.KeeperClient
		Logger       logger.Logger
		user         string
		token        string
		s3           Storage
		storedKeys   []Keys
		Mutex        sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
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
				storedKeys:   tt.fields.storedKeys,
				Mutex:        tt.fields.Mutex,
			}
			if err := c.SyncList(); (err != nil) != tt.wantErr {
				t.Errorf("SyncList() error = %v, wantErr %v", err, tt.wantErr)
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
