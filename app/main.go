package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
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

		cmd = cmd[:len(cmd)-1]

		cmd, args := extractCommand(cmd)

		switch cmd {
		case exitCMD:
			break repl
		case echoCMD:
			echo(args)
		case typeCMD:
			typeFunc(args)
		default:
			fmt.Printf("%s: %s\n", cmd, errorMsg)
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

func typeFunc(args string) {
	if lo.Contains(builtInCommands, args) {
		fmt.Printf("%v is a shell builtin\n", args)

		return
	}

	absPath, err := exec.LookPath(args)
	if err != nil {
		fmt.Printf("%s not found\n", args)

		return
	}

	fmt.Printf("%s is %s\n", args, absPath)
}
