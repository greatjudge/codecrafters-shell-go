package command

import (
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/tokenizer"
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
		defer cmd.Stdout.Close()
	}

	if cmd.stderrRedirected {
		defer cmd.Stderr.Close()
	}

	return nil
}

func ParseCommand(raw string) (*Command, error) {
	raw = strings.TrimSpace(raw)

	tokens, err := tokenizer.Tokenize(raw)
	if err != nil {
		return nil, err
	}

	args, redirectStdout, redirectStderr := splitTokens(tokens)

	if len(args) == 0 {
		return &Command{}, nil
	}

	cmd := &Command{
		Name:   args[0],
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = redirectStreams(cmd, redirectStdout, redirectStderr)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func splitTokens(tokens []tokenizer.Token) (args []string, redirectStdout, redirectStderr Redirection) {
	args = []string{}

	for _, t := range tokens {
		switch t.Type {
		case tokenizer.RedirectTokenType, tokenizer.AppendRedirectTokenType:
			redirectStdout = NewRedirectionFromToken(t)
		case tokenizer.RedirectErrTokenType:
			redirectStderr = NewRedirectionFromToken(t)
		case tokenizer.WordTokenType:
			args = append(args, t.Value)
		}
	}

	return args, redirectStdout, redirectStderr
}

func redirectStreams(cmd *Command, redirectStdout, redirectStderr Redirection) error {
	for _, r := range []Redirection{redirectStdout, redirectStderr} {
		if r.IsEmpty() {
			continue
		}

		var flag int

		switch r.WriteType {
		case RewriteRedirectionWriteType:
			flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
		case AppendRedirectionWriteType:
			flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
		}

		file, err := os.OpenFile(r.Filename, flag, 0666)
		if err != nil {
			return err
		}

		switch r.StreamType {
		case StdoutStreamType:
			cmd.stdoutRedirected = true
			cmd.Stdout = file
		case StderrStreamType:
			cmd.stderrRedirected = true
			cmd.Stderr = file
		}
	}

	return nil
}
