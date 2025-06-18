package webroutes

import (
	"context"
	"email_verify/db"
	"email_verify/respond"
	"email_verify/schema"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (m *WebRoutesHandler) getAllFiles(w http.ResponseWriter, r *http.Request) {
	query := `select id, file_name from files`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := m.db.PrepareContext(ctx, query)

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

	files := []schema.File{}

	for rows.Next() {
		var file schema.File

		if err := rows.Scan(&file.Id, &file.FileName); err != nil {
			respond.RespondErrMsg(w, err.Error())
			return
		}

		files = append(files, file)
	}

	res := struct {
		respond.ResponseStruct
		AllFiles []schema.File `json:"allFiles"`
	}{
		ResponseStruct: respond.SUCCESS,
		AllFiles:       files,
	}

	json.NewEncoder(w).Encode(&res)
}

func (m *WebRoutesHandler) getFileDetails(w http.ResponseWriter, r *http.Request) {
	fileId, err := parseInt64PathValue("fileId", r)
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	totalEmailCount, err := db.GetTotalEmailCount(m.db, fileId)
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	toVerifyCount, err := db.GetToVerifyCount(m.db, fileId)
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	res := struct {
		respond.ResponseStruct
		EmailsCount   int64 `json:"emailsCount"`
		ToVerifyCount int64 `json:"toVerifyCount"`
	}{
		ResponseStruct: respond.SUCCESS,
		EmailsCount:    totalEmailCount,
		ToVerifyCount:  toVerifyCount,
	}

	json.NewEncoder(w).Encode(&res)
}

func (m *WebRoutesHandler) getEmailDetailsList(w http.ResponseWriter, r *http.Request) {
	fileId, err := parseInt64PathValue("fileId", r)
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}
	from, err := parseInt64QueryValue("from", r)
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}
	limit, err := parseInt64QueryValue("limit", r)
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	query := `call sp_get_email_details_from_limit(?, ?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := m.db.PrepareContext(ctx, query)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, fileId, from, limit)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer rows.Close()

	details := []schema.EmailDetails{}

	for rows.Next() {
		var detail schema.EmailDetails

		if err := rows.Scan(
			&detail.FileId,
			&detail.EmailId,
			&detail.IsValidSyntax,
			&detail.Reachable,
			&detail.IsDeliverable,
			&detail.IsHostExists,
			&detail.HasMxRecords,
			&detail.IsDisposable,
			&detail.IsCatchAll,
			&detail.IsInboxFull,
			&detail.ErrorMsg,
		); err != nil {
			respond.RespondErrMsg(w, err.Error())
			return
		}

		details = append(details, detail)
	}

	res := struct {
		respond.ResponseStruct
		EmailDetailsList []schema.EmailDetails `json:"emailDetailsList"`
	}{
		ResponseStruct:   respond.SUCCESS,
		EmailDetailsList: details,
	}

	json.NewEncoder(w).Encode(&res)
}

func translateDetailsFieldToDBField(field string) string {
	switch field {
	case "isValidSyntax":
		return "is_valid_syntax"
	case "reachable":
		return "reachable" 
	case "isDeliverable":
		return "is_deliverable"
	case "isHostExists":
		return "is_host_exists"
	case "hasMxRecords":
		return "has_mx_records"
	case "isDisposable":
		return "is_disposable"
	case "isCatchAll":
		return "is_catch_all"
	case "isInboxFull":
		return "is_inbox_full"
	}
	return ""
}

func (m *WebRoutesHandler) filterEmails(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FileId       int64          `json:"fileId"`
		FilterFields map[string]any `json:"filterFields"`
		From         uint            `json:"from"`
		Limit        uint            `json:"limit"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	if body.Limit == 0 {
		body.Limit = 500
	}

	if body.From == 0 {
		body.From = 0
	}

	where := []string{}
	args := []any{}


	for k, v := range body.FilterFields {
		where = append(where, fmt.Sprintf("(%s = ?)", translateDetailsFieldToDBField(k)))
		args = append(args, v)
	}

	wh := strings.Join(where, "and")

	query := fmt.Sprintf(`select * from emails where ((file_id = %d) and %s) limit %d, %d`,
		body.FileId, wh, body.From, body.Limit,
	)

	queryCount := fmt.Sprintf(`select count(*) from emails where ((file_id = %d) and %s)`,
		body.FileId, wh,
	)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	rows, err := m.db.QueryContext(ctx, query, args...)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer rows.Close()

	details := []schema.EmailDetails{}

	for rows.Next() {
		var detail schema.EmailDetails

		if err := rows.Scan(
			&detail.FileId,
			&detail.EmailId,
			&detail.IsValidSyntax,
			&detail.Reachable,
			&detail.IsDeliverable,
			&detail.IsHostExists,
			&detail.HasMxRecords,
			&detail.IsDisposable,
			&detail.IsCatchAll,
			&detail.IsInboxFull,
			&detail.ErrorMsg,
		); err != nil {
			respond.RespondErrMsg(w, err.Error())
			return
		}

		details = append(details, detail)
	}

	row := m.db.QueryRowContext(ctx, queryCount, args...)

	if row.Err() != nil {
		respond.RespondErrMsg(w, row.Err().Error())
		return
	}

	var totalCount int64 = 0

	if err := row.Scan(&totalCount); err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	res := struct {
		respond.ResponseStruct
		EmailDetailsList []schema.EmailDetails `json:"emailDetailsList"`
		TotalEmailCount int64 `json:"totalEmailCount"`
	}{
		ResponseStruct:   respond.SUCCESS,
		EmailDetailsList: details,
		TotalEmailCount: totalCount,
	}

	json.NewEncoder(w).Encode(&res)

}
