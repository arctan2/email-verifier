package webroutes

import (
	"database/sql"
	"email_verify/respond"
	"net/http"
	"strconv"
)

type WebRoutesHandler struct {
	mux *http.ServeMux
	db  *sql.DB
}

func NewWebRoutesMux(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	m := WebRoutesHandler{mux, db}
	m.setupRoutes()
	return mux
}

func (m *WebRoutesHandler) deleteFileRoute(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		respond.RespondErrMsg(w, "file_id not provided.")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	if err := m.deleteFile(id); err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	respond.RespondSuccess(w)
}

func (m *WebRoutesHandler) setupProxyRoutes() {
	m.mux.HandleFunc("GET /get-proxy-list", m.getProxyList)
	m.mux.HandleFunc("POST /insert-proxy", m.insertProxy)
	m.mux.HandleFunc("PUT /{proxyId}/update-proxy", m.updateProxy)
	m.mux.HandleFunc("PUT /{proxyId}/update-proxy-is-enabled", m.updateProxyIsEnabled)
	m.mux.HandleFunc("DELETE /{proxyId}/delete-proxy", m.deleteProxy)
}

func (m *WebRoutesHandler) setupFileRoutes() {
	m.mux.HandleFunc("GET /get-all-files", m.getAllFiles)
	m.mux.HandleFunc("GET /{fileId}/get-file-details", m.getFileDetails)
	m.mux.HandleFunc("GET /get-file-list-stats", m.getFileListStatsLimit)
	m.mux.HandleFunc("GET /get-file-stats", m.getFileStats)

	m.mux.HandleFunc("POST /upload-file", m.uploadFile)

	m.mux.HandleFunc("DELETE /delete-file", m.deleteFileRoute)
}

func (m *WebRoutesHandler) setupEmailRoutes() {
	m.mux.HandleFunc("GET /{fileId}/get-email-details-list", m.getEmailDetailsList)

	m.mux.HandleFunc("POST /verify-emails", m.verifyEmails)
	m.mux.HandleFunc("POST /filter-emails", m.filterEmails)
}

func (m *WebRoutesHandler) setupRoutes() {
	m.setupFileRoutes()
	m.setupEmailRoutes()
	m.setupProxyRoutes()

	m.mux.HandleFunc("/{fileId}/verification-ws", m.verificationWsConn)
}
