package entity

import "time"

type Filter struct {
	DateFrom        time.Time
	DateTo          time.Time
	Categories      []int
	LabIDs          []int
	RecognizeStatus []string
	Tags            []string
	Sources         []string
	ValueFlags      []string
}
