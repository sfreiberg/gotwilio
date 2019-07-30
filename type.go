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
	a := true
	b := false
	True = &a
	False = &b
}
