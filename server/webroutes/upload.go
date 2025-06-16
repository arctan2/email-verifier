package webroutes

import (
	"bytes"
	"context"
	"email_verify/respond"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

func (m *WebRoutesHandler) insertFileDetails(fname string) (error, int64, string) {
	var fileId int64 = -1
	fileName := ""

	p := strings.Split(fname, ".")

	ext := p[len(p) - 1]

	switch(ext) {
	case "txt":
	default:
		return errors.New("Unsupported file extension."), fileId, fileName
	}

	query := `call sp_insert_file(?, ?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := m.db.PrepareContext(ctx, query)

	if err != nil {
		return err, fileId, fileName
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, "admin", fname)

	if err = row.Err(); row.Err() != nil {
		return err, fileId, fileName
	}

	if err := row.Scan(&fileId, &fileName); err != nil {
		return err, fileId, fileName
	}

	return nil, fileId, fileName
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

func (m *WebRoutesHandler) insertEmailsToDB(file multipart.File, fileId int64) (error) {
	buf, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	lines := bytes.Split(buf, []byte("\n"))

	idBytes := make([]byte, len(buf))
	// literally have no idea why I'm supposed to add the header
	// the load data infile query is ignoring first line even though I
	// am not doing IGNORE 1 LINES
	idBytes = append(idBytes, []byte("file_id,email_id\n")...)

	for i := range lines {
		if len(lines[i]) == 0 {
			continue
		}

		b := []byte(fmt.Sprintf("\"%d\",\"", fileId))
		lines[i] = append(b, bytes.ReplaceAll(lines[i], []byte(`"`), []byte(`""`))...)
		lines[i] = append(lines[i], []byte("\"\n")...)
		idBytes = append(idBytes, lines[i]...)
	}

	reader := bytes.NewReader(idBytes)

	handlerID := "upload_csv_data_" + strconv.FormatInt(fileId, 10)

	mysql.RegisterReaderHandler(handlerID, func() io.Reader {
		return reader
	})

	query := fmt.Sprintf(`LOAD DATA LOCAL INFILE 'Reader::%s'
		INTO TABLE emails
		FIELDS TERMINATED BY ',' 
		ENCLOSED BY '"'
		LINES TERMINATED BY '\n';`, handlerID)

	_, err = m.db.Exec(query)

	return err
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

	err, fileId, fileName := m.insertFileDetails(fileHeader.Filename);

	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	if err = m.insertEmailsToDB(file, fileId); err != nil {
		if e := m.deleteFile(fileId); e != nil {
			fmt.Println(e.Error())
		}
		respond.RespondErrMsg(w, err.Error())
		return
	}

	res := struct {
		respond.ResponseStruct
		FileName string `json:"fileName"`
		Id int64 `json:"id"`
	}{
		ResponseStruct: respond.SUCCESS,
		FileName: fileName,
		Id: fileId,
	}

	json.NewEncoder(w).Encode(&res)
}
