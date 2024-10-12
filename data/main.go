package data

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"os"
)

func OpenConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.data"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}

func CreateTable(c *cli.Context, header table.Row, arr []table.Row) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(c.App.Writer)
	t.AppendHeader(header)
	t.SetStyle(table.StyleLight)

	for _, value := range arr {
		t.AppendRow(value)
	}
	return t
}

func CreateKeyValueTable(c *cli.Context, arr []table.Row) table.Writer {
	t := CreateTable(c, table.Row{"Key", "Value"}, arr)
	return t
}

func CaptureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out), err
}
