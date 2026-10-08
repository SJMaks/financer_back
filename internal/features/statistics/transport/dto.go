package transport

import "time"

type AccountStatisticsDTO struct {
	AccountID int64     `json:"account_id" example:"1"`
	UserID    int64     `json:"user_id" example:"1"`
	Month     time.Time `json:"month" example:"2026-01-01T00:00:00Z"`
	Income    float64   `json:"income" example:"500"`
	Expense   float64   `json:"expense" example:"200"`
	Profit    float64   `json:"profit" example:"300"`
}

type UserStatisticsDTO struct {
	UserID  int64     `json:"user_id" example:"1"`
	Month   time.Time `json:"month" example:"2026-01-01T00:00:00Z"`
	Income  float64   `json:"income" example:"1000"`
	Expense float64   `json:"expense" example:"1200"`
	Profit  float64   `json:"profit" example:"-200"`
}

type UserStatisticsSliceDTO struct {
	Months []UserStatisticsDTO `json:"months"`
}

type AccountStatisticsSliceDTO struct {
	Months []AccountStatisticsDTO `json:"months"`
}
