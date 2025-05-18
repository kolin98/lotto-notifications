package worker

import (
	"context"
	"log/slog"
	"time"

	"lotto-notifications/internal/models"
	"lotto-notifications/internal/repository"
	"lotto-notifications/internal/service"
)

type ResultsWorker interface {
	Run(ctx context.Context)
}

type resultsWorker struct {
	game    models.Game
	repo    repository.Repository
	service service.Service
}

func NewResultsWorker(
	game models.Game,
	repo repository.Repository,
	service service.Service,
) (ResultsWorker, error) {
	if game.TiedTo != nil {
		return nil, ErrGameNotCheckable
	}
	if game.NextDrawDate == nil {
		return nil, ErrGameInfoNotSet
	}

	return &resultsWorker{
		game:    game,
		repo:    repo,
		service: service,
	}, nil
}

func (w *resultsWorker) Run(ctx context.Context) {
	slog.Debug(
		"Running worker",
		"game", w.game.GameType,
		"nextDrawDate", w.game.NextDrawDate,
	)
	now := time.Now()
	timeUntilNextDraw := w.game.NextDrawDate.Sub(now)

	if timeUntilNextDraw > 0 {
		slog.Info("Waiting for next draw",
			"game", w.game.GameType,
			"nextDrawDate", w.game.NextDrawDate,
			"timeUntilNextDraw", timeUntilNextDraw,
		)
		select {
		case <-ctx.Done():
			return
		case <-time.After(timeUntilNextDraw):
			w.do(ctx)
		}
	}
}

func (w *resultsWorker) do(ctx context.Context) {
	backoffDuration := 5 * time.Minute
	maxBackoff := 30 * time.Minute

	waitWithBackoff := func() bool {
		select {
		case <-ctx.Done():
			return false
		case <-time.After(backoffDuration):
			backoffDuration *= 2
			if backoffDuration > maxBackoff {
				backoffDuration = maxBackoff
			}
			return true
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			game, err := w.service.UpdateGame(ctx, w.game.GameType)
			if err != nil {
				slog.Error("Failed to update game", "error", err)
			}
			if err != nil || game.NextDrawDate == w.game.NextDrawDate {
				// we continue and wait for success and results to be available
				if !waitWithBackoff() {
					return
				}
				continue
			}

			_, err = w.service.GetAndSaveNewestResults(ctx, w.game.GameType, *game.NextDrawDate)
			if err != nil {
				if err != service.ErrResultsNotYetAvailable {
					slog.Error("Failed to get and save newest results", "error", err)
				}
				// we continue and wait for success and results to be available
				if !waitWithBackoff() {
					return
				}
				continue
			}

			// we successfully updated the game and saved the results
			// now the worker will wait for the next draw
			slog.Info("Successfully updated game and saved results", "game", game.GameType)
			w.game = game
			return
		}
	}
}
