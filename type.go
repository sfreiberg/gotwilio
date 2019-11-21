package gotwilio

// NewBoolean returns a boolean pointer value for a given boolean literal.
// This is important because for the Twilio API, booleans are really a ternary type, supporting true, false, or nil/null.
func NewBoolean(value bool) *bool {
	return &value
}
