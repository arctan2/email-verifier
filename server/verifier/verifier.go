package verifier

import (
	"database/sql"
	"email_verify/db"
	"email_verify/schema"
	"email_verify/socket"
	"maps"
	"strconv"
	"strings"
	"sync"
	"time"
	// emailverifier "github.com/AfterShip/email-verifier"
)

type ProgressData struct {
	Total int `json:"total"`
	Progress int `json:"progress"`
	Success int `json:"success"`
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

func (p *ProgressData) incProgress() {
	p.Lock()
	p.Progress++
	p.Unlock()
}

func (v *Verifier) incProgress(success int) {
	idx := len(v.CurrentProgressList) - 1
	
	p := v.CurrentProgressList[idx]
	p.incProgress()

	p.Success += success

	if p.Progress != 0 && p.Progress % 10 == 0 {
		v.Emit("progress", strconv.Itoa(p.Progress))
	}
}

func (v *Verifier) Emit(ev string, data string) {
	if v.ws != nil {
		v.ws.Emit(ev, data)
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
	v.CompletedBatches = make(map[int][]*ProgressData)
	v.CurrentProgressList = []*ProgressData{}

	truncateBatch := func() {
		for i := range v.CurrentBatch {
			v.CurrentBatch[i] = schema.EmailDetails{}
		}
	}

	i := 0

	socket.EmitWs(v.ws, "get-verifier-details-res", v.VerifierData)

	for ; i+batchSize < len(emails); i += batchSize {
		v.CurrentProgressList = []*ProgressData{NewProgressData(batchSize)}
		v.Emit("batch-start", strconv.Itoa(v.CurrentBatchNumber))

		v.verifyBatch(emails, i, i+batchSize, delay, retryRate)

		truncateBatch()

		v.CompletedBatches[v.CurrentBatchNumber] = make([]*ProgressData, len(v.CurrentProgressList))
		copy(v.CompletedBatches[v.CurrentBatchNumber], v.CurrentProgressList)

		v.Emit("delay", "")
		time.Sleep(time.Duration(delay) * time.Millisecond)
		v.CurrentBatchNumber++
	}

	if i < len(emails) {
		v.CurrentProgressList = []*ProgressData{NewProgressData(len(emails) - i)}
		v.Emit("batch-start", strconv.Itoa(v.CurrentBatchNumber))

		v.verifyBatch(emails, i, len(emails), delay, retryRate)

		v.CompletedBatches[v.CurrentBatchNumber] = make([]*ProgressData, len(v.CurrentProgressList))
		copy(v.CompletedBatches[v.CurrentBatchNumber], v.CurrentProgressList)
	}

	v.State = DONE

	socket.EmitWs(v.ws, "get-verifier-details-res", v.VerifierData)
	return nil
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
	maps.Copy(s.toRetryIdxs, s.idxs)
	s.idxs = make(map[int]int)
}

func (v *Verifier) verifyBatch(emails []string, from, to, delay, retryRate int) {
	retryState := RetryState{make(map[int]int), make(map[int]int), sync.Mutex{}}
	wg := sync.WaitGroup{}
	for i, j := from, 0; i < to; i, j = i+1, j+1 {
		email := strings.TrimSpace(emails[i])
		if email == "" {
			continue
		}
		wg.Add(1)
		go v.verifyEmail(email, j, i, &retryState, false, &wg)
	}
	wg.Wait()

	retryState.reset()

	if len(retryState.toRetryIdxs) == 0 {
		v.Emit("after-all-retries", "0")
		return
	}

	for i := 0; i < retryRate; i++ {
		v.Emit("delay", "")
		time.Sleep(time.Duration(delay) * time.Millisecond)
		l := len(retryState.toRetryIdxs)

		p := NewProgressData(l)

		v.CurrentProgressList = append(v.CurrentProgressList, p)
		socket.EmitWs(v.ws, "retry-begin", *p)
		v.retryBatch(emails, i+1 == retryRate, &retryState)
		if len(retryState.idxs) == 0 {
			v.Emit("after-all-retries", "0")
			return
		}
		retryState.reset()
	}

	v.Emit("after-all-retries", strconv.Itoa(len(retryState.toRetryIdxs)))
}

func (v *Verifier) retryBatch(emails []string, isLastRetry bool, retryState *RetryState) {
	wg := sync.WaitGroup{}
	for batchIdx, idx := range retryState.toRetryIdxs {
		email := strings.TrimSpace(emails[idx])
		if email == "" {
			continue
		}
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
	// 		EnableCatchAllCheck().
	// 		Proxy("socks5://139.59.24.173:1080?timeout=5s")
	// )

	var (
		verifier = NewTestVerifier().
			EnableSMTPCheck().
			EnableCatchAllCheck().
			Proxy("socks5://139.59.24.173:1080?timeout=5s")
	)

	checkIsLastRetry := func(e string) {
		if(isLastRetry) {
			v.CurrentBatch[batchIdx].ErrorMsg = sql.NullString{String: e, Valid: true}
		}
	}

	defer wg.Done()


	ret, err := verifier.Verify(email)

	emailDetails := &v.CurrentBatch[batchIdx]

	if err != nil {
		e := err.Error()
		if strings.Contains(e, "no such host") {
			checkIsLastRetry(e)
			retryState.add(batchIdx, idx)
		} else if strings.Contains(e, "has timed out") {
			checkIsLastRetry(e)
			retryState.add(batchIdx, idx)
		} else if strings.Contains(e, "i/o timeout") {
			checkIsLastRetry(e)
			retryState.add(batchIdx, idx)
		} else if strings.Contains(e, "temporarily unavailable") {
			checkIsLastRetry(e)
			retryState.add(batchIdx, idx)
		}
		emailDetails.ErrorMsg = sql.NullString{String: e, Valid: true}
		v.incProgress(0)
		return
	}

	if !ret.Syntax.Valid {
		emailDetails.IsValidSyntax = false
		v.incProgress(0)
		return
	}

	if ret == nil || ret.SMTP == nil {
		emailDetails.IsHostExists = false
		v.incProgress(0)
		return
	}

	emailDetails.FileId = v.File.Id
	emailDetails.EmailId = email
	emailDetails.IsValidSyntax = ret.Syntax.Valid
	emailDetails.IsReachable = ret.Reachable == "yes"
	emailDetails.IsDeliverable = ret.SMTP.Deliverable
	emailDetails.IsHostExists = ret.SMTP.HostExists
	emailDetails.HasMxRecords = ret.HasMxRecords
	emailDetails.IsDisposable = ret.Disposable
	emailDetails.IsCatchAll = ret.SMTP.CatchAll
	emailDetails.IsInboxFull = ret.SMTP.FullInbox
	emailDetails.ErrorMsg = sql.NullString{String: "", Valid: true}
	v.incProgress(1)
}
