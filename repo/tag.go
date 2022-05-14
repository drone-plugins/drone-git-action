package repo

import "os/exec"

const defaultTagMessage = "Release"

// CreateTag creates a new tag
func CreateTag(tag string, msg string) *exec.Cmd {
	if msg == "" {
		msg = defaultTagMessage + " " + tag
	}

	cmd := exec.Command(
		"git",
		"tag",
		tag,
		"-m",
		msg,
	)

	return cmd
}
