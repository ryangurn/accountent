package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"os"
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
	db.AutoMigrate(&logic.Address{}, &logic.Contact{}, &logic.Client{})

	app := &cli.App{
		EnableBashCompletion: true,
		Suggest:              true,
		Name:                 "accountant",
		Usage:                "A CLI for managing your business",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "addresses",
				Aliases: []string{"a"},
				Usage:   "Manage your addresses",
				Subcommands: []*cli.Command{
					{
						Name:      "create",
						Aliases:   []string{"c"},
						Usage:     "Create a new address",
						Action:    logic.CreateAddress,
						ArgsUsage: "<number> <street> <unit> <city> <state> <zip>",
					},
					{
						Name:      "read",
						Aliases:   []string{"r"},
						Usage:     "Read an address",
						Action:    logic.ReadAddress,
						ArgsUsage: "<id>",
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "List all addresses",
						Action:  logic.ListAddresses,
					},
					{
						Name:      "update",
						Aliases:   []string{"u"},
						Usage:     "Update an address",
						Action:    logic.UpdateAddress,
						ArgsUsage: "<id> <number> <street> <unit> <city> <state> <zip>",
					},
					{
						Name:      "delete",
						Aliases:   []string{"d"},
						Usage:     "Delete an address",
						Action:    logic.DeleteAddress,
						ArgsUsage: "<id>",
					},
				},
			},
			{
				Name:    "contacts",
				Aliases: []string{"cont"},
				Usage:   "Manage your contacts",
				Subcommands: []*cli.Command{
					{
						Name:      "create",
						Aliases:   []string{"c"},
						Usage:     "Create a new contact",
						Action:    logic.CreateContact,
						ArgsUsage: "<first_name> <last_name> <email> <phone>",
					},
					{
						Name:      "read",
						Aliases:   []string{"r"},
						Usage:     "Read a contact",
						Action:    logic.ReadContact,
						ArgsUsage: "<id>",
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "List all contacts",
						Action:  logic.ListContacts,
					},
					{
						Name:      "update",
						Aliases:   []string{"u"},
						Usage:     "Update a contact",
						Action:    logic.UpdateContact,
						ArgsUsage: "<id> <first_name> <last_name> <email> <phone>",
					},
					{
						Name:      "delete",
						Aliases:   []string{"d"},
						Usage:     "Delete a contact",
						Action:    logic.DeleteContact,
						ArgsUsage: "<id>",
					},
					{
						Name:      "link-address",
						Aliases:   []string{"la"},
						Usage:     "Associate an address to a contact",
						Action:    logic.AddAddressToContact,
						ArgsUsage: "<contact_id> <address_id>",
					},
					{
						Name:      "unlink-address",
						Aliases:   []string{"ua"},
						Usage:     "Remove an address from a contact",
						Action:    logic.RemoveAddressFromContact,
						ArgsUsage: "<contact_id> <address_id>",
					},
				},
			},
			{
				Name:    "clients",
				Aliases: []string{"cl"},
				Usage:   "Manage your clients",
				Subcommands: []*cli.Command{
					{
						Name:      "create",
						Aliases:   []string{"c"},
						Usage:     "Create a new client",
						Action:    logic.CreateClient,
						ArgsUsage: "<name>",
					},
					{
						Name:      "read",
						Aliases:   []string{"r"},
						Usage:     "Read a client",
						Action:    logic.ReadClient,
						ArgsUsage: "<id>",
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "List all client",
						Action:  logic.ListClients,
					},
					{
						Name:      "update",
						Aliases:   []string{"u"},
						Usage:     "Update a client",
						Action:    logic.UpdateClient,
						ArgsUsage: "<id> <name>",
					},
					{
						Name:      "delete",
						Aliases:   []string{"d"},
						Usage:     "Delete a client",
						Action:    logic.DeleteClient,
						ArgsUsage: "<id>",
					},
					{
						Name:      "link-address",
						Aliases:   []string{"la"},
						Usage:     "Associate an address to a client",
						Action:    logic.AddAddressToClient,
						ArgsUsage: "<client_id> <address_id>",
					},
					{
						Name:      "unlink-address",
						Aliases:   []string{"ua"},
						Usage:     "Remove an address from a client",
						Action:    logic.RemoveAddressFromClient,
						ArgsUsage: "<client_id> <address_id>",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
