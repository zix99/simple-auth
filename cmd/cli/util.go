package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func readPassword(prompt string) (string, error) {
	originalState, _ := terminal.GetState(int(syscall.Stdin))
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		_, ok := <-c
		if ok {
			terminal.Restore(int(syscall.Stdin), originalState)
			fmt.Println("Aborting on ctrl+c")
			os.Exit(0)
		}
	}()

	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()

	close(c)

	return string(bytePassword), nil
}

func readPasswordTwice() (string, error) {
	password1, err1 := readPassword("Password: ")
	if err1 != nil {
		return "", err1
	}
	password2, err2 := readPassword("Re-enter: ")
	if err2 != nil {
		return "", err2
	}

	if password1 == "" {
		return "", errors.New("Blank")
	}
	if password1 != password2 {
		return "", errors.New("Mismatched")
	}
	return password1, nil
}
