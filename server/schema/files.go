package schema

type File struct {
	Id       int64 `json:"id"`
	FileName string `json:"fileName"`
	UserId string `json:"userId"`
}

type FileStats struct {
	FileId int64 `json:"fileId"`
	FileName string `json:"fileName"`
	TotalEmails int64 `json:"totalEmails"`
	InvalidSyntax int64 `json:"invalidSyntax"`
	Reachable int64 `json:"reachable"`
	Unknown int64 `json:"unknown"`
	Deliverable int64 `json:"deliverable"`
	CatchAll int64 `json:"catchAll"`
	Disposable int64 `json:"disposable"`
	InboxFull int64 `json:"inboxFull"`
	HostExists int64 `json:"hostExists"`
	Errored int64 `json:"errored"`
}
