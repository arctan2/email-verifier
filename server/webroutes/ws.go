package webroutes

import (
	"database/sql"
	"email_verify/respond"
	"email_verify/socket"
	"email_verify/verifier"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func listenEvents(ws socket.Socket, fileId int64, db *sql.DB) {
	ws.On("get-verifier-details", func(_ []byte) {
		v := verifier.VerifierManager.Get(fileId)
		if v == nil {
			ws.EmitErr("get-verifier-details-res", "verifier not found.").Close()
			return
		}

		socket.EmitWs(ws, "get-verifier-details-res", v.VerifierData)
	})

	ws.On("create-verifier", func(b []byte) {
		var data struct {
			EmailCount int `json:"emailCount"`
			BatchSize int `json:"batchSize"`
			RetryCount int `json:"retryCount"`
			DelayMs int `json:"delayMs"`
			Proxies []string `json:"proxies"`
		}

		if err := json.Unmarshal(b, &data); err != nil {
			ws.EmitErr("create-verifier-res", err.Error()).Close()
			return
		}

		v := verifier.NewVerifier(
			fileId,
			data.EmailCount,
			data.BatchSize,
			data.RetryCount,
			data.DelayMs,
			data.Proxies,
			db,
			ws,
		)

		verifier.VerifierManager.Add(fileId, v)

		socket.EmitWs(ws, "create-verifier-res", respond.SUCCESS)
	})

	ws.On("remove-verifier", func(b []byte) {
		verifier.VerifierManager.Remove(fileId)
	})

	ws.On("run-verifier", func(_ []byte) {
		if v := verifier.VerifierManager.Get(fileId); v != nil {
			v.Run()
		} else {
			ws.EmitErr("run-verifier-res", "verifier not found.").Close()
		}
	})
}

func (m *WebRoutesHandler) verificationWsConn(w http.ResponseWriter, r *http.Request) {
	fileId, err := parseInt64PathValue("fileId", r)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	ws, err := socket.NewWebSocket(w, r)
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	fmt.Println("socket connected", fileId)

	if v := verifier.VerifierManager.Get(fileId); v != nil {
		ws.Emit("status", v.State)
		v.SetWs(ws)
	} else {
		ws.Emit("status", verifier.NOT_CREATED)
	}

	listenEvents(ws, fileId, m.db)

	ws.Listen()
	ws.Close()
}
