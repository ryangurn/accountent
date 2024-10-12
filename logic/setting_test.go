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
