package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

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

func (e *EmailDetails) ToCSV() string {
	em := ""

	if !e.ErrorMsg.Valid {
		em = "NULL"
	} else {
		em = e.ErrorMsg.String
	}

	return fmt.Sprintf(
		`"%d","%s","%d","%d","%d","%d","%d","%d","%d","%d","%s"`,
		e.FileId,
		strings.ReplaceAll(e.EmailId, `"`, `""`),
		boolToInt(e.IsValidSyntax),
		boolToInt(e.IsReachable),
		boolToInt(e.IsDeliverable),
		boolToInt(e.IsHostExists),
		boolToInt(e.HasMxRecords),
		boolToInt(e.IsDisposable),
		boolToInt(e.IsCatchAll),
		boolToInt(e.IsInboxFull),
		strings.ReplaceAll(em, `"`, `""`),
	)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (e *EmailDetails) ToCSVLn() string {
	em := ""

	if !e.ErrorMsg.Valid {
		em = "NULL"
	} else {
		em = e.ErrorMsg.String
	}

	return fmt.Sprintf(
		`"%d","%s","%d","%d","%d","%d","%d","%d","%d","%d","%s"
`,
		e.FileId,
		strings.ReplaceAll(e.EmailId, `"`, `""`),
		boolToInt(e.IsValidSyntax),
		boolToInt(e.IsReachable),
		boolToInt(e.IsDeliverable),
		boolToInt(e.IsHostExists),
		boolToInt(e.HasMxRecords),
		boolToInt(e.IsDisposable),
		boolToInt(e.IsCatchAll),
		boolToInt(e.IsInboxFull),
		strings.ReplaceAll(em, `"`, `""`),
	)
}
