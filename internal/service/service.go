package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lotto-notifications/internal/models"
	"lotto-notifications/internal/repository"
	"lotto-notifications/pkg/lotto"
)

var (
	ErrNoResultsAvailable     = errors.New("no results available")
	ErrNoResultsInDraw        = errors.New("no results in draw")
	ErrResultsNotYetAvailable = errors.New("results not yet available")
	ErrMainGameNotFound       = errors.New("main game not found")
)

type Service interface {
	UpdateAllGames(ctx context.Context) ([]models.Game, error)
	UpdateGame(ctx context.Context, gameType models.GameType) (models.Game, error)
	GetAndSaveNewestResults(ctx context.Context, gameType models.GameType, nextDrawDate time.Time) ([]models.Result, error)
}

type service struct {
	lottoClient lotto.Client
	repo        repository.Repository
}

func NewService(
	lottoClient lotto.Client,
	repo repository.Repository,
) Service {
	return &service{
		lottoClient: lottoClient,
		repo:        repo,
	}
}

func (s *service) UpdateAllGames(ctx context.Context) ([]models.Game, error) {
	games, err := s.repo.GetGames(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get all games: %w", err)
	}

	for idx, game := range games {
		gameInfo, err := s.lottoClient.GetGameInfo(ctx, string(game.GameType))
		if err != nil {
			return nil, fmt.Errorf("failed to get game info: %w", err)
		}
		game.NextDrawDate = &gameInfo.NextDrawDate
		game.ClosestPrizeValue = &gameInfo.ClosestPrizeValue
		game.Draws = &gameInfo.Draws
		game.CouponPrice = &gameInfo.CouponPrice
		game.ClosestPrizePool = &gameInfo.ClosestPrizePoolType
		games[idx] = game
	}

	err = s.repo.UpdateGames(ctx, games)
	if err != nil {
		return nil, fmt.Errorf("failed to update games: %w", err)
	}

	return games, nil
}

func (s *service) UpdateGame(ctx context.Context, gameType models.GameType) (models.Game, error) {
	gameInfo, err := s.lottoClient.GetGameInfo(ctx, string(gameType))
	if err != nil {
		return models.Game{}, fmt.Errorf("failed to get game info: %w", err)
	}

	game := models.Game{
		GameType:          gameType,
		NextDrawDate:      &gameInfo.NextDrawDate,
		ClosestPrizeValue: &gameInfo.ClosestPrizeValue,
		Draws:             &gameInfo.Draws,
		CouponPrice:       &gameInfo.CouponPrice,
		ClosestPrizePool:  &gameInfo.ClosestPrizePoolType,
	}

	err = s.repo.UpdateGames(ctx, []models.Game{game})
	if err != nil {
		return models.Game{}, fmt.Errorf("failed to update game: %w", err)
	}

	return game, nil
}

func (s *service) GetAndSaveNewestResults(
	ctx context.Context, gameType models.GameType, nextDrawDate time.Time,
) ([]models.Result, error) {
	draws, err := s.lottoClient.GetLastResults(ctx, string(gameType))
	if err != nil {
		return nil, fmt.Errorf("failed to get last results: %w", err)
	}

	if len(draws) == 0 {
		return nil, ErrNoResultsAvailable
	}

	mainGamePresent := false
	for _, result := range draws {
		if result.GameType == string(gameType) {
			mainGamePresent = true
			if result.DrawDate.Before(nextDrawDate) {
				return nil, ErrResultsNotYetAvailable
			}
		}
	}
	if !mainGamePresent {
		return nil, ErrMainGameNotFound
	}

	results := make([]models.Result, len(draws))
	for idx, draw := range draws {
		if len(draw.Results) == 0 {
			return nil, ErrNoResultsInDraw
		}
		results[idx] = models.Result{
			DrawID:         draw.DrawSystemID,
			GameType:       models.GameType(draw.GameType),
			DrawDate:       draw.DrawDate,
			Results:        draw.Results[0].Results,
			SpecialResults: draw.Results[0].SpecialResults,
			CreatedAt:      time.Now(),
		}
	}

	err = s.repo.InsertResults(ctx, results)
	if err != nil {
		return nil, fmt.Errorf("failed to insert results: %w", err)
	}

	return results, nil
}
