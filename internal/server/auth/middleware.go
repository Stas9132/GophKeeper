package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stas9132/GophKeeper/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func interceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		switch info.FullMethod {
		case //"/keeper.keeper/Health",
			"/keeper.keeper/Register",
			"/keeper.keeper/Login",
			"/keeper.keeper/Logout":
			return handler(ctx, req)
		default:
		}

		// Get metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		tmp, ok := md["authorization"]
		if !ok || len(tmp) < 1 {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		j := tmp[0]

		t, err := jwt.ParseWithClaims(j, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signed method")
			}
			return nil, nil
		})
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		if c, ok := t.Claims.(*jwt.MapClaims); ok && t.Valid {
			if u, ok := (*c)["iss"].(string); ok {
				ctx = context.WithValue(ctx, "", u)
			}
		} else {
			fmt.Println(c)
		}

		return handler(ctx, req)
	}
}

func UnaryServerInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return interceptor(logger)
}
