package forklift

import (
	"encoding/json"
	"math/rand"
	"time"
)

// today returns the current UTC date
func today() string {
	return time.Now().UTC().Format("2006-01-02")
}

// now returns the current UTC time
func now() time.Time {
	return time.Now().UTC()
}

// getYMDFromISOreturns the Y-m-d string from an ISO-formatted-ish date
// araddon/dateparse will be helpful here
func getYMDFromISO(d string) string {
	return d[:10]
}

// strftime returns the current time in the format specified
func strftime(format string, t time.Time) string {
	return t.Format(format)
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

// uniq returns a unique string of the specified length
func unique(length int) string {
	return randSeq(length)
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
