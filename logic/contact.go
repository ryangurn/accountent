package logic

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"ryangurnick.com/accountant/data"
)

var ContactSubcommands = []*cli.Command{
	{
		Name:      "create",
		Aliases:   []string{"c"},
		Usage:     "Create a new contact",
		Action:    CreateContact,
		ArgsUsage: "<first_name> <last_name> <email> <phone>",
	},
	{
		Name:      "read",
		Aliases:   []string{"r"},
		Usage:     "Read a contact",
		Action:    ReadContact,
		ArgsUsage: "<id>",
	},
	{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List all contacts",
		Action:  ListContacts,
	},
	{
		Name:      "update",
		Aliases:   []string{"u"},
		Usage:     "Update a contact",
		Action:    UpdateContact,
		ArgsUsage: "<id> <first_name> <last_name> <email> <phone>",
	},
	{
		Name:      "delete",
		Aliases:   []string{"d"},
		Usage:     "Delete a contact",
		Action:    DeleteContact,
		ArgsUsage: "<id>",
	},
	{
		Name:      "link-address",
		Aliases:   []string{"la"},
		Usage:     "Associate an address to a contact",
		Action:    AddAddressToContact,
		ArgsUsage: "<contact_id> <address_id>",
	},
	{
		Name:      "unlink-address",
		Aliases:   []string{"ua"},
		Usage:     "Remove an address from a contact",
		Action:    RemoveAddressFromContact,
		ArgsUsage: "<contact_id> <address_id>",
	},
}

type Contact struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Addresses []Address `gorm:"many2many:contact_addresses;"`
}

func (c Contact) ToString() string {
	return fmt.Sprintf("%s %s\n%s\n%s", c.FirstName, c.LastName, c.Email, c.Phone)
}

func CreateContact(c *cli.Context) error {
	if c.NArg() != 4 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	// map args
	var firstName = c.Args().Get(0)
	var lastName = c.Args().Get(1)
	var email = c.Args().Get(2)
	var phone = c.Args().Get(3)

	// create contact
	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var contact = Contact{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	t := data.CreateTable(c, table.Row{"First Name", "Last Name", "Email", "Phone"}, []table.Row{{contact.FirstName, contact.LastName, contact.Email, contact.Phone}})
	t.AppendFooter(table.Row{"Created", ""})
	t.Render()

	db.Create(&contact)
	return nil
}

func ReadContact(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var contact Contact
	tx := db.Preload("Addresses").First(&contact, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("contact not found")
	}

	t := data.CreateKeyValueTable(c, []table.Row{
		{"First Name", contact.FirstName},
		{"Last Name", contact.LastName},
		{"Email", contact.Email},
		{"Phone", contact.Phone},
	})
	for _, address := range contact.Addresses {
		t.AppendRow([]interface{}{fmt.Sprintf("Address %d", address.ID), address.ToString()})
	}
	t.Render()
	return nil
}

func ListContacts(c *cli.Context) error {
	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var contacts []Contact
	tx := db.Preload("Addresses").Find(&contacts)
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("no contacts found")
	}

	var items []table.Row
	for _, contact := range contacts {
		items = append(items, table.Row{contact.ID, contact.FirstName, contact.LastName, contact.Email, contact.Phone, len(contact.Addresses)})
	}
	t := data.CreateTable(c, table.Row{"ID", "First Name", "Last Name", "Email", "Phone", "Address"}, items)
	t.Render()

	return nil
}

func UpdateContact(c *cli.Context) error {
	if c.NArg() != 5 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var contact Contact
	tx := db.First(&contact, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("contact not found")
	}

	t := data.CreateTable(c, table.Row{"Key", "Existing", "New"}, []table.Row{
		{"First Name", contact.FirstName, c.Args().Get(1)},
		{"Last Name", contact.LastName, c.Args().Get(2)},
		{"Email", contact.Email, c.Args().Get(3)},
		{"Phone", contact.Phone, c.Args().Get(4)},
	})

	contact.FirstName = c.Args().Get(1)
	contact.LastName = c.Args().Get(2)
	contact.Email = c.Args().Get(3)
	contact.Phone = c.Args().Get(4)
	tx = db.Save(&contact)
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

func DeleteContact(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("missing required arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var contact Contact
	tx := db.First(&contact, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("contact not found")
	}

	t := data.CreateKeyValueTable(c, []table.Row{
		{"First Name", contact.FirstName},
		{"Last Name", contact.LastName},
		{"Email", contact.Email},
		{"Phone", contact.Phone},
	})

	tx = db.Delete(&contact)
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

func AddAddressToContact(c *cli.Context) error {
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

	var contact Contact
	tx = db.First(&contact, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("contact not found")
	}
	err = db.Model(&contact).Association("Addresses").Append([]Address{address})
	if err != nil {
		return err
	}

	return nil
}

func RemoveAddressFromContact(c *cli.Context) error {
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

	var contact Contact
	tx = db.First(&contact, c.Args().Get(0))
	if tx.RowsAffected == 0 || tx.Error != nil {
		return fmt.Errorf("contact not found")
	}

	err = db.Model(&contact).Association("Addresses").Delete(&address)
	if err != nil {
		return err
	}

	return nil
}
