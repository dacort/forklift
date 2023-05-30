package forklift

import "testing"

func TestFormatPath(t *testing.T) {
	jr := JSONRecord{Raw: `{"damon": "true", "dob": "2021-03-20"}`}
	jr2 := JSONRecord{Raw: `{"damon": "true", "dob": "2021-03-21"}`}

	// Test a static path
	static := "s3://bucket-name/prefix"
	jsonReplace := `s3://bucket-name/prefix/date={{json "dob"}}/data.json`
	o := jr.FormatPath(static)
	if o != static {
		t.Fatalf("Static failed, expected %s, got %s", static, o)
	}

	// Test JSON field parsing
	o = jr.FormatPath(jsonReplace)
	if o != "s3://bucket-name/prefix/date=2021-03-20/data.json" {
		t.Fatalf("JSON extract failed, expected %s, got %s", "s3://bucket-name/prefix/date=2021-03-20/data.json", o)
	}
	o = jr2.FormatPath(jsonReplace)
	if o != "s3://bucket-name/prefix/date=2021-03-21/data.json" {
		t.Fatalf("JSON extract failed, expected %s, got %s", "s3://bucket-name/prefix/date=2021-03-21/data.json", o)
	}

	// Test one-time field parsing
	onetime := "s3://bucket-name/prefix/{{unique 8}}/"
	o1 := jr.FormatPath(onetime)
	o2 := jr.FormatPath(onetime)
	if o1 != o2 {
		t.Fatalf("One-time field failed, expected %s, got %s", o1, o2)
	}
}
