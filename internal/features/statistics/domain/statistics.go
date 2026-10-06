package domain

import "time"

type AccountStatistics struct {
	AccountID int64     `json:"account_id"`
	UserID    int64     `json:"user_id"`
	Month     time.Time `json:"month"`
	Income    float64   `json:"income"`
	Expense   float64   `json:"expense"`
	Profit    float64   `json:"profit"`
}

type UserStatistics struct {
	UserID  int64     `json:"user_id"`
	Month   time.Time `json:"month"`
	Income  float64   `json:"income"`
	Expense float64   `json:"expense"`
	Profit  float64   `json:"profit"`
}
