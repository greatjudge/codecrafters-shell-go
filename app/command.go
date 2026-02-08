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

	parsed = append(parsed, parseArgsV2(args)...)

	return parsed
}

const (
	StateOutside        = 1
	StateInArg          = 2
	StateInSingleQuotes = 3
)

func parseArgsV2(raw string) []string {
	raw = strings.Replace(raw, "''", "", -1)

	var (
		args  = []string{}
		arg   = []rune{}
		state = StateOutside
	)

	for _, ch := range []rune(raw) {
		switch state {
		case StateOutside:
			switch {
			case ch == ' ':
				continue
			case ch == '\'':
				state = StateInSingleQuotes
			default:
				arg = append(arg, ch)
				state = StateInArg
			}
		case StateInArg:
			switch {
			case ch == ' ':
				args = append(args, string(arg))
				arg = []rune{}
				state = StateOutside
			case ch == '\'':
				state = StateInSingleQuotes
			default:
				arg = append(arg, ch)
			}
		case StateInSingleQuotes:
			switch {
			case ch == '\'':
				args = append(args, string(arg))
				arg = []rune{}
				state = StateOutside
			default:
				arg = append(arg, ch)
			}
		}
	}

	if (state == 2 || state == 3) && len(arg) != 0 {
		args = append(args, string(arg))
	}

	return args
}
