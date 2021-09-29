package forklift

import (
	"encoding/json"
	"time"
)

// today returns the current UTC date
func today() string {
	return time.Now().UTC().Format("2006-01-02")
}

// getYMDFromISOreturns the Y-m-d string from an ISO-formatted-ish date
// araddon/dateparse will be helpful here
func getYMDFromISO(d string) string {
	return d[:10]
}

func getYMDFromUnixUTC(ts float64) string {
	unixTime := time.Unix(int64(ts), 0)
	return unixTime.Format("2006-01-02")
}

func getYMDFromUnixNano(ts float64) string {
	unixTime := time.Unix(int64(ts)/1000, 0)
	return unixTime.Format("2006-01-02")
}

// jsonFromKey retrieves the string value from the root key provided
func jsonFromKey(key string, s string) string {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		panic(err)
	}
	return m[key].(string)
}
