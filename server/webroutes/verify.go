package webroutes

import (
	"email_verify/respond"
	"email_verify/schema"
	"email_verify/verifier"
	"encoding/json"
	"net/http"
	"sync"
)

func verify(email string, results *[]schema.EmailDetails, wg *sync.WaitGroup, mutex *sync.Mutex) {
	v := verifier.VerifyEmail(email)
	mutex.Lock()
	*results = append(*results, v)
	mutex.Unlock()
	wg.Done()
}

func (m *WebRoutesHandler) verifyEmails(w http.ResponseWriter, r *http.Request) {
	var emails []string
	err := json.NewDecoder(r.Body).Decode(&emails)
	if err != nil {
		respond.RespondErrMsg(w, err.Error())
		return
	}

	var results []schema.EmailDetails

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, email := range emails {
		wg.Add(1)
		go verify(email, &results, &wg, &mutex)
	}

	wg.Wait()

	res := struct {
		respond.ResponseStruct
		Results []schema.EmailDetails `json:"results"`
	} {
		Results: results,
	}

	json.NewEncoder(w).Encode(res)
}
