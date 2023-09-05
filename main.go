package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "git-action plugin"
	app.Usage = "git-action plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "actions",
			Usage:  "actions to execute",
			EnvVar: "PLUGIN_ACTIONS",
		},

		cli.StringFlag{
			Name:   "commit.author.name",
			Usage:  "git author name",
			EnvVar: "PLUGIN_AUTHOR_NAME,DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "git author email",
			EnvVar: "PLUGIN_AUTHOR_EMAIL,DRONE_COMMIT_AUTHOR_EMAIL",
		},

		cli.StringFlag{
			Name:   "netrc.machine",
			Usage:  "netrc machine",
			EnvVar: "PLUGIN_NETRC_MACHINE,DRONE_NETRC_MACHINE",
		},
		cli.StringFlag{
			Name:   "netrc.username",
			Usage:  "netrc username",
			EnvVar: "PLUGIN_NETRC_USERNAME,DRONE_NETRC_USERNAME",
		},
		cli.StringFlag{
			Name:   "netrc.password",
			Usage:  "netrc password",
			EnvVar: "PLUGIN_NETRC_PASSWORD,DRONE_NETRC_PASSWORD",
		},
		cli.StringFlag{
			Name:   "ssh-key",
			Usage:  "private ssh key",
			EnvVar: "PLUGIN_SSH_KEY",
		},

		cli.StringFlag{
			Name:   "remote",
			Usage:  "url of the repo",
			EnvVar: "PLUGIN_REMOTE",
		},
		cli.StringFlag{
			Name:   "branch",
			Usage:  "name of branch",
			EnvVar: "PLUGIN_BRANCH",
			Value:  "master",
		},

		cli.StringFlag{
			Name:   "path",
			Usage:  "path to git repo",
			EnvVar: "PLUGIN_PATH",
		},

		cli.StringFlag{
			Name:   "message",
			Usage:  "commit message",
			EnvVar: "PLUGIN_MESSAGE",
		},

		cli.StringFlag{
			Name:   "tag",
			Usage:  "tag to create",
			EnvVar: "PLUGIN_TAG,DRONE_TAG",
		},

		cli.StringFlag{
			Name:   "tag-message",
			Usage:  "tag message",
			EnvVar: "PLUGIN_TAG_MESSAGE",
		},

		cli.BoolFlag{
			Name:   "force",
			Usage:  "force push to remote",
			EnvVar: "PLUGIN_FORCE",
		},
		cli.BoolFlag{
			Name:   "followtags",
			Usage:  "push to remote with tags",
			EnvVar: "PLUGIN_FOLLOWTAGS",
		},
		cli.BoolFlag{
			Name:   "skip-verify",
			Usage:  "skip ssl verification",
			EnvVar: "PLUGIN_SKIP_VERIFY",
		},
		cli.BoolFlag{
			Name:   "empty-commit",
			Usage:  "allow empty commits",
			EnvVar: "PLUGIN_EMPTY_COMMIT",
		},
		cli.BoolFlag{
			Name:   "no-verify",
			Usage:  "bypasses commit hooks",
			EnvVar: "PLUGIN_NO_VERIFY",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Netrc: Netrc{
			Login:    c.String("netrc.username"),
			Machine:  c.String("netrc.machine"),
			Password: c.String("netrc.password"),
		},
		Commit: Commit{
			Author: Author{
				Name:  c.String("commit.author.name"),
				Email: c.String("commit.author.email"),
			},
		},
		Config: Config{
			Actions:     c.StringSlice("actions"),
			Key:         c.String("ssh-key"),
			Remote:      c.String("remote"),
			Branch:      c.String("branch"),
			Path:        c.String("path"),
			Message:     c.String("message"),
			Tag:         c.String("tag"),
			TagMessage:  c.String("tag-message"),
			Force:       c.Bool("force"),
			FollowTags:  c.Bool("followtags"),
			SkipVerify:  c.Bool("skip-verify"),
			EmptyCommit: c.Bool("empty-commit"),
			NoVerify:    c.Bool("no-verify"),
		},
	}

	return plugin.Exec()
}
