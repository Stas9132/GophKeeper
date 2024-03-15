package main

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	r2 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/server/api"
	"github.com/stas9132/GophKeeper/internal/server/auth"
	"github.com/stas9132/GophKeeper/keeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var ll *slog.Logger

func run(ctx context.Context) <-chan error {

	interceptorLogger := func(l *slog.Logger) logging.Logger {
		return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(lvl), msg, fields...)
		})
	}

	res := make(chan error, 1)
	time.AfterFunc(100*time.Millisecond, func() {
		s := grpc.NewServer(grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(ll),
				[]logging.Option{logging.WithLogOnEvents(logging.StartCall, logging.FinishCall)}...),
			auth.UnaryServerInterceptor(ll),
		))
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			lis, err := net.Listen("tcp", config.ListenAddress)
			if err != nil {
				log.Fatalln(lis, err)
			}
			a, err := api.NewAPI(ll)
			if err != nil {
				log.Fatalln(err)
			}
			keeper.RegisterKeeperServer(s, a)
			reflection.Register(s)
			log.Println("gRPC control service starts", config.ListenAddress)
			if err = s.Serve(lis); err != nil {
				log.Fatalln(err)
			}
			log.Println("gRPC control service closed")
			wg.Done()
		}()
		<-ctx.Done()
		s.GracefulStop()
		wg.Wait()
		res <- nil
	})
	return res
}

func runREST(ctx context.Context) <-chan error {
	res := make(chan error, 1)
	time.AfterFunc(300*time.Millisecond, func() {
		mux := r2.NewServeMux()
		s := &http.Server{
			Addr:    config.ListenAddressR,
			Handler: mux,
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
			err := keeper.RegisterKeeperHandlerFromEndpoint(ctx, mux, config.ListenAddress, opts)
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("REST server listening at " + config.ListenAddressR)
			if err = s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalln(err)
			}
			log.Println("REST control service closed")
			wg.Done()
		}()
		<-ctx.Done()
		s.Close()
		wg.Wait()
		res <- nil
	})
	return res
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ll = slog.Default()

	eg := run(ctx)
	er := runREST(ctx)

	<-eg
	<-er
}
