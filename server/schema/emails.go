package schema

import "database/sql"

type EmailDetails struct {
	FileId int64 `json:"fileId"`
	EmailId string `json:"emailId"`
	IsValidSyntax bool `json:"isValidSyntax"`
	IsReachable bool `json:"isReachable"`
	IsDeliverable bool `json:"isDeliverable"`
	IsHostExists bool `json:"isHostExists"`
	HasMxRecords bool `json:"hasMxRecords"`
	IsDisposable bool `json:"isDisposable"`
	IsCatchAll bool `json:"isCatchAll"`
	IsInboxFull bool `json:"isInboxFull"`
	ErrorMsg sql.NullString `json:"errorMsg"`
}

