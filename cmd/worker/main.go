package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"lotto-notifications/internal/config"
	"lotto-notifications/internal/database"
	"lotto-notifications/internal/logging"
	"lotto-notifications/internal/repository"
	"lotto-notifications/internal/service"
	"lotto-notifications/internal/worker"
	"lotto-notifications/pkg/lotto"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logging.Init(cfg.Environment)

	err = database.Initialize(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		return
	}
	defer database.Close()

	db, err := database.GetDB()
	if err != nil {
		slog.Error("Failed to get database", "error", err)
		return
	}

	lottoClient := lotto.NewClient(cfg.LottoAPIKey)
	repo := repository.NewRepository(db)
	service := service.NewService(lottoClient, repo)

	games, err := service.UpdateAllGames(context.Background())
	if err != nil {
		slog.Error("Failed to update all games", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	for _, game := range games {
		w, err := worker.NewResultsWorker(game, repo, service)
		if err != nil {
			slog.Error("Failed to create worker", "error", err)
			return
		}

		wg.Add(1)
		go func(w worker.ResultsWorker) {
			defer wg.Done()
			w.Run(ctx)
		}(w)
	}

	<-ctx.Done()
	slog.Info("Shutting down gracefully...")

	wg.Wait()
	slog.Info("Shutdown complete")
}
