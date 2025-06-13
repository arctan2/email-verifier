package webroutes

import (
	"context"
	"email_verify/respond"
	"email_verify/schema"
	"encoding/json"
	"net/http"
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

	query := `select count(*) from emails where file_id = ?`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := m.db.PrepareContext(ctx, query)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, fileId)

	if err = row.Err(); row.Err() != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	var count int64 = 0
	if err := row.Scan(&count); err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	res := struct {
		respond.ResponseStruct
		EmailsCount int64 `json:"emailsCount"`
	}{
		ResponseStruct: respond.SUCCESS,
		EmailsCount:    count,
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
			&detail.IsReachable,
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
		ResponseStruct: respond.SUCCESS,
		EmailDetailsList:       details,
	}

	json.NewEncoder(w).Encode(&res)
}
