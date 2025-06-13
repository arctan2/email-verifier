package common

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"themoneybro/respond"
	"themoneybro/schema"
	"time"
)

func UploadRecords(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		respond.RespondErrMsg(w, "Invalid body.")
		return
	}

	query := `CALL sp_insert_records_json(?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, "Prateek", string(body))

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	respond.RespondSuccess(w)
}

func GetAllAccounts(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `CALL sp_get_all_accounts(?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, "Prateek")

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}
	defer rows.Close() 

	var accounts = []schema.Account{}
	for rows.Next() {
		var acc schema.Account
		if err := rows.Scan(&acc.AccountId, &acc.BankName, &acc.Balance); err != nil {
			respond.RespondErrMsg(w, err.Error())
			return 
		}
		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	res := struct {
		respond.ResponseStruct
		Accounts []schema.Account `json:"accounts"`
	} { respond.SUCCESS, accounts }

	json.NewEncoder(w).Encode(&res)
}

func GetAllKeywords(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	query := `select keyword from keywords`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}
	defer rows.Close() 

	keywords := []string{}
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			respond.RespondErrMsg(w, err.Error())
			return 
		}
		keywords = append(keywords, keyword)
	}

	if err := rows.Err(); err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	res := struct {
		respond.ResponseStruct
		Keywords []string `json:"keywords"`
	} { respond.SUCCESS, keywords }

	json.NewEncoder(w).Encode(&res)
}
