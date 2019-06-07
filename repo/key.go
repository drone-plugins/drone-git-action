package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const netrcFile = `
machine %s
login %s
password %s
`

const configFile = `
Host *
StrictHostKeyChecking no
UserKnownHostsFile=/dev/null
`

// WriteKey writes the private key.
func WriteKey(privateKey string) error {
	if privateKey == "" {
		return nil
	}

	home := "/root"

	if currentUser, err := user.Current(); err == nil {
		home = currentUser.HomeDir
	}

	sshpath := filepath.Join(
		home,
		".ssh")

	if err := os.MkdirAll(sshpath, 0700); err != nil {
		return err
	}

	confpath := filepath.Join(
		sshpath,
		"config")

	if err := ioutil.WriteFile(
		confpath,
		[]byte(configFile),
		0700,
	); err != nil {
		return err
	}

	privpath := filepath.Join(
		sshpath,
		"id_rsa",
	)

	if err := ioutil.WriteFile(
		privpath,
		[]byte(privateKey),
		0600,
	); err != nil {
		return err
	}

	return nil
}

// WriteNetrc writes the netrc file.
func WriteNetrc(machine, login, password string) error {
	if machine == "" {
		return nil
	}

	netrcContent := fmt.Sprintf(
		netrcFile,
		machine,
		login,
		password,
	)

	home := "/root"

	if currentUser, err := user.Current(); err == nil {
		home = currentUser.HomeDir
	}

	netpath := filepath.Join(
		home,
		".netrc",
	)

	return ioutil.WriteFile(
		netpath,
		[]byte(netrcContent),
		0600,
	)
}
