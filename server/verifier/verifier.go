package verifier

import (
	"context"
	"database/sql"
	"email_verify/db"
	"email_verify/schema"
	"email_verify/socket"
	"fmt"
	"io"
	"maps"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	// emailverifier "github.com/AfterShip/email-verifier"
)

type ProgressData struct {
	Total int `json:"total"`
	Progress int `json:"progress"`
	Success int `json:"success"`
	Failed int `json:"failed"`
	Retry int `json:"retry"`
	sync.Mutex `json:"-"`
}

const (
	NOT_CREATED = "not created"
	CREATED = "created"
	RUNNING = "running"
	DONE = "done"
)

type VerifierData struct {
	State string `json:"state"`
	EmailCount int `json:"emailCount"`
	BatchSize int `json:"batchSize"`
	RetryCount int `json:"retryCount"`
	DelayMs int `json:"delayMs"`
	Proxies []string `json:"proxies"`
	CurProxyIdx int `json:"curProxyIdx"`

	CompletedBatches map[int][]*ProgressData `json:"completedBatches"`

	CurrentBatch []schema.EmailDetails `json:"-"`
	CurrentBatchNumber int `json:"currentBatchNumber"`
	CurrentBatchSize int `json:"currentBatchSize"`
	CurrentProgressList []*ProgressData `json:"currentProgressList"`
}

type Verifier struct {
	db *sql.DB
	ws socket.Socket
	File schema.File

	VerifierData
}

func NewVerifier(
	fileId int64,
	emailCount,
	batchSize,
	retryCount,
	delayMs int,
	proxies []string,
	db *sql.DB,
	ws socket.Socket,
) *Verifier {
	v := Verifier{}

	v.File.Id = fileId
	v.EmailCount = emailCount
	v.BatchSize = batchSize
	v.RetryCount = retryCount
	v.DelayMs = delayMs
	v.Proxies = proxies
	v.State = CREATED
	v.db = db
	v.ws = ws

	return &v
}

func NewProgressData(total int) *ProgressData {
	d := &ProgressData{}
	d.Total = total
	return d
}

func (v *Verifier) SetWs(ws socket.Socket) {
	v.ws = ws
}

func (v *Verifier) incProgress(success, failed, retry int) {
	idx := len(v.CurrentProgressList) - 1
	p := v.CurrentProgressList[idx]

	p.Lock()
	
	p.Progress++

	p.Success += success
	p.Failed += failed
	p.Retry += retry

	if p.Progress != 0 && p.Progress % 10 == 0 {
		socket.EmitWs(v.ws, "progress", p)
	}

	p.Unlock()
}

func (v *Verifier) Emit(ev string, data string) {
	if v.ws != nil {
		v.ws.Emit(ev, data)
	}
}

func (v *Verifier) updateProxy() {
	if len(v.Proxies) == 0 {
		v.CurProxyIdx = -1
		return
	}

	v.CurProxyIdx++

	if v.CurProxyIdx >= len(v.Proxies) {
		v.CurProxyIdx = 0
	}
}

func (v *Verifier) Run() error {
	v.State = RUNNING
	emails, err := db.GetEmailsForVerification(v.db, v.File.Id)
	if err != nil {
		return err
	}

	batchSize := v.BatchSize
	delay := v.DelayMs
	retryRate := v.RetryCount

	s := batchSize

	if len(emails) < s {
		s = len(emails)
	}

	v.CurrentBatch = make([]schema.EmailDetails, s)
	for i := range v.CurrentBatch {
		v.CurrentBatch[i] = schema.NewEmailDetails()
	}

	v.CompletedBatches = make(map[int][]*ProgressData)
	v.CurrentProgressList = []*ProgressData{}
	v.CurProxyIdx = -1

	v.updateProxy()

	truncateAndConvCsvBatch := func(batchSize int) string {
		csvStr := ""

		for i := 0; i < batchSize; i++ {
			csvStr += v.CurrentBatch[i].ToCSVLn()
			v.CurrentBatch[i] = schema.NewEmailDetails()
		}

		return strings.TrimSuffix(csvStr, "\n")
	}

	i := 0

	socket.EmitWs(v.ws, "get-verifier-details-res", v.VerifierData)

	for ; i+batchSize < len(emails); i += batchSize {
		v.CurrentProgressList = []*ProgressData{NewProgressData(batchSize)}
		v.Emit("batch-start", strconv.Itoa(v.CurrentBatchNumber))

		v.verifyBatch(emails, i, i+batchSize, delay, retryRate)

		v.ws.Emit("update-db-start", "")
		csvStr := truncateAndConvCsvBatch(batchSize)
		if err := v.updateCsvStrToDB(csvStr); err != nil {
			return err
		}
		v.ws.Emit("update-db-done", "")

		v.CompletedBatches[v.CurrentBatchNumber] = make([]*ProgressData, len(v.CurrentProgressList))
		for i := range v.CompletedBatches[v.CurrentBatchNumber] {
			p := v.CurrentProgressList[i]
			n := ProgressData{
				Total: p.Total,
				Progress: p.Progress,
				Success: p.Success,
				Failed: p.Failed,
				Retry: p.Retry,
			}
			v.CompletedBatches[v.CurrentBatchNumber][i] = &n
		}

		v.Emit("batch-delay", "")
		time.Sleep(time.Duration(delay) * time.Millisecond)
		v.CurrentBatchNumber++
		v.updateProxy()
	}

	if i < len(emails) {
		v.CurrentProgressList = []*ProgressData{NewProgressData(len(emails) - i)}
		v.Emit("batch-start", strconv.Itoa(v.CurrentBatchNumber))

		v.verifyBatch(emails, i, len(emails), delay, retryRate)

		v.ws.Emit("update-db-start", "")
		csvStr := truncateAndConvCsvBatch(len(emails) - i)
		if err := v.updateCsvStrToDB(csvStr); err != nil {
			return err
		}
		v.ws.Emit("update-db-done", "")

		v.CompletedBatches[v.CurrentBatchNumber] = make([]*ProgressData, len(v.CurrentProgressList))
		for i := range v.CompletedBatches[v.CurrentBatchNumber] {
			p := v.CurrentProgressList[i]
			n := ProgressData{
				Total: p.Total,
				Progress: p.Progress,
				Success: p.Success,
				Failed: p.Failed,
				Retry: p.Retry,
			}
			v.CompletedBatches[v.CurrentBatchNumber][i] = &n
		}
	}

	v.State = DONE

	socket.EmitWs(v.ws, "get-verifier-details-res", v.VerifierData)
	return nil
}

