package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	errorMsg = "command not found"

	exitCMD = "exit"
	echoCMD = "echo"
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
