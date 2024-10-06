package logic

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"ryangurnick.com/accountant/data"
)

var AddressSubcommands = []*cli.Command{
	{
		Name:      "create",
		Aliases:   []string{"c"},
		Usage:     "Create a new address",
		Action:    CreateAddress,
		ArgsUsage: "<number> <street> <unit> <city> <state> <zip>",
	},
	{
		Name:      "read",
		Aliases:   []string{"r"},
		Usage:     "Read an address",
		Action:    ReadAddress,
		ArgsUsage: "<id>",
	},
	{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List all addresses",
		Action:  ListAddresses,
	},
	{
		Name:      "update",
		Aliases:   []string{"u"},
		Usage:     "Update an address",
		Action:    UpdateAddress,
		ArgsUsage: "<id> <number> <street> <unit> <city> <state> <zip>",
	},
	{
		Name:      "delete",
		Aliases:   []string{"d"},
		Usage:     "Delete an address",
		Action:    DeleteAddress,
		ArgsUsage: "<id>",
	},
}

type Address struct {
	gorm.Model
	Number   string
	Street   string
	Unit     string
	City     string
	State    string
	Zip      string
	Clients  []*Client  `gorm:"many2many:client_addresses;"`
	Contacts []*Contact `gorm:"many2many:contact_addresses;"`
}

func (a Address) ToString() string {
	return fmt.Sprintf("%s %s %s\n%s, %s %s", a.Number, a.Street, a.Unit, a.City, a.State, a.Zip)
}

func CreateAddress(c *cli.Context) error {
	if c.NArg() != 6 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	// map args
	var number = c.Args().Get(0)
	var street = c.Args().Get(1)
	var unit = c.Args().Get(2)
	var city = c.Args().Get(3)
	var state = c.Args().Get(4)
	var zip = c.Args().Get(5)

	// create address
	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var address = Address{
		Number: number,
		Street: street,
		Unit:   unit,
		City:   city,
		State:  state,
		Zip:    zip,
	}

	t := data.CreateTable(c, table.Row{"Key", "Value"}, []table.Row{
		{"Number", number},
		{"Street", street},
		{"Unit", unit},
		{"City", city},
		{"State", state},
		{"Zip", zip},
	})
	t.AppendFooter(table.Row{"Created", ""})
	t.Render()

	db.Create(&address)
	return nil
}

func ReadAddress(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var address Address
	tx := db.Preload("Clients").Preload("Contacts").First(&address, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("address not found")
	}

	var items = []table.Row{
		{"ID", address.ID},
		{"Number", address.Number},
		{"Street", address.Street},
		{"Unit", address.Unit},
		{"City", address.City},
		{"State", address.State},
		{"Zip", address.Zip},
	}

	t := data.CreateKeyValueTable(c, items)
	// clients
	for _, client := range address.Clients {
		t.AppendRow([]interface{}{fmt.Sprintf("Client #%d", client.ID), client.ToString()})
	}

	// contacts
	for _, contact := range address.Contacts {
		t.AppendRow([]interface{}{fmt.Sprintf("Contact #%d", contact.ID), contact.ToString()})
	}

	t.AppendFooter(table.Row{"Found", ""})
	t.Render()
	return nil
}

func ListAddresses(c *cli.Context) error {
	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var addresses []Address
	db.Find(&addresses)

	if len(addresses) == 0 {
		return fmt.Errorf("no addresses found")
	} else {
		var items []table.Row
		for _, address := range addresses {
			items = append(items, table.Row{address.ID, address.Number, address.Street, address.Unit, address.City, address.State, address.Zip})
		}
		t := data.CreateTable(c, table.Row{"ID", "Number", "Street", "Unit", "City", "State", "Zip"}, items)
		t.Render()
	}

	return nil
}

func UpdateAddress(c *cli.Context) error {
	if c.NArg() != 7 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var address Address
	tx := db.First(&address, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("address not found")
	}

	t := data.CreateTable(c, table.Row{"Key", "Existing", "New"}, []table.Row{
		{"Number", address.Number, c.Args().Get(1)},
		{"Street", address.Street, c.Args().Get(2)},
		{"Unit", address.Unit, c.Args().Get(3)},
		{"City", address.City, c.Args().Get(4)},
		{"State", address.State, c.Args().Get(5)},
		{"Zip", address.Zip, c.Args().Get(6)},
	})

	address.Number = c.Args().Get(1)
	address.Street = c.Args().Get(2)
	address.Unit = c.Args().Get(3)
	address.City = c.Args().Get(4)
	address.State = c.Args().Get(5)
	address.Zip = c.Args().Get(6)

	tx = db.Save(&address)
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

func DeleteAddress(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var address Address
	tx := db.First(&address, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("address not found")
	}

	t := data.CreateKeyValueTable(c, []table.Row{
		{"ID", address.ID},
		{"Number", address.Number},
		{"Street", address.Street},
		{"Unit", address.Unit},
		{"City", address.City},
		{"State", address.State},
		{"Zip", address.Zip},
	})

	tx = db.Delete(&address)
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
