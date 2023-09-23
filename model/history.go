package model

import "time"

const (
	DateFormat = "01-02-2006"
	TimeFormat = "15:04"
)

type History struct {
	CurrentDate string
	Length      string
	Start       time.Time
	End         time.Time
	Task        string
}
