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

func ParseCommand(raw string) Command {
	raw = strings.TrimSpace(raw)

	name, args := extractCommand(raw)

	return Command{
		Name: name,
		Args: parseArgs(name, args),
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

	for _, arg := range strings.Split(args, " ") {
		arg = strings.Trim(arg, "'")
		arg = strings.Replace(arg, "''", "", -1)

		if arg != "" {
			parsed = append(parsed, arg)
		}
	}

	return parsed
}

func parseArgsV2(raw string) []string {
	var (
		args = []string{}
		arg  = []rune{}

		state = 1
	)

	for _, ch := range []rune(raw) {
		switch state {
		case 1:
			if ch != ' ' {
				state = 2
				arg = append(arg, ch)
			}
		case 2:
			if ch == ' ' {
				args = append(args, string(arg))
				arg = []rune{}
				state = 1
			} else {
				arg = append(arg, ch)
			}
		}
	}

	if state == 2 {
		args = append(args, string(arg))
	}

	return args
}
