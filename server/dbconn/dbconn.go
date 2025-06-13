package dbconn

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func DbConn(isProdMode bool) *sql.DB {
	if !isProdMode {
		fmt.Print("Please confirm that you're using prod db.")
		fmt.Scanln()
	}

	var db *sql.DB

    cfg := mysql.NewConfig()
    cfg.User = "root"
    cfg.Passwd = "@arctan2"
    cfg.Net = "tcp"
    cfg.Addr = "127.0.0.1:3306"
    cfg.DBName = "email_verifier"

    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
		fmt.Println("Unable to make DB connection")
		return nil
    }

    pingErr := db.Ping()
    if pingErr != nil {
		fmt.Println("Ping failed")
		return nil
    }

	fmt.Println("connected to", cfg.DBName)

	return db
}

func DbTestConn() *sql.DB {
	var db *sql.DB

    cfg := mysql.NewConfig()
    cfg.User = "root"
    cfg.Passwd = "@arctan2"
    cfg.Net = "tcp"
    cfg.Addr = "127.0.0.1:3306"
    cfg.DBName = "email_verifier"

    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
		fmt.Println("Unable to make DB connection")
		return nil
    }

    pingErr := db.Ping()
    if pingErr != nil {
		fmt.Println("Ping failed")
		return nil
    }

	fmt.Println("connected to", cfg.DBName)

	return db
}
