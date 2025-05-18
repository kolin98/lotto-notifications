-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS games (
    type                TEXT PRIMARY KEY,
    next_draw_date      TIMESTAMP DEFAULT NULL,
    closest_prize_value REAL DEFAULT NULL,
    draws               TEXT DEFAULT NULL,
    coupon_price        TEXT DEFAULT NULL,
    closest_prize_pool  TEXT DEFAULT NULL,
    tied_to             TEXT DEFAULT NULL REFERENCES games(type)
);

INSERT INTO games (type, tied_to) VALUES
    ('Lotto', NULL),
    ('LottoPlus', 'Lotto'),
    ('EuroJackpot', NULL),
    ('MultiMulti', NULL),
    ('MiniLotto', NULL),
    ('Kaskada', NULL),
    ('EkstraPensja', NULL)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS games;
-- +goose StatementEnd
