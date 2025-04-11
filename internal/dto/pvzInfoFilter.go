package dto

import "time"

type PvzInfoFilter struct {
	StartDate time.Time
	EndDate   time.Time
	Page      int
	Limit     int
}