func (v *Verifier) updateCsvStrToDB(csvStr string) error {
	tmpTableId := "tmp_tbl_" + strconv.FormatInt(v.File.Id, 10) + "_" + strconv.Itoa(v.CurrentBatchNumber)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := v.db.ExecContext(ctx, fmt.Sprintf(`CREATE TEMPORARY TABLE IF NOT EXISTS %s (
		file_id int NOT NULL,
		email_id varchar(320) NOT NULL,
		is_valid_syntax tinyint NOT NULL DEFAULT '0',
		reachable varchar(10) NOT NULL DEFAULT '',
		is_deliverable tinyint NOT NULL DEFAULT '0',
		is_host_exists tinyint NOT NULL DEFAULT '0',
		has_mx_records tinyint NOT NULL DEFAULT '0',
		is_disposable tinyint NOT NULL DEFAULT '0',
		is_catch_all tinyint NOT NULL DEFAULT '0',
		is_inbox_full tinyint NOT NULL DEFAULT '0',
		error_msg text DEFAULT NULL,
		PRIMARY KEY (email_id)
	)`, tmpTableId))

	if err != nil {
		return err
	}

	reader := strings.NewReader(csvStr)

	handlerID := "load_file_to_tmp_tbl_" + strconv.FormatInt(v.File.Id, 10) + "_" + strconv.Itoa(v.CurrentBatchNumber)

	mysql.RegisterReaderHandler(handlerID, func() io.Reader {
		return reader
	})

	_, err = v.db.ExecContext(ctx, fmt.Sprintf(`
	LOAD DATA LOCAL INFILE 'Reader::%s'
	INTO TABLE %s
	FIELDS TERMINATED BY ',' 
	ENCLOSED BY '"' 
	LINES TERMINATED BY '\n'`, handlerID, tmpTableId))

	if err != nil {
		return err
	}

	_, err = v.db.ExecContext(ctx, fmt.Sprintf(`
		UPDATE emails e
		JOIN %s t ON e.file_id = t.file_id and e.email_id = t.email_id
		set
			e.is_valid_syntax = t.is_valid_syntax,
			e.reachable = t.reachable,
			e.is_deliverable = t.is_deliverable,
			e.is_host_exists = t.is_host_exists,
			e.has_mx_records = t.has_mx_records,
			e.is_disposable = t.is_disposable,
			e.is_catch_all = t.is_catch_all,
			e.is_inbox_full = t.is_inbox_full,
			e.error_msg = t.error_msg
	`, tmpTableId))

	if err != nil {
		return err
	}

	_, err = v.db.ExecContext(ctx, fmt.Sprintf(`DROP TEMPORARY TABLE %s`, tmpTableId))

	return err
}

type RetryState struct {
	idxs        map[int]int
	toRetryIdxs map[int]int
	sync.Mutex
}

func (s *RetryState) add(batchIdx, idx int) {
	s.Lock()
	s.idxs[batchIdx] = idx
	s.Unlock()
}

func (s *RetryState) reset() {
	s.toRetryIdxs = make(map[int]int)
	maps.Copy(s.toRetryIdxs, s.idxs)
	s.idxs = make(map[int]int)
}

