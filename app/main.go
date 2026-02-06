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

type handlerFunc func(args []string)

const (
	errorMsg = "command not found"
	errNotFoundInPATH

	exitCMD = "exit"
	echoCMD = "echo"
	typeCMD = "type"
	pwdCMD  = "pwd"
)

var (
	builtInCommands = []string{
		exitCMD,
		echoCMD,
		typeCMD,
		pwdCMD,
	}
)

func main() {
repl:
	for {
		fmt.Print("$ ")

		reader := bufio.NewReader(os.Stdin)

		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}

		cmd = strings.TrimSpace(cmd)

		cmd, args := extractCommand(cmd)

		switch cmd {
		case exitCMD:
			break repl
		case echoCMD:
			echo(args)
		case typeCMD:
			typeFunc(args)
		case pwdCMD:
			pwd()
		default:
			runExecutableHandler(cmd, args)
		}
	}
}

func extractCommand(cmd string) (string, string) {
	parts := strings.SplitN(cmd, " ", 2)
	if len(parts) < 2 {
		return cmd, ""
	}

	return parts[0], parts[1]
}

func echo(args string) {
	fmt.Println(args)
}

func runExecutableHandler(cmd string, args string) {
	cmdPath, err := lookPath(cmd)
	if err != nil {
		fmt.Printf("%s: %s\n", cmd, errorMsg)

		return
	}

	cmdArgs := parseArgs(cmd, args)

	executable := exec.Command(cmdPath, cmdArgs[1:]...)
	executable.Stdout = os.Stdout
	executable.Stderr = os.Stderr
	executable.Args = cmdArgs // Set args to use original command name as args[0]

	err = executable.Run()
	if err != nil {
		fmt.Println("err: ", err)

		return
	}
}

func parseArgs(cmd, args string) []string {
	parsed := []string{cmd}

	if len(args) == 0 {
		return parsed
	}

	parsed = append(parsed, strings.Split(args, " ")...)

	return parsed
}

func typeFunc(args string) {
	if lo.Contains(builtInCommands, args) {
		fmt.Printf("%v is a shell builtin\n", args)

		return
	}

	absPath, err := lookPath(args)
	if err != nil {
		fmt.Printf("%s not found\n", args)

		return
	}

	fmt.Printf("%s is %s\n", args, absPath)
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

func pwd() {
	d, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println(d)
}
