package time

import "time"

const (
	yyyyMMddHHmmss = "2006-01-02 15:04:05"
	//yyyy-MM-dd HH:mm:ss
)

const (
	KipleTime = "2006-01-02 15:04:05"
	KipleDate = "2006-01-02"
)

var KipleTimeZone = time.FixedZone("CST", 8*3600)

func GetNowTime() time.Time {
	return time.Now().In(KipleTimeZone)
}
