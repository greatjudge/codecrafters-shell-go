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

const (
	errorMsg = "command not found"

	exitCMD = "exit"
	echoCMD = "echo"
	typeCMD = "type"
)

var (
	builtInCommands = []string{
		exitCMD,
		echoCMD,
		typeCMD,
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
	// absPath, err := lookPath(cmd)
	// if err != nil {
	// 	fmt.Printf("%s: %s\n", cmd, errorMsg)

	// 	return
	// }

	cmdArgs := parseArgs(args)

	executable := exec.Command(cmd, cmdArgs...)

	out, err := executable.Output()
	if err != nil {
		fmt.Println("err: ", err)

		return
	}

	fmt.Print(string(out))
}

func parseArgs(args string) []string {
	if len(args) == 0 {
		return []string{}
	}

	return strings.Split(args, " ")
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
