package main

import (
	"example-interpreter/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bonvenon, %s. This is the Monkey programming language.\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
