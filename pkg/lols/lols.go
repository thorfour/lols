package lols

var commands = map[string]func([]string) (string, error){
	"new": newLol,
}

// Handle responds to given commands
func Handle(cmd string, args []string) (string, error) {
	f, ok := commands[cmd]
	if !ok {
		return getLol(args)
	}

	return f(args)
}

func newLol(args []string) (string, error) { return "", nil }
func getLol(args []string) (string, error) { return "", nil }
