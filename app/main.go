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
		"echo",
		"type",
		"pwd",
		"cd",
	}
)

func main() {
	for {
		fmt.Print("$ ")

		reader := bufio.NewReader(os.Stdin)

		raw, err := reader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}

		cmd := ParseCommand(raw)

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
	fmt.Println(cmd.RawArgs)
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
	if lo.Contains(builtIn, cmd.RawArgs) {
		fmt.Printf("%v is a shell builtin\n", cmd.RawArgs)

		return
	}

	absPath, err := lookPath(cmd.RawArgs)
	if err != nil {
		fmt.Printf("%s not found\n", cmd.RawArgs)

		return
	}

	fmt.Printf("%s is %s\n", cmd.RawArgs, absPath)
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
	dir := cmd.RawArgs

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
