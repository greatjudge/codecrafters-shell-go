package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	errorMsg = "command not found"
)

func main() {

	for {
		fmt.Print("$ ")

		reader := bufio.NewReader(os.Stdin)

		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}

		fmt.Printf("%s: %s\n", cmd[:len(cmd)-1], errorMsg)
	}
}
