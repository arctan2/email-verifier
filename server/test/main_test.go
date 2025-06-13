package test

import (
	"fmt"
	"log"
	"testing"
	"themoneybro/dbconn"
)

func ClearDataFromDatabase() (error) {
	db := dbconn.DbTestConn()

	query := `call sp_truncate_tables_for_test`

	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	fmt.Println("Truncated all tables successfully.")
	return nil
}

func TestMain(m *testing.M) {
	if err := ClearDataFromDatabase(); err != nil {
		log.Fatal(err.Error())
	}
	m.Run()
}
