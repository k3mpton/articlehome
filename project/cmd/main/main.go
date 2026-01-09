package main

import (
	"legi/newspapers/project/internal/app"
	"legi/newspapers/project/internal/config"
	logger "legi/newspapers/project/utils/Logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustReadCfg()
	log := logger.InitLogger(cfg.Env)

	app := app.NewApp(log, cfg.Grpc.Port)

	go app.GrpcApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("Stopped server", slog.Int("port", cfg.Grpc.Port))

	app.GrpcApp.StopApp()

	log.Info("Server stop!")
}
