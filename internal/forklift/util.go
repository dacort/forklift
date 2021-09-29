package forklift

import (
	"encoding/json"
	"time"
	"fmt"
	"strconv"
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

func getYMDFromUnixUTC(ts string) string {
	i, _ := strconv.ParseInt(ts, 10, 64)
	unixTime := time.Unix(i, 0)
	return unixTime.Format("2006-01-02")
}

func getYMDFromUnixNano(ts string) string {
	i, _ := strconv.ParseInt(ts, 10, 64)
	unixTime := time.Unix(i/1000, 0)
	return unixTime.Format("2006-01-02")
}

// jsonFromKey retrieves the string value from the root key provided
func jsonFromKey(key string, s string) string {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		panic(err)
	}
	v := m[key]
	// fmt.Printf("%v", m[key])
	// return fmt.Sprintf("%q", m[key])

	// Figured we would have to start adding type detection at some point...
	switch v.(type) {
	case string:
		return v.(string)
	case float64:
		if (v.(float64) == float64(int64(v.(float64)))) {
			return fmt.Sprintf("%d", int64(v.(float64)))
		} else {
			return fmt.Sprintf("%f", v)
		}
	default:
		return fmt.Sprintf("%v", v)
	}
}
