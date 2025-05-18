package repository

import (
	"context"
	"fmt"
	"log/slog"
	"lotto-notifications/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetGames(ctx context.Context, independentOnly bool) ([]models.Game, error)
	GetGame(ctx context.Context, gameType string) (models.Game, error)
	GetResults(ctx context.Context, gameType string) ([]models.Result, error)
	GetNewestResult(ctx context.Context, gameType string) (models.Result, error)
	UpdateGames(ctx context.Context, games []models.Game) error
	InsertResults(ctx context.Context, results []models.Result) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetGames(ctx context.Context, independentOnly bool) ([]models.Game, error) {
	stmt := `SELECT * FROM games`
	if independentOnly {
		stmt += ` WHERE tied_to IS NULL`
	}
	games := []models.Game{}
	err := r.db.SelectContext(ctx, &games, stmt)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (r *repository) GetGame(ctx context.Context, gameType string) (models.Game, error) {
	stmt := `SELECT * FROM games WHERE type = ?`
	game := models.Game{}
	err := r.db.GetContext(ctx, &game, stmt, gameType)
	if err != nil {
		return models.Game{}, err
	}
	return game, nil
}

func (r *repository) GetResults(ctx context.Context, gameType string) ([]models.Result, error) {
	stmt := `SELECT * FROM results WHERE game_type = ?`
	results := []models.Result{}
	err := r.db.SelectContext(ctx, &results, stmt, gameType)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *repository) GetNewestResult(ctx context.Context, gameType string) (models.Result, error) {
	stmt := `SELECT * FROM results WHERE game_type = ? ORDER BY draw_date DESC LIMIT 1`
	result := models.Result{}
	err := r.db.GetContext(ctx, &result, stmt, gameType)
	if err != nil {
		return models.Result{}, err
	}
	return result, nil
}

func (r *repository) UpdateGames(ctx context.Context, games []models.Game) error {
	slog.Debug("Updating games", "games", len(games))
	stmt := `UPDATE games SET
		next_draw_date = :next_draw_date,
		closest_prize_value = :closest_prize_value,
		draws = :draws,
		coupon_price = :coupon_price,
		closest_prize_pool = :closest_prize_pool
	WHERE type = :type`

	trx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer trx.Rollback()

	// should be done in batch but sqlx has a bug with batched updates
	// not a big deal
	for _, game := range games {
		_, err := trx.NamedExecContext(ctx, stmt, game)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}
	err = trx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *repository) InsertResults(ctx context.Context, results []models.Result) error {
	stmt := `INSERT INTO results (draw_id, game_type, draw_date, prize_pool, prize_value)
		VALUES (:draw_id, :game_type, :draw_date, :prize_pool, :prize_value)`
	_, err := r.db.NamedExecContext(ctx, stmt, results)
	if err != nil {
		return err
	}

	return nil
}
