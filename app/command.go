package main

import "strings"

type Command struct {
	Name    string
	RawArgs string
	Args    []string
}

func NewCommand(
	name string,
	rawArgs string,
	args []string,
) Command {
	return Command{
		Name:    name,
		RawArgs: rawArgs,
		Args:    args,
	}
}

func ParseCommand(raw string) Command {
	raw = strings.TrimSpace(raw)

	name, args := extractCommand(raw)

	return Command{
		Name:    name,
		RawArgs: args,
		Args:    parseArgs(name, args),
	}
}

func extractCommand(cmd string) (string, string) {
	parts := strings.SplitN(cmd, " ", 2)
	if len(parts) < 2 {
		return cmd, ""
	}

	return parts[0], parts[1]
}

func parseArgs(cmd, args string) []string {
	parsed := []string{cmd}

	if len(args) == 0 {
		return parsed
	}

	parsed = append(parsed, strings.Split(args, " ")...)

	return parsed
}
