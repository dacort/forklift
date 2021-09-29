package forklift

import (
    "testing"
)

func TestGetYMDFromISO(t *testing.T) {
    isotime := "2021-09-29T20:32:37+00:00"
    expected := "2021-09-29"
    ymd := getYMDFromISO(isotime)
    if ymd != expected {
        t.Errorf("Error converting ISO to Y-m-d, expected: %s, got: %s", expected, ymd)
    }
}

func TestGetYMDFromUnixUTC(t *testing.T) {
    ts := float64(1601305234)
    expected := "2020-09-28"
    ymd := getYMDFromUnixUTC(ts)

    if ymd != expected {
        t.Errorf("Error converting Unix timestamp to Y-m-d, expected: %s, got: %s", expected, ymd)
    }
}

func TestGetYMDFromUnixNano(t *testing.T) {
    ts := float64(1601305234262)
    expected := "2020-09-28"
    ymd := getYMDFromUnixNano(ts)

    if ymd != expected {
        t.Errorf("Error converting Unix timestamp to Y-m-d, expected: %s, got: %s", expected, ymd)
    }
}

func TestJSONFromKeyForString(t *testing.T) {
    testJSON := `{"user": "damon", "ts": 1601305234, "tsnano": 1601305234262}`
    result := jsonFromKey("user", testJSON)
    if result != "damon" {
        t.Errorf("Error getting JSON value, expected: %s, got: %s", "damon", result)
    }
}

func TestJSONFromKeyForNumber(t *testing.T) {
    testJSON := `{"user": "damon", "ts": 1601305234, "tsnano": 1601305234262}`
    result := jsonFromKey("ts", testJSON)
    if result != "1601305234" {
        t.Errorf("Error getting JSON value, expected: %s, got: %s", "1601305234", result)
    }
}