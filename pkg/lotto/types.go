package lotto

import "time"

type GameInfo struct {
	GameType             string    `json:"gameType"`
	NextDrawDate         time.Time `json:"nextDrawDate"`
	ClosestPrizeValue    float64   `json:"closestPrizeValue"`
	Draws                string    `json:"draws"`
	CouponPrice          string    `json:"couponPrice"`
	ClosestPrizePoolType string    `json:"closestPrizePoolType"`
}

type Draw struct {
	DrawSystemID         uint      `json:"drawSystemId"`
	DrawDate             time.Time `json:"drawDate"`
	GameType             string    `json:"gameType"`
	MultiplierValue      uint      `json:"multiplierValue"`
	Results              []Result  `json:"results"`
	ShowSpecialResults   bool      `json:"showSpecialResults"`
	IsNewEuroJackpotDraw bool      `json:"isNewEuroJackpotDraw"`
}

type Result struct {
	DrawSystemID   uint      `json:"drawSystemId"`
	DrawDate       time.Time `json:"drawDate"`
	GameType       string    `json:"gameType"`
	Results        []int     `json:"resultsJson"`
	SpecialResults []int     `json:"specialResults"`
}
