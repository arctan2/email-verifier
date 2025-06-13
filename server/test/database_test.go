package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"themoneybro/dbconn"
	"themoneybro/mobileroutes"
	"themoneybro/respond"
	"themoneybro/schema"
)

func insertAccount(db *sql.DB, acc schema.Account) error {
	q := "insert into bank_accounts (account_id, bank_name, user_name, balance) values (?, ?, ?, ?)"
	_, err := db.Exec(q, acc.AccountId, acc.BankName, "Prateek", acc.Balance)
	return err
}

func roundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func randFloat64(min float64, max float64) float64 {
	return min + (max - min) * rand.Float64()
}

func randBool() bool {
    return rand.IntN(2) == 1
}

func cmpBalanceWithDb(db *sql.DB, accountId string, expected float64) error {
	q := `select balance from bank_accounts where account_id = ?`

	row := db.QueryRow(q, accountId)

	if row == nil {
		return errors.New("Row is nil")
	}

	var balance float64

	err := row.Scan(&balance)

	if err != nil {
		return err
	}

	if balance != expected {
		return errors.New(fmt.Sprintf("Expected '%f' but got '%f'", expected, balance))
	}

	return nil
}

func isResponseValid(rr *httptest.ResponseRecorder) (error, []byte) {
	if rr.Code != http.StatusOK {
		return errors.New("expected status 200 OK, got " + strconv.Itoa(rr.Code)), nil
	}
	
	body, err := io.ReadAll(rr.Body)

	if err != nil {
		return err, nil
	}

	var res respond.ResponseStruct

	if err = json.Unmarshal(body, &res); err != nil {
		return err, nil
	}

	if res.Err {
		return errors.New(res.Msg), nil
	}

	return nil, body
}

func TestUploadRecords(t *testing.T) {
	if err := ClearDataFromDatabase(); err != nil {
		log.Fatal(err.Error())
	}

	db := dbconn.DbTestConn()

	mux := mobileroutes.NewMobileRoutesMux(db)

	acc := schema.Account{ AccountId: "1234", BankName: "some", Balance: 0 }

	if err := insertAccount(db, acc); err != nil {
		t.Fatalf("Insert Account failed: %s", err.Error())
	}

	var records []schema.TransactRecord

	var sum float64 = 0

	transacts := []float64{10000, 100, -200, 300, 40, 35, -600, -200, -2000, 80}

	for _, v := range transacts {
		var transferType string

		if v > 0 {
			transferType = "credit"
		} else {
			transferType = "debit"
		}

		record := schema.TransactRecord{
			AccountId: acc.AccountId,
			Amount: math.Abs(v),
			TransferType: transferType,
			FromTo: "",
			Date: "17/05/2025",
			Time: "21:28",
			Keywords: "test_keyword",
			RawSMS: "raw raw raw sms\nwith stuff",
		}
		records = append(records, record)

		sum += v
	}

	j, err := json.Marshal(records)

	if err != nil {
		t.Fatalf("json marshal failed: %s", err.Error())
	}

	req := httptest.NewRequest("POST", "/upload-records", bytes.NewReader(j))
	rr := httptest.NewRecorder()

    mux.ServeHTTP(rr, req)

	if err, _ = isResponseValid(rr); err != nil {
		t.Fatalf(err.Error())
	}

	if err := cmpBalanceWithDb(db, acc.AccountId, sum); err != nil {
        t.Fatalf(err.Error())
	}
}
