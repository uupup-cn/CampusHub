package repository

import "time"

type ReleaseLog struct {
	ID          int64
	Amount      int64
	FromPool    string
	ToPool      string
	OperatorID  int64
	Reason     *string
	IsAuto      bool
	CreatedAt   time.Time
}
