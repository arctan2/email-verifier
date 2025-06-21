package webroutes

import (
	"context"
	"email_verify/respond"
	"email_verify/schema"
	"encoding/json"
	"net/http"
	"time"
)

func (m *WebRoutesHandler) getProxyList(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	if userId == "" {
		respond.RespondErrMsg(w, "userId invalid")
		return
	}

	q := `call sp_get_proxy_list(?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	rows, err := m.db.QueryContext(ctx, q, userId)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}
	defer rows.Close()

	list := []schema.ProxyDetails{}

	for rows.Next() {
		var p schema.ProxyDetails

		if err := rows.Scan(
			&p.Id,
			&p.Proto,
			&p.Host,
			&p.Port,
			&p.Name,
			&p.Password,
			&p.IsInUse,
		); err != nil {
			respond.RespondErrMsg(w, err.Error())
			return
		}

		list = append(list, p)
	}

	res := struct {
		respond.ResponseStruct
		ProxyList []schema.ProxyDetails `json:"proxyList"`
	} {
		respond.SUCCESS,
		list,
	}

	json.NewEncoder(w).Encode(&res)
}

func (m *WebRoutesHandler) insertProxy(w http.ResponseWriter, r *http.Request) {
	var body schema.ProxyDetails

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	q := `call sp_insert_proxy(?, ?, ?, ?, ?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err = m.db.ExecContext(
		ctx, q,
		body.UserId,
		body.Proto,
		body.Host,
		body.Port,
		body.Name,
		body.Password,
	)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(&respond.SUCCESS)
}

func (m *WebRoutesHandler) deleteProxy(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	if userId == "" {
		respond.RespondErrMsg(w, "userId invalid")
		return
	}

	proxyId, err := parseInt64PathValue("proxyId", r)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	q := `call sp_delete_proxy(?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err = m.db.ExecContext(ctx, q, userId, proxyId)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(&respond.SUCCESS)
}

func (m *WebRoutesHandler) updateProxy(w http.ResponseWriter, r *http.Request) {
	proxyId, err := parseInt64PathValue("proxyId", r)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	var body schema.ProxyDetails

	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	if body.UserId == "" {
		respond.RespondErrMsg(w, "Invalid userId")
		return
	}

	q := `call sp_update_proxy(?, ?, ?, ?, ?, ?, ?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err = m.db.ExecContext(
		ctx, q,
		body.UserId,
		proxyId,
		body.Proto,
		body.Host,
		body.Port,
		body.Name,
		body.Password,
		false,
	)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(&respond.SUCCESS)
}
