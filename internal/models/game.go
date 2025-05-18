package models

import "time"

type Game struct {
	GameType          GameType   `db:"type"`
	NextDrawDate      *time.Time `db:"next_draw_date"`
	ClosestPrizeValue *float64   `db:"closest_prize_value"`
	Draws             *string    `db:"draws"`
	CouponPrice       *string    `db:"coupon_price"`
	ClosestPrizePool  *string    `db:"closest_prize_pool"`
	TiedTo            *string    `db:"tied_to"`
}
