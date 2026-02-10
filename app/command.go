package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

const (
	redirectArg = ">"
)

type Command struct {
	Name             string
	Args             []string
	Stdout           io.WriteCloser
	stdoutRedirected bool
}

func (cmd Command) Close() error {
	if cmd.stdoutRedirected {
		return cmd.Stdout.Close()
	}

	return nil
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

	cmd := Command{
		Name: args[0],
		Args: args,
	}

	var (
		filename       string
		redirectArgIdx = 0
		deleteTwo      = false
	)

	for redirectArgIdx = 0; redirectArgIdx < len(args); redirectArgIdx++ {
		a := args[redirectArgIdx]

		idx := strings.Index(a, redirectArg)
		if idx == -1 {
			continue
		}

		// 1. > filename
		// 2. 1> filename
		// 3. >filename
		// 4. 1>filename

		if idx+1 >= len(a) {
			if redirectArgIdx+1 < len(args) {
				filename = args[redirectArgIdx+1]
				deleteTwo = true
			} else {
				return Command{}, fmt.Errorf("> expect filename after it")
			}
		} else {
			filename = a[idx+1:]
		}

		break
	}

	filename = strings.TrimSpace(filename)

	if filename != "" {
		file, err := os.Create(filename)
		if err != nil {
			return Command{}, err
		}

		cmd.stdoutRedirected = true
		cmd.Stdout = file
	} else {
		cmd.Stdout = os.Stdout
	}

	if redirectArgIdx < len(args) {
		j := redirectArgIdx + 1
		if deleteTwo {
			j = redirectArgIdx + 2
		}

		args = slices.Delete(args, redirectArgIdx, j)
	}

	return cmd, nil
}
