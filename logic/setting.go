package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"ryangurnick.com/accountant/data"
)

var SettingSubcommands = []*cli.Command{
	{
		Name:    "setup",
		Aliases: []string{"s"},
		Usage:   "Setup default settings",
		Action:  SetupSettings,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force setup, even if settings exist (deletes existing settings)",
			},
		},
	},
	{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List all settings",
		Action:  ListSettings,
	},
	{
		Name:      "update",
		Aliases:   []string{"u"},
		Usage:     "Update a setting",
		Action:    UpdateSetting,
		ArgsUsage: "<namespace> <key> <newValue>",
	},
	{
		Name:    "populate-missing",
		Aliases: []string{"pm"},
		Usage:   "Add missing settings, will not replace any existing values",
		Action:  AddMissingSettings,
	},
	{
		Name:      "export",
		Aliases:   []string{"e"},
		Usage:     "Export settings to a file",
		Action:    ExportSettings,
		ArgsUsage: "<filename>",
	},
}

type Setting struct {
	gorm.Model
	Namespace string
	Key       string
	Value     string
}

var defaults = []Setting{
	// business basic info
	{Namespace: "business", Key: "name", Value: "Business Name"},
	{Namespace: "business", Key: "phone", Value: "(555) 555-5555"},
	{Namespace: "business", Key: "email", Value: "example@email.com"},
	{Namespace: "business", Key: "address_number", Value: "1234"},
	{Namespace: "business", Key: "address_street", Value: "Main St"},
	{Namespace: "business", Key: "address_unit", Value: "Suite 101"},
	{Namespace: "business", Key: "address_city", Value: "Portland"},
	{Namespace: "business", Key: "address_state", Value: "Oregon"},
	{Namespace: "business", Key: "address_zip", Value: "97217"},
	{Namespace: "business", Key: "address_country", Value: "United States"},
	{Namespace: "business", Key: "timezone", Value: "America/Los_Angeles"},

	// tax and financial information
	{Namespace: "financial", Key: "currency", Value: "USD"},
	{Namespace: "financial", Key: "year_end_month", Value: "12"},
	{Namespace: "financial", Key: "year_end_day", Value: "31"},
	{Namespace: "financial", Key: "standard_rate", Value: "50.00"},
}

func (s Setting) ToString() string {
	return strings.ToLower(s.Namespace) + "." + strings.ToLower(s.Key)
}

func SetupSettings(c *cli.Context) error {
	// open connection
	force := c.Bool("force")
	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	// delete existing if force
	if force { // not sure
		db.Exec("DELETE FROM settings")
		db.Exec("DELETE FROM SQLITE_SEQUENCE WHERE name='settings';") // reset autoincrement
	}

	// check if matching count
	var settings []Setting
	db.Find(&settings)
	if len(settings) > 0 && !force {
		return fmt.Errorf("there appear to be some settings already, consider running list")
	}

	// create settings
	for _, setting := range defaults {
		db.Create(&setting)
	}

	return nil
}

func AddMissingSettings(c *cli.Context) error {
	items := []table.Row{}
	for _, setting := range defaults {
		// open connection
		db, err := data.OpenConnection()
		if err != nil {
			return err
		}

		// check if setting exists
		var count int64
		db.Table("settings").Where("namespace = ? AND key = ?", setting.Namespace, setting.Key).Count(&count)

		if count == 0 {
			tx := db.Create(&setting)
			if tx.RowsAffected == 1 {
				items = append(items, table.Row{setting.Namespace, setting.Key, "created"})
			} else {
				items = append(items, table.Row{setting.Namespace, setting.Key, "error"})
			}
		} else {
			items = append(items, table.Row{setting.Namespace, setting.Key, "exists"})
		}
	}

	t := data.CreateTable(c, table.Row{"Namespace", "Key", "Status"}, items)
	t.Render()

	return nil
}

func ListSettings(c *cli.Context) error {
	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var settings []Setting
	db.Find(&settings)

	if len(settings) == 0 {
		return fmt.Errorf("no settings found, consider running setup")
	} else {
		var items []table.Row
		for _, setting := range settings {
			items = append(items, table.Row{setting.ID, setting.Namespace, setting.Key, setting.Value})
		}
		t := data.CreateTable(c, table.Row{"ID", "Namespace", "Key", "Value"}, items)
		t.Render()
	}

	return nil
}

func UpdateSetting(c *cli.Context) error {
	if c.NArg() != 3 {
		return fmt.Errorf("incorrect number of arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	namespace := c.Args().Get(0)
	key := c.Args().Get(1)
	newValue := c.Args().Get(2)

	var setting Setting
	db.Where("namespace = ? AND key = ?", namespace, key).First(&setting)
	if setting.ID == 0 {
		return fmt.Errorf("setting not found")
	}

	t := data.CreateTable(c, table.Row{"Key", "Existing", "New"}, []table.Row{
		{"Namespace", setting.Namespace, namespace},
		{"Key", setting.Key, key},
		{"Value", setting.Value, newValue},
	})

	setting.Value = newValue
	tx := db.Save(&setting)
	if tx.RowsAffected == 1 {
		t.AppendFooter(table.Row{"Saved", "", ""})
		t.Render()
	} else {
		t.AppendFooter(table.Row{"Unable to save", "", ""})
		t.Render()
	}

	return nil
}

func ExportSettings(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("incorrect number of arguments: " + c.Command.ArgsUsage)
	}

	db, err := data.OpenConnection()
	if err != nil {
		return err
	}

	var settings []Setting
	db.Find(&settings)

	bytee, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if os.Stat(c.Args().Get(0)); os.IsExist(err) {
		return fmt.Errorf("file already exists: " + c.Args().Get(0))
	}

	os.WriteFile(c.Args().Get(0), bytee, 0644)

	return nil
}
