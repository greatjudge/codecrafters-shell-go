package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/command"
	"github.com/samber/lo"
)

type handlerFunc func(cmd *command.Command)

const (
	errorMsg = "command not found"

	exitCMD = "exit"
)

var (
	handlerByName = map[string]handlerFunc{
		"echo": echo,
		"type": typeHandler,
		"pwd":  pwd,
		"cd":   cd,
	}

	builtIn = []string{
		"exit",
	}
)

func main() {
	for k := range handlerByName {
		builtIn = append(builtIn, k)
	}

	completeItems := make([]readline.PrefixCompleterInterface, 0, len(builtIn))
	for _, c := range builtIn {
		completeItems = append(completeItems, readline.PcItem(c))
	}

	completer := readline.NewPrefixCompleter(
		completeItems...,
	)

	l, err := readline.NewEx(&readline.Config{
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	l.SetPrompt("$ ")

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		cmd, err := command.ParseCommand(line)
		if err != nil {
			log.Panic(err)
		}

		if cmd.Name == exitCMD {
			break
		}

		handler, ok := handlerByName[cmd.Name]
		if ok {
			handler(cmd)
		} else {
			runExecutableHandler(cmd)
		}

		cmd.Close()
	}
}

func echo(cmd *command.Command) {
	if len(cmd.Args) < 2 {
		fmt.Fprintf(cmd.Stdout, "\n")

		return
	}

	fmt.Fprintln(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
}

func runExecutableHandler(cmd *command.Command) {
	cmdPath, err := lookPath(cmd.Name)
	if err != nil {
		fmt.Printf("%s: %s\n", cmd.Name, errorMsg)

		return
	}

	executable := exec.Command(cmdPath, cmd.Args[1:]...)
	executable.Stdout = cmd.Stdout
	executable.Stderr = cmd.Stderr
	executable.Args = cmd.Args // Set args to use original command name as args[0]

	executable.Run()
}

func typeHandler(cmd *command.Command) {
	for _, arg := range cmd.Args[1:] {
		if lo.Contains(builtIn, arg) {
			fmt.Fprintf(cmd.Stdout, "%v is a shell builtin\n", arg)

			return
		}

		absPath, err := lookPath(arg)
		if err != nil {
			fmt.Fprintf(cmd.Stdout, "%s not found\n", arg)

			return
		}

		fmt.Fprintf(cmd.Stdout, "%s is %s\n", arg, absPath)
	}
}

// find path to executable command
func lookPath(file string) (string, error) {
	path := os.Getenv("PATH")

	pathes := strings.Split(path, string(os.PathListSeparator))

	var foundPath string

	for _, p := range pathes {
		err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.Name() != file {
				return nil
			}

			info, err := os.Lstat(path)
			if err != nil {
				return nil
			}

			if IsExecutedAny(info.Mode()) {
				foundPath = path

				return filepath.SkipAll
			}

			return nil
		})
		if err != nil {
			return "", err
		}

		if foundPath != "" {
			break
		}
	}

	if foundPath == "" {
		return "", fmt.Errorf("not found")
	}

	return foundPath, nil
}

// mode - bitmask rwxrwxrwx (in octal 0777)
// 3 user, 3 group, 3 other
// & 001001001 (in octal in go: 0111)
func IsExecutedAny(mode fs.FileMode) bool {
	return mode&0111 != 0
}

func pwd(cmd *command.Command) {
	d, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(cmd.Stdout, err.Error())

		return
	}

	fmt.Fprintln(cmd.Stdout, d)
}

func cd(cmd *command.Command) {
	var dir string

	switch {
	case len(cmd.Args) < 2:
		dir = "~"
	case len(cmd.Args) == 2:
		dir = cmd.Args[1]
	default:
		fmt.Fprintf(cmd.Stdout, "cd: too many arguments")

		return
	}

	dir = insertHomeDir(dir)

	err := os.Chdir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(cmd.Stdout, "cd: %v: No such file or directory\n", dir)

			return
		}

		fmt.Fprintln(cmd.Stdout, err.Error())

		return
	}
}

func insertHomeDir(dir string) string {
	if !strings.Contains(dir, "~") {
		return dir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	return strings.Replace(dir, "~", homeDir, 1)
}
