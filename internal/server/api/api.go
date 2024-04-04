package api

import (
	"bytes"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/minio/minio-go/v7"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/internal/server/db"
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
	s3 S3
	db DB
}

type S3 interface {
	PutObject(context.Context, string, string, io.Reader, int64, minio.PutObjectOptions) (minio.UploadInfo, error)
	MakeBucket(context.Context, string, minio.MakeBucketOptions) error
}

type DB interface {
	Register(ctx context.Context, user, password string) (bool, error)
	Login(ctx context.Context, user, password string) (bool, error)
	Logout(ctx context.Context) (bool, error)
}

func NewAPI(logger logger.Logger) (*API, error) {
	s3, err := storage.NewS3(logger)
	if err != nil {
		return nil, err
	}
	db, err := db.NewDB(logger)
	if err != nil {
		return nil, err
	}
	return &API{
		Logger: logger,
		s3:     s3,
		db:     db,
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
	if err := a.s3.MakeBucket(ctx, in.GetUser(), minio.MakeBucketOptions{}); err != nil {
		a.Error("Unable to make bucket: " + in.GetUser() + ". Error: " + err.Error())
		return nil, status.Error(codes.Unknown, err.Error())
	}
	if _, err := a.db.Register(ctx, in.GetUser(), in.GetPassword()); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	j, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"iss": in.GetUser(), "exp": time.Now().Add(TTL).Unix()},
	).SignedString(config.JwtKey)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	if err = grpc.SetHeader(ctx, metadata.Pairs("authorization", j)); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &keeper.Empty{}, nil
}

func (a *API) Login(ctx context.Context, in *keeper.AuthMain) (*keeper.Empty, error) {
	if ok, err := a.db.Login(ctx, in.GetUser(), in.GetPassword()); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	} else if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	j, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{"iss": in.GetUser(), "exp": time.Now().Add(TTL).Unix()},
	).SignedString(config.JwtKey)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	if err = grpc.SetHeader(ctx, metadata.Pairs("authorization", j)); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &keeper.Empty{}, nil
}

func (a *API) Logout(ctx context.Context, in *keeper.Empty) (*keeper.Empty, error) {
	if _, err := a.db.Logout(ctx); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
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
	//ids := a.db[iss]
	var u string
	//var f bool
	//for _, id := range ids {
	//	if id.Name == in.GetName() {
	//		u = id.u
	//		f = true
	//		break
	//	}
	//}
	//if !f {
	//	u = uuid.NewString()
	//	a.db[iss] = append(ids, Key{Name: in.GetName(), Type: int(in.GetType()), u: u})
	//}
	a.Unlock()

	info, err := a.s3.PutObject(ctx, iss, u, bytes.NewReader(in.GetEncData()), int64(len(in.GetEncData())), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	a.Info(info.Bucket + ":" + info.Key + " stored in s3")
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
	log.Println(iss)
	//var f bool
	//ids := a.db[iss]
	//for _, id = range ids {
	//	if id.Name == in.GetName() {
	//		f = true
	//		break
	//	}
	//}
	//if !f {
	//	retErr = status.Error(codes.Unknown, "not found")
	//}
	a.Unlock()

	errMsg := "noerror"
	if retErr != nil {
		errMsg = retErr.Error()
	}
	a.Info("find " + in.GetName() + ": " + errMsg)

	in.S3Link = id.u
	in.Type = keeper.TypeCode(id.Type)
	return in, retErr
}
