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
	Stderr           io.WriteCloser
	stdoutRedirected bool
	stderrRedirected bool
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
		filename, errFilename string
		outIDX, errIDX        = -1, -1
	)

	for i, a := range args {
		switch {
		case a == ">" || a == "1>":
			if i+1 >= len(args) {
				return Command{}, fmt.Errorf("expect filename after >")
			}

			filename = args[i+1]
			outIDX = i
		case a == "2>":
			if i+1 >= len(args) {
				return Command{}, fmt.Errorf("expect filename after 2>")
			}

			errFilename = args[i+1]
			errIDX = i
		}
	}

	filename = strings.TrimSpace(filename)
	errFilename = strings.TrimSpace(errFilename)

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

	if outIDX != -1 {
		cmd.Args = slices.Delete(cmd.Args, outIDX, outIDX+2)
	}

	if errFilename != "" {
		file, err := os.Create(errFilename)
		if err != nil {
			return Command{}, err
		}

		cmd.stderrRedirected = true
		cmd.Stderr = file
	} else {
		cmd.Stderr = os.Stderr
	}

	if errIDX != -1 {
		cmd.Args = slices.Delete(cmd.Args, errIDX, errIDX+2)
	}

	return cmd, nil
}
