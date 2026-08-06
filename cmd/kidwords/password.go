// License: MIT Open Source
// Copyright (c) Joe Linoff 2016
// Wrap around golang.org/x/crypto/ssh/terminal to handle ^C interrupts based on a suggestion by Konstantin Shaposhnikov in
// this thread: https://groups.google.com/forum/#!topic/golang-nuts/kTVAbtee9UA.
// Correctly resets terminal echo after ^C interrupts.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

func scanPassword(prompt string) ([]byte, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Restore state in the event of an interrupt.
	// CITATION: Konstantin Shaposhnikov - https://groups.google.com/forum/#!topic/golang-nuts/kTVAbtee9UA
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-c
		os.Exit(1)
	}()

	// Now get the password.
	fmt.Print(prompt)
	p, err := term.ReadPassword(syscall.Stdin)
	fmt.Println("")
	if err != nil {
		return nil, err
	}

	// Stop looking for ^C on the channel.
	signal.Stop(c)
	return p, nil
}
