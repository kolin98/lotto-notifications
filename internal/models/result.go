package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type IntSlice []int

type Result struct {
	DrawID         uint      `db:"draw_id"`
	GameType       GameType  `db:"game_type"`
	DrawDate       time.Time `db:"draw_date"`
	Results        IntSlice  `db:"results"`
	SpecialResults IntSlice  `db:"special_results"`
	CreatedAt      time.Time `db:"created_at"`
}

func (s IntSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	strs := make([]string, len(s))
	for i, num := range s {
		strs[i] = strconv.Itoa(num)
	}
	return strings.Join(strs, ","), nil
}

func (s *IntSlice) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan Results: expected string, got %T", value)
	}

	if str == "" {
		*s = []int{}
		return nil
	}

	parts := strings.Split(str, ",")
	*s = make([]int, len(parts))
	for i, part := range parts {
		num, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return fmt.Errorf("failed to scan Results: invalid number %q", part)
		}
		(*s)[i] = num
	}
	return nil
}
