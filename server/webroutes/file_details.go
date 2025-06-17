package webroutes

import (
	"context"
	"email_verify/respond"
	"email_verify/schema"
	"encoding/json"
	"net/http"
	"time"
)

func (m *WebRoutesHandler) getFileListStatsLimit(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		respond.RespondErrMsg(w, "invalid userId")
		return
	}
	from, err := parseInt64QueryValue("from", r)
	if err != nil {
		from = 0
	}
	limit, err := parseInt64QueryValue("limit", r)
	if err != nil {
		limit = -1
	}

	query := `call sp_file_list_stats_limit(?, ?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := m.db.PrepareContext(ctx, query)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userId, from, limit)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer rows.Close()

	list := []schema.FileStats{}

	for rows.Next() {
		var d schema.FileStats

		if err := rows.Scan(
			&d.FileId,
			&d.FileName,
			&d.TotalEmails,
			&d.InvalidSyntax,
			&d.Reachable,
			&d.Unknown,
			&d.Deliverable,
			&d.CatchAll,
			&d.Disposable,
			&d.InboxFull,
			&d.HostExists,
			&d.Errored,
		); err != nil {
			respond.RespondErrMsg(w, err.Error())
			return
		}

		list = append(list, d)
	}

	res := struct {
		respond.ResponseStruct
		StatsList []schema.FileStats `json:"statsList"`
	}{
		ResponseStruct:   respond.SUCCESS,
		StatsList: list,
	}

	json.NewEncoder(w).Encode(&res)
}

func (m *WebRoutesHandler) getFileStats(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		respond.RespondErrMsg(w, "invalid userId")
		return
	}
	fileId, err := parseInt64QueryValue("fileId", r)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	query := `call sp_get_file_stats(?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := m.db.PrepareContext(ctx, query)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userId, fileId)

	if row == nil {
		respond.RespondErrMsg(w, "Record not found.")
		return
	}

	var d schema.FileStats

	if err := row.Scan(
		&d.FileName,
		&d.TotalEmails,
		&d.InvalidSyntax,
		&d.Reachable,
		&d.Unknown,
		&d.Deliverable,
		&d.CatchAll,
		&d.Disposable,
		&d.InboxFull,
		&d.HostExists,
		&d.Errored,
	); err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	res := struct {
		respond.ResponseStruct
		FileStats schema.FileStats `json:"fileStats"`
	}{
		ResponseStruct:   respond.SUCCESS,
		FileStats: d,
	}

	json.NewEncoder(w).Encode(&res)
}
