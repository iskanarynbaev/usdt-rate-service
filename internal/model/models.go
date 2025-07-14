package models

import "time"

type Rate struct {
	Ask       float64
	Bid       float64
	Timestamp time.Time
}
