package forklift

import (
	"bytes"
	"log"
	"text/template"
)

// Record is an individual record that will get routed to an output destination.
// We can get ra bytes from it and each record can format its output path.
type Record interface {
	FormatPath(string) string
	GetBytes() []byte
}

// JSONRecord parses arbitrary JSON strings
type JSONRecord struct {
	Raw      string
	parsed   map[string]interface{}
	template *template.Template
}

// JSONValue returns the string value from the provided key
func (jr JSONRecord) JSONValue(key string) string {
	return jsonFromKey(key, jr.Raw)
}

// FormatPath formats the parent output path template into an actual location
// We also currently parse the JSON here if it hasn't been already
func (jr JSONRecord) FormatPath(uri string) string {
	var formattedPath bytes.Buffer

	// Can we do this elsewhere?
	// Also how can we change this so some commands (like today) don't execute for **EVERY** row?
	// Nm, I think if we want a standard prefix that's just part of the original URL
	funcMap := template.FuncMap{
		"json": jr.JSONValue,
		// "jsonFromKey":   jsonFromKey,
		"today":         today,
		"getYMDFromISO": getYMDFromISO,
		"now":           now,
		"strftime":      strftime,
		"unique":        unique,
	}
	tmpl, err := template.New("record").Funcs(funcMap).Parse(uri)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&formattedPath, jr.Raw)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

	return formattedPath.String()
}

// GetBytes returns the bytes for the raw record
func (jr JSONRecord) GetBytes() []byte {
	return []byte(jr.Raw + "\n")
}
