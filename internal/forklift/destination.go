package forklift

// Destination is an output destination where forklift will deliver it's results.
// It knows how to add records and "close".
type Destination interface {
	AddRecord(Record)
	Close()
}
