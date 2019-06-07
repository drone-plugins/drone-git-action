package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// helper function to simply wrap os execte command.
func execute(cmd *exec.Cmd) error {
	fmt.Println("+", strings.Join(cmd.Args, " "))

	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// helper function returns true if directory dir is empty.
func isDirEmpty(dir string) bool {
	f, err := os.Open(dir)

	if err != nil {
		return true
	}

	defer f.Close()

	_, err = f.Readdir(1)
	return err == io.EOF
}
