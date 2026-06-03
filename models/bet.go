package models

import "time"

type Bet struct {
	ID        int
	UserID    int
	Username  string
	Amount    int
	Game      string
	Result    bool
	CreatedAt time.Time
}
