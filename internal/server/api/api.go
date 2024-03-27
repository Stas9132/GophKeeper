package api

import (
	"bytes"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/internal/storage"
	"github.com/stas9132/GophKeeper/keeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"sync"
	"time"
)

type Key struct {
	Name string
	Type int
	u    string
}

type API struct {
	sync.Mutex
	logger.Logger
	keeper.UnimplementedKeeperServer
	storage Storage
	db      map[string][]Key
}

type Storage interface {
	Register(ctx context.Context, user, password string) (bool, error)
	Login(ctx context.Context, user, password string) (bool, error)
	Logout(ctx context.Context) (bool, error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
}

func NewAPI(logger logger.Logger) (*API, error) {
	s3, err := storage.NewS3(logger)
	if err != nil {
		return nil, err
	}
	return &API{
		Logger:  logger,
		storage: s3,
		db:      make(map[string][]Key),
	}, nil
}

func (a *API) Health(ctx context.Context, in *keeper.Empty) (*keeper.HealthMain, error) {
	return &keeper.HealthMain{
		Status:  "ok",
		Version: config.Version,
		Message: "Service is healthy",
	}, nil
}

const TTL = time.Hour

func (a *API) Register(ctx context.Context, in *keeper.AuthMain) (*keeper.Empty, error) {
	if ok, err := a.storage.Register(ctx, in.GetUser(), in.GetPassword()); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	} else if !ok {
		return nil, status.Error(codes.AlreadyExists, "already exist")
	}

	j, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"iss": in.GetUser(), "exp": time.Now().Add(TTL).Unix()},
	).SignedString([]byte("123"))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	if err = grpc.SetHeader(ctx, metadata.Pairs("authorization", j)); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &keeper.Empty{}, nil
}

func (a *API) Login(ctx context.Context, in *keeper.AuthMain) (*keeper.Empty, error) {
	if ok, err := a.storage.Login(ctx, in.GetUser(), in.GetPassword()); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	} else if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	j, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"iss": in.GetUser(), "exp": time.Now().Add(TTL).Unix()},
	).SignedString([]byte("123"))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	if err = grpc.SetHeader(ctx, metadata.Pairs("authorization", j)); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &keeper.Empty{}, nil
}

func (a *API) Logout(ctx context.Context, in *keeper.Empty) (*keeper.Empty, error) {
	if err := grpc.SetHeader(ctx, metadata.Pairs("authorization", "")); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &keeper.Empty{}, nil
}

func (a *API) Put(ctx context.Context, in *keeper.ObjMain) (*keeper.Empty, error) {
	iss, ok := ctx.Value("iss").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	a.Lock()
	ids := a.db[iss]
	var u string
	for _, id := range ids {
		if id.Name == in.GetName() {
			u = id.u
			goto lunlock
		}
	}
	u = uuid.NewString()
	a.db[iss] = append(ids, Key{Name: in.GetName(), Type: int(in.GetType()), u: u})
lunlock:
	a.Unlock()

	info, err := a.storage.PutObject(ctx, iss, u, bytes.NewReader(in.GetEncData()), int64(len(in.GetEncData())), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	log.Println(u, info)
	return &keeper.Empty{}, nil
}
func (a *API) Get(ctx context.Context, in *keeper.ObjMain) (*keeper.ObjMain, error) {
	iss, ok := ctx.Value("iss").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	a.Lock()
	var retErr error
	var id Key
	ids := a.db[iss]
	for _, id = range ids {
		if id.Name == in.GetName() {
			goto lunlock
		}
	}
	retErr = status.Error(codes.Unknown, "not found")
lunlock:
	a.Unlock()

	in.S3Link = id.u
	return in, retErr
}

func (a *API) Sync(ctx context.Context, in *keeper.SyncMain) (*keeper.SyncMain, error) {
	a.Lock()
	defer a.Unlock()
	u, ok := ctx.Value("iss").(string)
	if !ok || len(u) == 0 {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	inKeys := make([]Key, len(in.GetKeys()))
	for i, key := range in.GetKeys() {
		inKeys[i] = Key{
			Name: key.GetName(),
			Type: int(key.GetType()),
		}
	}
	var outKeys []Key
	oldKeys := a.db[u]
	hash := make(map[Key]struct{})
	for _, s := range append(oldKeys, inKeys...) {
		hash[s] = struct{}{}
	}
	outKeys = make([]Key, 0, len(hash))
	for key := range hash {
		outKeys = append(outKeys, key)
	}
	a.db[u] = outKeys

	in.Keys = make([]*keeper.SyncMain_KeysMain, len(outKeys))
	for i, key := range outKeys {
		in.Keys[i] = &keeper.SyncMain_KeysMain{
			Name: key.Name,
			Type: keeper.TypeCode(key.Type),
		}
	}
	return in, nil
}
