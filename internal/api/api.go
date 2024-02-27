package api

import (
	"context"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
	"github.com/stas9132/GophKeeper/keeper"
)

type API struct {
	logger.Logger
	keeper.UnimplementedKeeperServer
}

func NewAPI(logger logger.Logger) *API {
	return &API{Logger: logger}
}

func (a *API) Health(ctx context.Context, in *keeper.Empty) (*keeper.HealthMain, error) {
	return &keeper.HealthMain{
		Status:  "ok",
		Version: config.Version,
		Message: "Service is healthy",
	}, nil
}
