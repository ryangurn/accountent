package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"ryangurnick.com/accountant/data"
	"ryangurnick.com/accountant/logic"
)

var db *gorm.DB
var err error

func main() {
	db, err := data.OpenConnection()
	if err != nil {
		fmt.Println(err)
	}

	// Migrate the schema
	db.AutoMigrate(&logic.Setting{}, &logic.Address{}, &logic.Contact{}, &logic.Client{})

	app := &cli.App{
		EnableBashCompletion: true,
		Suggest:              true,
		Name:                 "accountant",
		Usage:                "A CLI for managing your business",
		Commands: []*cli.Command{
			{
				Name:        "addresses",
				Aliases:     []string{"a"},
				Usage:       "Manage your addresses",
				Subcommands: logic.AddressSubcommands,
			},
			{
				Name:        "contacts",
				Aliases:     []string{"cont"},
				Usage:       "Manage your contacts",
				Subcommands: logic.ContactSubcommands,
			},
			{
				Name:        "clients",
				Aliases:     []string{"cl"},
				Usage:       "Manage your clients",
				Subcommands: logic.ClientSubcommands,
			},
			{
				Name:        "settings",
				Aliases:     []string{"s"},
				Usage:       "Manage your settings",
				Subcommands: logic.SettingSubcommands,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
