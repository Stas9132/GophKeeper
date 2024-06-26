package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
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
			return config.JwtKey, nil
		})
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		if c, ok := t.Claims.(*jwt.MapClaims); ok && t.Valid {
			if u, ok := (*c)["iss"].(string); ok {
				ctx = context.WithValue(ctx, "iss", u)
			}
		} else {
			fmt.Println(c)
		}

		return handler(ctx, req)
	}
}

type serverStreamWrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *serverStreamWrapper) Context() context.Context { return w.ctx }

func interceptorStream(logger logger.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		switch info.FullMethod {
		case //"/keeper.keeper/Health",
			"/keeper.keeper/Register",
			"/keeper.keeper/Login",
			"/keeper.keeper/Logout":
			return handler(srv, ss)
		default:
		}

		// Get metadata from context
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		tmp, ok := md["authorization"]
		if !ok || len(tmp) < 1 {
			return status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		j := tmp[0]

		t, err := jwt.ParseWithClaims(j, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signed method")
			}
			return config.JwtKey, nil
		})
		if err != nil {
			return status.Errorf(codes.Unauthenticated, err.Error())
		}

		ctx := ss.Context()
		if c, ok := t.Claims.(*jwt.MapClaims); ok && t.Valid {
			if u, ok := (*c)["iss"].(string); ok {
				ctx = context.WithValue(ctx, "iss", u)
			}
		} else {
			log.Println(c)
		}

		return handler(srv, &serverStreamWrapper{
			ServerStream: ss,
			ctx:          ctx,
		})
	}
}

func UnaryServerInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return interceptor(logger)
}

func StreamServerInterceptor(logger logger.Logger) grpc.StreamServerInterceptor {
	return interceptorStream(logger)
}
