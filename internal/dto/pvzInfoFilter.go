package dto

import "time"

type PvzInfoFilterDto struct {
	StartDate      time.Time
	StartDateGiven bool
	EndDate        time.Time
	EndDateGiven   bool
	Page           int
	Limit          int
}

func (d PvzInfoFilterDto) isValidPage() bool {
	return d.Page >= 1
}

func (d PvzInfoFilterDto) isValidLimit() bool {
	return d.Limit >= 1 && d.Limit <= 30
}

func (d PvzInfoFilterDto) IsValidParams() bool {
	return d.isValidLimit() && d.isValidPage()
}
