package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"themoneybro/dbconn"
	"themoneybro/respond"
	"themoneybro/schema"
	"themoneybro/webroutes"
)

func webRoutesInitDefaults(t *testing.T) (*http.ServeMux, *sql.DB, schema.Account) {
	if err := ClearDataFromDatabase(); err != nil {
		log.Fatal(err.Error())
	}

	db := dbconn.DbTestConn()

	mux := webroutes.NewWebRoutesMux(db)

	acc := schema.Account{ AccountId: "1234", BankName: "some", Balance: 0 }

	if err := insertAccount(db, acc); err != nil {
		t.Fatalf("Insert Account failed: %s", err.Error())
	}

	return mux, db, acc
}

func TestInsertTransact(t *testing.T) {
	mux, db, acc := webRoutesInitDefaults(t)

	var records []schema.TransactOfAccountId

	var sum float64 = 0

	transacts := []float64{10000, 100, -200, 300, 40, 35, -600, -200, -2000, 80}

	for _, v := range transacts {
		var transferType string

		if v > 0 {
			transferType = "credit"
		} else {
			transferType = "debit"
		}

		record := schema.TransactOfAccountId{
			Amount: math.Abs(v),
			TransferType: transferType,
			FromTo: "",
			Date: "17/05/2025",
			Time: "21:28",
			Keywords: "test_keyword",
			RawSMS: "raw raw raw sms\nwith stuff",
		}
		records = append(records, record)

		j, err := json.Marshal(&record)
		
		if err != nil {
			t.Fatalf(err.Error())
		}

		req := httptest.NewRequest("POST", "/" + acc.AccountId + "/insert-transaction-of-acc", bytes.NewReader(j))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if err, _ = isResponseValid(rr); err != nil {
			t.Fatalf(err.Error())
		}

		sum += v
	}

	if err := cmpBalanceWithDb(db, acc.AccountId, sum); err != nil {
        t.Fatalf(err.Error())
	}
}

func TestInsertRandomShitTranscat(t *testing.T) {
	mux, db, acc := webRoutesInitDefaults(t)

	var records []schema.TransactOfAccountId

	var sum float64 = 0

	transacts := []float64{}

	for i := 0; i < 100; i++ {
		transacts = append(transacts, randFloat64(1, 10000))
	}

	for _, v := range transacts {
		var transferType string

		if v > 0 {
			transferType = "credit"
		} else {
			transferType = "debit"
		}

		record := schema.TransactOfAccountId{
			Amount: math.Abs(v),
			TransferType: transferType,
			FromTo: "",
			Date: "17/05/2025",
			Time: "21:28",
			Keywords: "test_keyword",
			RawSMS: "raw raw raw sms\nwith stuff",
		}
		records = append(records, record)

		j, err := json.Marshal(&record)
		
		if err != nil {
			t.Fatalf(err.Error())
		}

		req := httptest.NewRequest("POST", "/" + acc.AccountId + "/insert-transaction-of-acc", bytes.NewReader(j))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if err, _ = isResponseValid(rr); err != nil {
			t.Fatalf(err.Error())
		}

		sum += roundFloat(v, 2)
	}

	if err := cmpBalanceWithDb(db, acc.AccountId, sum); err != nil {
        t.Fatalf(err.Error())
	}
}

