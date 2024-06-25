package command

// Interface for all cli commands.
type CliCommand interface {
	// Command name.
	CommandName() string
	// Command logic.
	Execute(args string)
}
