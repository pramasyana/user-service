package helper

import (
	"time"
)

// interval is number of days since epoch
// sample 16463
func DateSinceEpoch(interval int32) time.Time {
	return time.Unix(int64(interval)*24*60*60, 0)
}
