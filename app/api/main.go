package main

import (
	"fmt"
	"log/slog"

	"github.com/dibikhairurrazi/audio-storage/config"
	"github.com/dibikhairurrazi/audio-storage/db"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/converter"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/handler"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/repository"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/service"
	"github.com/dibikhairurrazi/audio-storage/module/audio-storage/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config file", "err", err.Error())
		panic(err)
	}

	database, err := db.Initialize(&cfg.Database)
	if err != nil {
		slog.Error("failed to initialize db connection", "err", err.Error())
		panic(err)
	}

	err = db.MigrateUp(database.MasterConn, "migration", cfg.Database.Master.DBName)
	if err != nil {
		slog.Error("failed to run migration on db", "err", err.Error())
		panic(err)
	}

	if cfg.Server.Mode != "release" {
		slog.Info("running seed on database", "mode", cfg.Server.Mode)
		db.SeedDB(database.MasterConn)
	}

	repo := repository.NewPostgreSQL(*database)
	converter := converter.New()
	storage := storage.NewLocalStorage(cfg.Storage.RootFolder)

	phraseService := service.NewPhraseServiceProvider(repo, converter, storage)
	userService := service.NewUserServiceProvider(repo)
	httpHandler := handler.NewHTTPHandler(phraseService, userService)

	e := handler.New()
	handler.SetupRoute(e, httpHandler)
	slog.Error("failed to run server", "err", e.Start(fmt.Sprintf(":%v", cfg.Server.Port)))
	// e.Start(fmt.Sprintf(":%v", cfg.Server.Port))
}
