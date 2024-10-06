package logic

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"ryangurnick.com/accountant/data"
)

var ClientSubcommands = []*cli.Command{
	{
		Name:      "create",
		Aliases:   []string{"c"},
		Usage:     "Create a new client",
		Action:    CreateClient,
		ArgsUsage: "<name>",
	},
	{
		Name:      "read",
		Aliases:   []string{"r"},
		Usage:     "Read a client",
		Action:    ReadClient,
		ArgsUsage: "<id>",
	},
	{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List all client",
		Action:  ListClients,
	},
	{
		Name:      "update",
		Aliases:   []string{"u"},
		Usage:     "Update a client",
		Action:    UpdateClient,
		ArgsUsage: "<id> <name>",
	},
	{
		Name:      "delete",
		Aliases:   []string{"d"},
		Usage:     "Delete a client",
		Action:    DeleteClient,
		ArgsUsage: "<id>",
	},
	{
		Name:      "link-address",
		Aliases:   []string{"la"},
		Usage:     "Associate an address to a client",
		Action:    AddAddressToClient,
		ArgsUsage: "<client_id> <address_id>",
	},
	{
		Name:      "unlink-address",
		Aliases:   []string{"ua"},
		Usage:     "Remove an address from a client",
		Action:    RemoveAddressFromClient,
		ArgsUsage: "<client_id> <address_id>",
	},
}

type Client struct {
	gorm.Model
	Name      string
	Addresses []Address `gorm:"many2many:client_addresses;"`
	Contacts  []Contact `gorm:"many2many:client_contacts;"`
}

func (c Client) ToString() string {
	return c.Name
}

func CreateClient(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	// map args
	var name = c.Args().Get(0)

	// create client
	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var client = Client{
		Name: name,
	}

	t := data.CreateTable(c, table.Row{"Key", "Value"}, []table.Row{
		{"Name", name},
	})
	t.AppendFooter(table.Row{"Created", ""})
	t.Render()

	db.Create(&client)
	return nil
}

func ReadClient(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var client Client
	tx := db.Preload("Addresses").First(&client, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("address not found")
	}

	t := data.CreateKeyValueTable(c, []table.Row{
		{"ID", client.ID},
		{"Name", client.Name},
	})

	for _, address := range client.Addresses {
		t.AppendRow([]interface{}{fmt.Sprintf("Address %d", address.ID), address.ToString()})
	}
	t.Render()

	return nil
}

func ListClients(c *cli.Context) error {
	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var clients []Client
	db.Find(&clients)

	if len(clients) == 0 {
		return fmt.Errorf("no client found")
	} else {
		var items []table.Row
		for _, client := range clients {
			items = append(items, table.Row{client.ID, client.Name})
		}
		t := data.CreateTable(c, table.Row{"ID", "Name"}, items)

		t.Render()
	}

	return nil
}

func UpdateClient(c *cli.Context) error {
	if c.NArg() != 2 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var client Client
	tx := db.First(&client, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("client not found")
	}

	t := data.CreateTable(c, table.Row{"Key", "Existing", "New"}, []table.Row{
		{"Name", client.Name, c.Args().Get(1)},
	})

	client.Name = c.Args().Get(1)

	tx = db.Save(&client)
	if tx.RowsAffected == 1 {
		t.AppendFooter(table.Row{"Saved", "", ""})
		t.Render()
		return nil
	} else {
		t.AppendFooter(table.Row{"Unable to save", "", ""})
		t.Render()
	}
	return nil
}

func DeleteClient(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var client Client
	tx := db.First(&client, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("client not found")
	}

	t := data.CreateKeyValueTable(c, []table.Row{
		{"ID", client.ID},
		{"Name", client.Name},
	})

	tx = db.Delete(&client)
	if tx.RowsAffected == 1 {
		t.AppendFooter(table.Row{"Deleted", ""})
		t.Render()
		return nil
	} else {
		t.AppendFooter(table.Row{"Unable to delete", ""})
		t.Render()
	}
	return nil
}

func AddAddressToClient(c *cli.Context) error {
	if c.NArg() != 2 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var address Address
	tx := db.First(&address, c.Args().Get(1))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("address not found")
	}

	var client Client
	tx = db.First(&client, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("contact not found")
	}
	err = db.Model(&client).Association("Addresses").Append([]Address{address})
	if err != nil {
		return err
	}

	return nil
}

func RemoveAddressFromClient(c *cli.Context) error {
	if c.NArg() != 2 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var address Address
	tx := db.First(&address, c.Args().Get(1))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("address not found")
	}

	var client Client
	tx = db.First(&client, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("client not found")
	}

	err = db.Model(&client).Association("Addresses").Delete(&address)
	if err != nil {
		return err
	}

	return nil
}
