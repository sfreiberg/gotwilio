package gotwilio

// Boolean is a custom ternary nullable bool type
type Boolean *bool

var (
	// True represents "true" in our ternary optional bool
	True Boolean
	// False represents "false" in our ternary optional bool
	False Boolean
)

// initialize our Boolean type
func init() {
	// we have to do this because we can't reference the address of a boolean literal
	// e.g. &true doesn't work
	a := true
	b := false
	True = &a
	False = &b
}
