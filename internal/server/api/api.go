package api

import (
	"bytes"
	"context"
	"fmt"
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
	"time"
)

type Key struct {
	Name string
	Type int
	u    string
}

type API struct {
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
	PutMeta(context.Context, string, string, int) (*db.Meta, error)
	GetMeta(context.Context, string, string) (*db.Meta, error)
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

func (a *API) Put(server keeper.Keeper_PutServer) error {
	fmt.Println("a", server.Context())

	iss, ok := server.Context().Value("iss").(string)
	if !ok {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}

	fmt.Println("aa")

	pr, pw := io.Pipe()
	obj, err := server.Recv()
	if err != nil {
		a.Error("Recv 1-st chunk", "error", err)
		return status.Error(codes.Unknown, err.Error())
	}
	name := obj.GetName()
	Type := obj.GetType()
	size := obj.GetSize()
	if _, err = io.Copy(pw, bytes.NewReader(obj.GetEncData())); err != nil {
		a.Error("io.Copy()", "error", err)
		return status.Error(codes.Unknown, err.Error())
	}
	go func() {
		for obj2, err2 := server.Recv(); err2 == nil; obj2, err2 = server.Recv() {
			if _, err2 = io.Copy(pw, bytes.NewReader(obj2.GetEncData())); err != nil {
				a.Error("io.Copy()", "error", err2)
				return
			}
		}
	}()

	meta, err := a.db.PutMeta(server.Context(), iss, name, int(Type))
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}

	info, err := a.s3.PutObject(server.Context(), iss, meta.ObjId, pr, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}
	a.Info(info.Bucket + ":" + info.Key + " stored in s3")
	return nil
}
func (a *API) Get(ctx context.Context, in *keeper.ObjMain) (*keeper.ObjMain, error) {
	iss, ok := ctx.Value("iss").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	meta, err := a.db.GetMeta(ctx, iss, in.GetName())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	in.S3Link = meta.ObjId
	in.Type = keeper.TypeCode(meta.ObjType)
	return in, nil
}
