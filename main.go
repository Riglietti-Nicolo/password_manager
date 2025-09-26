// package passmanager

package main

import (
	"fmt"
	"github.com/Riglietti-Nicolo/password_manager/cmd"
	"log"
)

func main() {
	fmt.Printf("\n\nBenvenuto in password manager!\n\n")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}