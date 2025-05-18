-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS results (
    draw_id         INTEGER NOT NULL,
    game_type       TEXT NOT NULL REFERENCES games(type),
    draw_date       TIMESTAMP NOT NULL,
    results         TEXT NOT NULL,
    special_results TEXT DEFAULT NULL,
    created_at      TIMESTAMP NOT NULL,
    PRIMARY KEY (draw_id, game_type)
);
CREATE INDEX IF NOT EXISTS idx_results_draw_date ON results (draw_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_results_draw_date;
DROP TABLE IF EXISTS results;
-- +goose StatementEnd
