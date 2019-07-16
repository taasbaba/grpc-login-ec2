package grpc

import (
	"context"
	"google.golang.org/grpc"
	auth "grpc-login-server/server/internal/api/v1"
	"grpc-login-server/server/internal/logger"
	"grpc-login-server/server/internal/protocol/grpc/middleware"
	"net"
	"os"
	"os/signal"
)

func RunServer(ctx context.Context, authServer auth.AuthServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// middleware add here
	// gRPC server statup options
	opts := []grpc.ServerOption{}

	// add middleware
	opts = middleware.AddLogging(logger.Log, opts)
	// END middleware add here

	server := grpc.NewServer(opts...)
	auth.RegisterAuthServer(server, authServer)
	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down gRPC AuthServer...")

			server.Stop()

			<-ctx.Done()
		}
	}()
	// start gRPC server
	logger.Log.Info("starting gRPC AuthServer...")
	return server.Serve(listen)
}
