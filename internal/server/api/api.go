package api

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/internal/storage"
	"github.com/stas9132/GophKeeper/keeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type API struct {
	logger.Logger
	keeper.UnimplementedKeeperServer
	storage Storage
	db      *sync.Map
}

type Storage interface {
	Register(ctx context.Context, user, password string) (bool, error)
	Login(ctx context.Context, user, password string) (bool, error)
	Logout(ctx context.Context) (bool, error)
}

func NewAPI(logger logger.Logger) (*API, error) {
	s3, err := storage.NewS3(logger)
	if err != nil {
		return nil, err
	}
	return &API{
		Logger:  logger,
		storage: s3,
		db:      &sync.Map{},
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
	).SigningString()
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
	).SigningString()
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

func (a *API) Sync(ctx context.Context, in *keeper.SyncMain) (*keeper.SyncMain, error) {

	return in, nil
}
