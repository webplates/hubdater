package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func credentials() (string, string) {
	username := strings.TrimSpace(os.Getenv("DOCKERHUB_USERNAME"))
	password := strings.TrimSpace(os.Getenv("DOCKERHUB_PASSWORD"))

	if username == "" || password == "" {
		reader := bufio.NewReader(os.Stdin)

		if username == "" {
			fmt.Print("Enter Username: ")
			username, err := reader.ReadString('\n')

			if err != nil {
				panic(err)
			}

			username = strings.TrimSpace(username)
		}

		if password == "" {
			fmt.Print("Enter Password: ")
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))

			if err != nil {
				panic(err)
			}

			password = strings.TrimSpace(string(bytePassword))
		}
	}

	return username, password
}
