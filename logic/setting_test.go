package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"ryangurnick.com/accountant/data"
)

func setupSuite(tb testing.TB) func(tb testing.TB) {
	// truncate the table
	db, dbErr := data.OpenConnection()
	if dbErr != nil {
		tb.Fatal(dbErr)
	}

	db.Exec("DELETE FROM settings")
	db.Exec("DELETE FROM SQLITE_SEQUENCE WHERE name='settings';")

	return func(tb testing.TB) {

	}
}

func TestSetting_ToString(t *testing.T) {
	// setup
	setting := Setting{
		Namespace: "test",
		Key:       "test",
		Value:     "Test",
	}

	// test + assert
	assert.Equal(t, "test.test", setting.ToString())
}

func TestSetup(t *testing.T) {
	// setup
	suite := setupSuite(t)
	defer suite(t)

	// test
	app := cli.App{}
	app.Commands = SettingSubcommands
	err := app.Run([]string{"setting", "setup"})

	// assert
	assert.Nil(t, err, "Error should be nil")

	db, dbErr := data.OpenConnection()
	if dbErr != nil {
		t.Fatal(dbErr)
	}

	var settings []Setting
	db.Find(&settings)

	assert.Equal(t, len(defaults), len(settings))
}

func TestSetup_Force(t *testing.T) {
	// setup
	suite := setupSuite(t)
	defer suite(t)

	// test
	app := cli.App{}
	app.Commands = SettingSubcommands
	err := app.Run([]string{"setting", "setup", "--force"})

	// assert
	db, dbErr := data.OpenConnection()
	assert.Nil(t, dbErr, "DB should not return an error")

	var setting []Setting
	db.Find(&setting)

	assert.NotEqual(t, 0, len(setting), "Settings should not be empty")
	assert.Nil(t, err, "Error should be nil")
}

func TestAddMissingSettings(t *testing.T) {
	// setup
	suite := setupSuite(t)
	defer suite(t)

	// test
	app := cli.App{}
	app.Commands = SettingSubcommands
	output, err := data.CaptureOutput(func() error {
		err := app.Run([]string{"setting", "populate-missing"})
		return err
	})

	// assert
	assert.Nil(t, err, "Error should be nil")

	db, dbErr := data.OpenConnection()
	assert.Nil(t, dbErr, "DB should not return an error")

	var setting []Setting
	db.Find(&setting)

	assert.Equal(t, len(defaults), len(setting), "Settings should not be empty")
	assert.Nil(t, err, "Error should be nil")

	assert.NotContains(t, output, "error")
	assert.Contains(t, output, "created")
}

func TestAddMissingSettings_Exists(t *testing.T) {
	// setup
	// not including suite setup, to retain prior data

	// test
	app := cli.App{}
	app.Commands = SettingSubcommands
	output, err := data.CaptureOutput(func() error {
		err := app.Run([]string{"setting", "populate-missing"})
		return err
	})

	assert.Nil(t, err, "Error should be nil")
	db, dbErr := data.OpenConnection()
	assert.Nil(t, dbErr, "DB should not return an error")
	var setting []Setting
	db.Find(&setting)

	assert.Equal(t, len(defaults), len(setting), "Settings should not be empty")
	assert.Nil(t, err, "Error should be nil")
	assert.Contains(t, output, "exists")
	assert.NotContains(t, output, "created")
}

func TestListSettings(t *testing.T) {
	suite := setupSuite(t)
	defer suite(t)

	app := cli.App{}
	app.Commands = SettingSubcommands
	output, err := data.CaptureOutput(func() error {
		err := app.Run([]string{"setting", "list"})
		return err
	})

	assert.NotNil(t, err, "no settings found, consider running setup")
	assert.Emptyf(t, output, "output should ")
}

func TestListSettings_WithData(t *testing.T) {
	suite := setupSuite(t)
	defer suite(t)

	app := cli.App{}
	app.Commands = SettingSubcommands
	output, err := data.CaptureOutput(func() error {
		app.Run([]string{"setting", "setup", "--force"})

		err := app.Run([]string{"setting", "list"})
		return err
	})

	assert.Nil(t, err, "no settings found, consider running setup")
	assert.NotEmpty(t, output, "output should contain settings list")
}

func TestUpdateSetting(t *testing.T) {
	suite := setupSuite(t)
	defer suite(t)

	app := cli.App{}
	app.Commands = SettingSubcommands
	output, err := data.CaptureOutput(func() error {
		app.Run([]string{"setting", "setup", "--force"})

		err := app.Run([]string{"setting", "update", "business", "name", "Test"})
		return err
	})

	assert.Nil(t, err, "no settings found, consider running setup")
	assert.Contains(t, output, "Test")
	assert.Contains(t, output, "SAVED")
	assert.NotContains(t, output, "UNABLE TO SAVE")
}

func TestUpdateSetting_Failed(t *testing.T) {
	suite := setupSuite(t)
	defer suite(t)

	app := cli.App{}
	app.Commands = SettingSubcommands
	output, err := data.CaptureOutput(func() error {
		app.Run([]string{"setting", "setup", "--force"})

		err := app.Run([]string{"setting", "update", "business", "invalid_KEY_123", "Test"})
		return err
	})

	assert.NotNil(t, err, "setting not found")
	assert.Empty(t, output, "output should be empty")
}

func TestExportSettings(t *testing.T) {
	suite := setupSuite(t)
	defer suite(t)

	app := cli.App{}
	app.Commands = SettingSubcommands
	output, err := data.CaptureOutput(func() error {
		app.Run([]string{"setting", "setup", "--force"})

		err := app.Run([]string{"setting", "export", "export.json"})
		return err
	})

	assert.Nil(t, err, "no settings found, consider running setup")
	assert.Empty(t, output, "output is not expected")
}