func TestUpdateTransact(t *testing.T) {
	mux, db, acc := webRoutesInitDefaults(t)

	var records []schema.TransactRecord

	var sum float64 = 0

	amounts := []float64{10000, 100, -200, 300, 40, 35, -600, -200, -2000, 80}

	for _, v := range amounts {
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

	var resBody struct {
		respond.ResponseStruct
		Transactions []schema.TransactOfAccountId `json:"transactions"`
	}

	req = httptest.NewRequest("GET", "/" + acc.AccountId + "/get-transactions", nil)
    rr = httptest.NewRecorder()

    mux.ServeHTTP(rr, req)

	err, body := isResponseValid(rr)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if err = json.Unmarshal(body, &resBody); err != nil {
		t.Fatalf(err.Error())
	}

	transacts := resBody.Transactions
	amountUpdates := []float64{570, 100, 20, -450, 3456, -55, 600, -300, 20, -40}
	dateUpdates := []string{}
	keywords := []string{}
	transferTypes := []string{}

	for i := 0; i < 10; i++ {
		dateUpdates = append(dateUpdates, "28/03/2025")
	}

	for i := 0; i < 10; i++ {
		keywords = append(keywords, "keyword1\nkeyword2")
	}

	sum = 0

	for i, tr := range transacts {
		tr.Keywords = keywords[i]
		
		var transferType string

		if amountUpdates[i] > 0 {
			transferType = "credit"
		} else {
			transferType = "debit"
		}

		transferTypes = append(transferTypes, transferType)

		tr.Amount = math.Abs(amountUpdates[i])
		tr.TransferType = transferType
		tr.Date = dateUpdates[i]

		j, err := json.Marshal(&tr) 

		if err != nil {
			t.Fatalf(err.Error())
		}

		req = httptest.NewRequest("POST", "/" + acc.AccountId + "/update-transaction-of-acc", bytes.NewReader(j))
		rr = httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if err, _ = isResponseValid(rr); err != nil {
			t.Fatalf(err.Error())
		}

		sum += amountUpdates[i]
	}

	req = httptest.NewRequest("GET", "/" + acc.AccountId + "/get-transactions", nil)
    rr = httptest.NewRecorder()

    mux.ServeHTTP(rr, req)

	err, body = isResponseValid(rr)

	if err != nil {
		t.Fatalf(err.Error())
	}

	var resBody2 struct {
		respond.ResponseStruct
		Transactions []schema.TransactOfAccountId `json:"transactions"`
	}

	if err = json.Unmarshal(body, &resBody2); err != nil {
		t.Fatalf(err.Error())
	}

	transacts2 := resBody2.Transactions

	for i := range transacts {
		t2 := transacts2[i]

		if math.Abs(amountUpdates[i]) != t2.Amount {
			t.Fatalf("Expected %f got %f", amountUpdates[i], t2.Amount)
		}
		if transferTypes[i] != t2.TransferType {
			t.Fatalf("Expected %s got %s", transferTypes[i], t2.TransferType)
		}
		if dateUpdates[i] != t2.Date {
			t.Fatalf("Expected %s got %s", dateUpdates[i], t2.Date)
		}
		if keywords[i] != t2.Keywords {
			t.Fatalf("Expected %s got %s", keywords[i], t2.Keywords)
		}
	}

	if err := cmpBalanceWithDb(db, acc.AccountId, sum); err != nil {
        t.Fatalf(err.Error())
	}
}

func TestDeleteTransact(t *testing.T) {
	mux, db, acc := webRoutesInitDefaults(t)

	var records []schema.TransactOfAccountId

	var sum float64 = 0

	transacts := []float64{10000, 100, -200, 300, 40, 35, -600, -200, -2000, 80}

	for _, v := range transacts {
		var transferType string

		if v > 0 {
			transferType = "credit"
		} else {
			transferType = "debit"
		}

		record := schema.TransactOfAccountId{
			Amount: math.Abs(v),
			TransferType: transferType,
			FromTo: "",
			Date: "17/05/2025",
			Time: "21:28",
			Keywords: "test_keyword",
			RawSMS: "raw raw raw sms\nwith stuff",
		}
		records = append(records, record)

		j, err := json.Marshal(&record)
		
		if err != nil {
			t.Fatalf(err.Error())
		}

		req := httptest.NewRequest("POST", "/" + acc.AccountId + "/insert-transaction-of-acc", bytes.NewReader(j))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if err, _ = isResponseValid(rr); err != nil {
			t.Fatalf(err.Error())
		}

		sum += v
	}

	if err := cmpBalanceWithDb(db, acc.AccountId, sum); err != nil {
        t.Fatalf(err.Error())
	}

	for i := range transacts {
		req := httptest.NewRequest("DELETE", "/delete-transaction-record?transact_id=" + strconv.Itoa(i + 1), nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if err, _ := isResponseValid(rr); err != nil {
			t.Fatalf(err.Error())
		}
	}

	if err := cmpBalanceWithDb(db, acc.AccountId, 0); err != nil {
        t.Fatalf(err.Error())
	}
}
