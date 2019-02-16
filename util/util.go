package util

import (
	"fmt"
	"time"
)


func ElapsedTime(start time.Time, end time.Time) string {
	return fmt.Sprintf("%v", end.Sub(start))
}
