package main

import "fmt"

type printer func(string)

func handleCommand(command string, cb printer) {
	cb(fmt.Sprintf("DEBUG: Got command %s", command))
}
