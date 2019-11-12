package karma

// Command my own string type for commands (think of it as an enum)
type Command string

const (
	help   string = "help"
	me     string = "me"
	status string = "status"
	add    string = "++"
	sub    string = "--"
	top    string = "top"
)

// Commands a set of the support commands by this processor
var Commands = map[string]struct{}{
	help:   struct{}{},
	me:     struct{}{},
	status: struct{}{},
	add:    struct{}{},
	sub:    struct{}{},
	top:    struct{}{},
}

// ProcConfig processor config object to contain all of these customizations
// SingleLimit one time karma swings are capped at 5 (default)
// DailyLimit this is the default daily limit for giving/ taking karma
// used by top function as guard rails
// used by top function as guard rails
type ProcConfig struct {
	SingleLimit    int
	DailyLimit     int
	TopUserDefault int
	TopUserMax     int
}

// DefaultConfig default settings for the Processor
var DefaultConfig ProcConfig = ProcConfig{
	SingleLimit:    5,
	DailyLimit:     25,
	TopUserDefault: 3,
	TopUserMax:     10,
}