func (v *Verifier) verifyBatch(emails []string, from, to, delay, retryRate int) {
	retryState := RetryState{make(map[int]int), make(map[int]int), sync.Mutex{}}
	wg := sync.WaitGroup{}
	i := from
	batchIdx := 0
	for i < to {
		email := emails[i]
		wg.Add(1)
		go v.verifyEmail(email, batchIdx, i, &retryState, false, &wg)
		i++
		batchIdx++
	}
	wg.Wait()

	retryState.reset()

	if len(retryState.toRetryIdxs) == 0 {
		socket.EmitWs(v.ws, "after-all-retries", v.CurrentProgressList[len(v.CurrentProgressList) - 1])
		return
	}

	for i := 0; i < retryRate; i++ {
		socket.EmitWs(v.ws, "retry-delay", v.CurrentProgressList[len(v.CurrentProgressList) - 1])
		time.Sleep(time.Duration(delay) * time.Millisecond)
		l := len(retryState.toRetryIdxs)

		p := NewProgressData(l)

		v.CurrentProgressList = append(v.CurrentProgressList, p)
		socket.EmitWs(v.ws, "retry-begin", *p)
		v.retryBatch(emails, i+1 == retryRate, &retryState)
		if len(retryState.idxs) == 0 {
			socket.EmitWs(v.ws, "after-all-retries", *v.CurrentProgressList[len(v.CurrentProgressList) - 1])
			return
		}
		retryState.reset()
	}

	socket.EmitWs(v.ws, "after-all-retries", *v.CurrentProgressList[len(v.CurrentProgressList) - 1])
}

func (v *Verifier) retryBatch(emails []string, isLastRetry bool, retryState *RetryState) {
	wg := sync.WaitGroup{}
	for batchIdx, idx := range retryState.toRetryIdxs {
		email := emails[idx]
		wg.Add(1)
		go v.verifyEmail(email, batchIdx, idx, retryState, isLastRetry, &wg)
	}
	wg.Wait()
}

func (v *Verifier) verifyEmail(email string, batchIdx, idx int, retryState *RetryState, isLastRetry bool, wg *sync.WaitGroup) {
	// var (
	// 	verifier = emailverifier.
	// 		NewVerifier().
	// 		EnableSMTPCheck().
	// 		EnableCatchAllCheck()
	// )

	var (
		verifier = NewTestVerifier().
			EnableSMTPCheck().
			EnableCatchAllCheck()
	)

	if v.CurProxyIdx != -1 {
		verifier = verifier.Proxy(v.Proxies[v.CurProxyIdx])
	}

	checkIsLastRetry := func(e string) {
		if(isLastRetry) {
			v.CurrentBatch[batchIdx].ErrorMsg = sql.NullString{String: e, Valid: true}
		}
	}

	defer wg.Done()

	ret, err := verifier.Verify(email)

	emailDetails := &v.CurrentBatch[batchIdx]

	emailDetails.FileId = v.File.Id
	emailDetails.EmailId = email

	if ret != nil {
		emailDetails.ErrorMsg = sql.NullString{String: "", Valid: true}

		if !ret.Syntax.Valid {
			emailDetails.IsValidSyntax = false
			v.incProgress(1, 1, 0)
			return
		}

		emailDetails.IsValidSyntax = ret.Syntax.Valid
		emailDetails.HasMxRecords = ret.HasMxRecords
	}

	if err != nil {
		e := err.Error()
		retry := 0
		if strings.Contains(e, "no such host") {
			emailDetails.IsHostExists = false
			v.incProgress(1, 1, 0)
			return
		} else if strings.Contains(e, "has timed out") {
			checkIsLastRetry(e)
			retryState.add(batchIdx, idx)
			retry = 1
		} else if strings.Contains(e, "i/o timeout") {
			checkIsLastRetry(e)
			retryState.add(batchIdx, idx)
			retry = 1
		} else if strings.Contains(e, "temporarily unavailable") {
			checkIsLastRetry(e)
			retryState.add(batchIdx, idx)
			retry = 1
		} else {
			emailDetails.ErrorMsg = sql.NullString{String: e, Valid: true}
		}
		v.incProgress(0, 1, retry)
		return
	}

	if ret == nil || ret.SMTP == nil {
		emailDetails.IsHostExists = false
		emailDetails.ErrorMsg = sql.NullString{String: "ret or SMTP is nil", Valid: true}
		v.incProgress(1, 1, 0)
		return
	}

	emailDetails.IsCatchAll = ret.SMTP.CatchAll
	emailDetails.Reachable = ret.Reachable
	emailDetails.IsDeliverable = ret.SMTP.Deliverable
	emailDetails.IsHostExists = ret.SMTP.HostExists
	emailDetails.IsDisposable = ret.Disposable
	emailDetails.IsInboxFull = ret.SMTP.FullInbox

	v.incProgress(1, 0, 0)
}
