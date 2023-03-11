package db

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found")

func OpenDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()
	postgresURL, ok := os.LookupEnv("BOOKSWAP_DB_URL")
	require.True(t, ok)
	db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
	require.Nil(t, err)
	require.NotNil(t, db)
	return db, cleanUpDB(db)
}

func cleanUpDB(gdb *gorm.DB) func() {
	name := "cleanUpDB"
	type record struct {
		table string
		id    string
	}
	var records []record
	gdb.Callback().Create().After("gorm:create").Register(name, func(d *gorm.DB) {
		table := d.Statement.Schema.Table
		model := reflect.ValueOf(d.Statement.Model)
		id := reflect.Indirect(model).FieldByName("ID").String()
		records = append(records, record{table: table, id: id})
	})

	return func() {
		defer gdb.Callback().Create().Remove(name)
		tx := gdb.Begin()
		for _, r := range records {
			tx.Table(r.table).Where("id = ?", r.id).Delete("")
		}
	}
}
