package main

import (
	"github.com/codegangsta/cli"
	c "github.com/dlapiduz/govcode.org/common"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "govcode"
	app.Usage = "Tools for govcode"
	app.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "run database migrations",
			Action: func(context *cli.Context) {
				c.DB.AutoMigrate(c.Organization{})
				c.DB.AutoMigrate(c.Repository{})
				c.DB.AutoMigrate(c.Commit{})
				c.DB.AutoMigrate(c.User{})
				c.DB.AutoMigrate(c.Pull{})
				c.DB.AutoMigrate(c.CommitOrgStats{})
				c.DB.AutoMigrate(c.RepoStat{})
			},
		},
		{
			Name:  "import",
			Usage: "run github import",
			Action: func(context *cli.Context) {
				runImport()
			},
		},
		{
			Name:  "stats",
			Usage: "generate stats",
			Action: func(context *cli.Context) {
				generateStats()
			},
		},
	}

	app.Run(os.Args)
}
