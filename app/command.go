package main

import (
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func NewCommand(
	name string,
	args []string,
) Command {
	return Command{
		Name: name,
		Args: args,
	}
}

func ParseCommand(raw string) (Command, error) {
	raw = strings.TrimSpace(raw)

	args, err := parseArgs(raw)
	if err != nil {
		return Command{}, err
	}

	if len(args) == 0 {
		return Command{}, nil
	}

	return Command{
		Name: args[0],
		Args: args,
	}, nil
}
