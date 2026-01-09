package app

import (
	grpcapp "legi/newspapers/project/internal/app/grpc"
	"legi/newspapers/project/internal/service"
	"legi/newspapers/project/internal/storage"
	"log/slog"
)

type App struct {
	GrpcApp *grpcapp.GrpcApp
}

func NewApp(log *slog.Logger, port int) *App {
	storage := storage.NewStorage()

	service := service.NewService(log, storage, storage)

	app := grpcapp.NewApp(port, log, &service)
	return &App{
		GrpcApp: app,
	}

}
