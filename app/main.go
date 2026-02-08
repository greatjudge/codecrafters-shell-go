package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type handlerFunc func(cmd Command)

const (
	errorMsg = "command not found"
	errNotFoundInPATH

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

	for {
		fmt.Print("$ ")

		reader := bufio.NewReader(os.Stdin)

		raw, err := reader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}

		cmd, err := ParseCommand(raw)
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
	}
}

func echo(cmd Command) {
	if len(cmd.Args) < 2 {
		fmt.Printf("\n")

		return
	}

	fmt.Println(strings.Join(cmd.Args[1:], " "))
}

func runExecutableHandler(cmd Command) {
	cmdPath, err := lookPath(cmd.Name)
	if err != nil {
		fmt.Printf("%s: %s\n", cmd.Name, errorMsg)

		return
	}

	executable := exec.Command(cmdPath, cmd.Args[1:]...)
	executable.Stdout = os.Stdout
	executable.Stderr = os.Stderr
	executable.Args = cmd.Args // Set args to use original command name as args[0]

	err = executable.Run()
	if err != nil {
		fmt.Println("err: ", err)

		return
	}
}

func typeHandler(cmd Command) {
	for _, arg := range cmd.Args[1:] {
		if lo.Contains(builtIn, arg) {
			fmt.Printf("%v is a shell builtin\n", arg)

			return
		}

		absPath, err := lookPath(arg)
		if err != nil {
			fmt.Printf("%s not found\n", arg)

			return
		}

		fmt.Printf("%s is %s\n", arg, absPath)
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

func pwd(_ Command) {
	d, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println(d)
}

func cd(cmd Command) {
	var dir string

	switch {
	case len(cmd.Args) < 2:
		dir = "~"
	case len(cmd.Args) == 2:
		dir = cmd.Args[1]
	default:
		fmt.Printf("cd: too many arguments")

		return
	}

	dir = insertHomeDir(dir)

	err := os.Chdir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("cd: %v: No such file or directory\n", dir)

			return
		}

		fmt.Println(err)

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
