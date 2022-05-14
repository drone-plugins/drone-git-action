package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/drone-plugins/drone-git-action/repo"
)

type (
	Netrc struct {
		Machine  string
		Login    string
		Password string
	}

	Commit struct {
		Author Author
	}

	Author struct {
		Name  string
		Email string
	}

	Config struct {
		Actions     []string
		Key         string
		Remote      string
		Branch      string
		Path        string
		Message     string
		Tag         string
		TagMessage  string
		Force       bool
		FollowTags  bool
		SkipVerify  bool
		EmptyCommit bool
		NoVerify    bool
	}

	Plugin struct {
		Netrc  Netrc
		Commit Commit
		Config Config
	}
)

func (p *Plugin) Exec() error {
	if err := p.HandlePath(); err != nil {
		return err
	}

	if err := p.WriteConfig(); err != nil {
		return err
	}

	if err := p.WriteKey(); err != nil {
		return err
	}

	if err := p.WriteNetrc(); err != nil {
		return err
	}

	for _, action := range p.Config.Actions {
		switch action {
		case "clone":
			if err := p.InitRepo(); err != nil {
				return err
			}

			if err := p.AddRemote(); err != nil {
				return err
			}

			if err := p.FetchSource(); err != nil {
				return err
			}

			if err := p.CheckoutHead(); err != nil {
				return err
			}
		case "commit":
			if err := p.HandleCommit(); err != nil {
				return err
			}
		case "push":
			if err := p.HandlePush(); err != nil {
				return err
			}
		case "tag":
			if err := p.HandleTag(); err != nil {
				return err
			}
		case "push-tag":
			if err := p.HandlePushTag(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unknown action %s", action)
		}
	}

	return nil
}

// HandlePath changes to a different directory if required
func (p Plugin) HandlePath() error {
	if p.Config.Path != "" {
		if err := os.MkdirAll(p.Config.Path, os.ModePerm); err != nil {
			return err
		}

		if err := os.Chdir(p.Config.Path); err != nil {
			return err
		}
	}

	return nil
}

// WriteConfig writes all required configurations.
func (p Plugin) WriteConfig() error {
	if err := repo.GlobalName(p.Commit.Author.Name).Run(); err != nil {
		return err
	}

	if err := repo.GlobalUser(p.Commit.Author.Email).Run(); err != nil {
		return err
	}

	if p.Config.SkipVerify {
		if err := repo.SkipVerify().Run(); err != nil {
			return err
		}
	}

	return nil
}

// WriteKey writes the private SSH key.
func (p Plugin) WriteKey() error {
	return repo.WriteKey(
		p.Config.Key,
	)
}

// WriteNetrc writes the netrc config.
func (p Plugin) WriteNetrc() error {
	return repo.WriteNetrc(
		p.Netrc.Machine,
		p.Netrc.Login,
		p.Netrc.Password,
	)
}

// InitRepo initializes the repository.
func (p Plugin) InitRepo() error {
	if isDirEmpty(filepath.Join(p.Config.Path, ".git")) {
		return execute(exec.Command(
			"git",
			"init",
		))
	}

	return nil
}

// AddRemote adds a remote to repository.
func (p Plugin) AddRemote() error {
	if p.Config.Remote != "" {
		if err := execute(repo.RemoteAdd("origin", p.Config.Remote)); err != nil {
			return err
		}
	}

	return nil
}

// FetchSource fetches the source from remote.
func (p Plugin) FetchSource() error {
	return execute(exec.Command(
		"git",
		"fetch",
		"origin",
		fmt.Sprintf("+%s:", p.Config.Branch),
	))
}

// CheckoutHead handles branch checkout.
func (p Plugin) CheckoutHead() error {
	return execute(exec.Command(
		"git",
		"checkout",
		"-qf",
		p.Config.Branch,
	))
}

// HandleCommit commits changes locally.
func (p Plugin) HandleCommit() error {
	if err := execute(repo.Add()); err != nil {
		return err
	}

	if err := execute(repo.TestCleanTree()); err != nil {
		if err := execute(repo.ForceCommit(p.Config.Message, p.Config.NoVerify)); err != nil {
			return err
		}
	} else {
		if p.Config.EmptyCommit {
			if err := execute(repo.EmptyCommit(p.Config.Message, p.Config.NoVerify)); err != nil {
				return err
			}
		}
	}

	return nil
}

// HandlePush pushs changes to remote.
func (p Plugin) HandlePush() error {
	return execute(repo.RemotePushNamedBranch(
		"origin",
		p.Config.Branch,
		p.Config.Branch,
		p.Config.Force,
		p.Config.FollowTags,
	))
}

// HandleTag creates a tag on the current commit.
func (p Plugin) HandleTag() error {
	return execute(repo.CreateTag(
		p.Config.Tag,
		p.Config.TagMessage,
	))
}

// HandlePushTag pushes a tag to remote.
func (p Plugin) HandlePushTag() error {
	return execute(repo.RemotePushTag(
		"origin",
		p.Config.Tag,
	))
}
