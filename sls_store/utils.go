package sls_store

import (
	"time"
)

func buildSearchingData(lookback time.Duration) (int64, int64) {
	currentTime := time.Now()
	to := currentTime.Unix()
	from := currentTime.Add(-1 * lookback).Unix()
	return from, to
}
