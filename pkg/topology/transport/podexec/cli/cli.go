package cli

type CLI interface {
	CreateCommand(lua string, args ...any) (*Command, error)
	Unmarshal(res string, target any) error
}

type Command struct {
	Command []string
	StdIn   string
}
