package webroutes

import (
	"context"
	"email_verify/respond"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

func (m *WebRoutesHandler) insertFileDetails(fname string) (error, int64, string, string) {
	var fileId int64 = -1
	fileName := ""

	p := strings.Split(fname, ".")

	ext := p[len(p)-1]

	switch ext {
	case "txt":
	case "csv":
	default:
		return errors.New("Unsupported file extension."), fileId, fileName, ext
	}

	query := `call sp_insert_file(?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := m.db.PrepareContext(ctx, query)

	if err != nil {
		return err, fileId, fileName, ext
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, "admin", fname)

	if err = row.Err(); row.Err() != nil {
		return err, fileId, fileName, ext
	}

	if err := row.Scan(&fileId, &fileName); err != nil {
		return err, fileId, fileName, ext
	}

	return nil, fileId, fileName, ext
}

func (m *WebRoutesHandler) parseTextFileEmails(file multipart.File, fileId int64) ([]string, error) {
	buf, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(buf), "\n")

	return lines, nil
}

func (m *WebRoutesHandler) parseCSVFileEmails(file multipart.File, fileId int64) ([]string, error) {
	reader := csv.NewReader(file)

	reader.ReuseRecord = true
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("File is empty.")
		}
		return nil, err
	}

	re, err := regexp.Compile("(?i)email")
	if err != nil {
		return nil, err
	}

	idx := -1

	for i, head := range header {
		if re.MatchString(head) {
			idx = i
			break
		}
	}

	lines := []string{}

	if idx == -1 {
		idx = 0
		lines = append(lines, header[0])
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		lines = append(lines, record[idx])
	}

	return lines, nil
}

func (m *WebRoutesHandler) insertEmailsToDB(file multipart.File, fileId int64, ext string) (int64, error) {
	var lines []string

	switch ext {
	case "txt":
		l, err := m.parseTextFileEmails(file, fileId)
		if err != nil {
			return 0, err
		}
		lines = l
	case "csv":
		l, err := m.parseCSVFileEmails(file, fileId)
		if err != nil {
			return 0, err
		}
		lines = l
	}

	records := ""
	// literally have no idea why I'm supposed to add the header
	// the load data infile query is ignoring first line even though I
	// am not doing IGNORE 1 LINES
	records += "file_id,email_id"

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		records += fmt.Sprintf("\n"+`"%d","%s"`, fileId, strings.ReplaceAll(line, `"`, `""`))
	}

	reader := strings.NewReader(records)

	handlerID := "upload_csv_data_" + strconv.FormatInt(fileId, 10)

	mysql.RegisterReaderHandler(handlerID, func() io.Reader {
		return reader
	})

	query := fmt.Sprintf(`LOAD DATA LOCAL INFILE 'Reader::%s'
		INTO TABLE emails
		FIELDS TERMINATED BY ',' 
		ENCLOSED BY '"'
		LINES TERMINATED BY '\n';`, handlerID)

	res, err := m.db.Exec(query)
	
	if err != nil {
		return 0, err
	}

	a, err := res.RowsAffected()

	return a, err
}

func (m *WebRoutesHandler) uploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}
	defer file.Close()

	err, fileId, fileName, ext := m.insertFileDetails(fileHeader.Filename)

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	linesCount, err := m.insertEmailsToDB(file, fileId, ext)

	if err != nil {
		if e := m.deleteFile(fileId); e != nil {
			fmt.Println(e.Error())
		}
		respond.RespondErrMsg(w, err.Error())
		return
	}

	res := struct {
		respond.ResponseStruct
		FileName   string `json:"fileName"`
		EmailCount int64    `json:"emailCount"`
		Id         int64  `json:"id"`
	}{
		ResponseStruct: respond.SUCCESS,
		FileName:       fileName,
		EmailCount:     linesCount,
		Id:             fileId,
	}

	json.NewEncoder(w).Encode(&res)
}

func (m *WebRoutesHandler) deleteFile(id int64) error {
	query := `call sp_delete_file_by_id(?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := m.db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, "admin", id)

	return err
}
