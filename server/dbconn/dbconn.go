package dbconn

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func getUserPwd() (string, string, error) {
	f, err := os.Open("./auth.txt")

	if err != nil {
		return "", "", err
	}

	buf, err := io.ReadAll(f)

	if err != nil {
		return "", "", err
	}

	lines := strings.Split(string(buf), "\n")

	if len(lines) < 2 {
		return "", "", errors.New("Too few lines in the auth file.")
	}

	return lines[0], lines[1], nil
}

func DbConn(isProdMode bool) *sql.DB {
	if !isProdMode {
		fmt.Print("Please confirm that you're using prod db.")
		fmt.Scanln()
	}

	var db *sql.DB

	user, pwd, err := getUserPwd()

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

    cfg := mysql.NewConfig()
    cfg.User = user
	cfg.Passwd = pwd
    cfg.Net = "tcp"
    cfg.Addr = "127.0.0.1:3306"
    cfg.DBName = "email_verifier"

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

	user, pwd, err := getUserPwd()

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

    cfg := mysql.NewConfig()
    cfg.User = user
    cfg.Passwd = pwd
    cfg.Net = "tcp"
    cfg.Addr = "127.0.0.1:3306"
    cfg.DBName = "email_verifier"

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
