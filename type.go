package gotwilio

// NewBoolean returns a boolean pointer value for a given boolean literal.
// This is important because for the Twilio API, booleans are really a ternary type, supporting true, false, or nil/null.
func NewBoolean(value bool) *bool {
	return &value
}

// Option contains a key/value pair to define optional request parameters.
type Option struct {
	Key   string
	Value string
}
