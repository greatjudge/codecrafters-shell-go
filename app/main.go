package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	errorMsg = "command not found"

	exitCmd = "exit"
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

		switch cmd {
		case exitCmd:
			break repl
		default:
			fmt.Printf("%s: %s\n", cmd, errorMsg)
		}
	}
}
