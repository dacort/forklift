package forklift

import (
	"os"
)

// StdOutPassthrough just passes all input data through to stdout
type StdOutPassthrough struct{}

// AddRecord determines where each record should get written to
// And then sends a Write command
func (sop StdOutPassthrough) AddRecord(r Record) {
	os.Stdout.Write(r.GetBytes())
}

// Close is a no-op for stdout
func (sop StdOutPassthrough) Close() {}

// NewStdOutPassthrough initializes the S3 Uploader
func NewStdOutPassthrough() StdOutPassthrough {
	return StdOutPassthrough{}
}
