package grpcapp

import (
	"fmt"
	"legi/newspapers/project/internal/server"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcApp struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func NewApp(port int, log *slog.Logger, publ server.PublisherArticles) *GrpcApp {
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	server.NewServer(grpcServer, publ)

	return &GrpcApp{
		log:        log,
		port:       port,
		grpcServer: grpcServer,
	}
}

func (a *GrpcApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *GrpcApp) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		"op", op,
		"port", a.port,
	)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%v", a.port))
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	log.Info("start server", slog.String("addr", conn.Addr().String()))

	if err := a.grpcServer.Serve(conn); err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	return nil
}

func (a *GrpcApp) StopApp() {
	const op = "grpcapp.StopApp"

	a.log.With(
		"op", op,
	).Info("stopped grpc server...", slog.Int("port", a.port))

	a.grpcServer.GracefulStop()
}
