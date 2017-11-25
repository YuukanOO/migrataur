package migrataur

import (
	"strconv"
	"time"
)

func currentTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
