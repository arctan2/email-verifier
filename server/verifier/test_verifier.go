package verifier

import (
	"errors"
	"math/rand"
	"time"

	emailverifier "github.com/AfterShip/email-verifier"
)

type TestVerifier struct {
}

func NewTestVerifier() *TestVerifier {
	return &TestVerifier{}
}

func (t *TestVerifier) EnableSMTPCheck() *TestVerifier {
	return t
}

func (t *TestVerifier) EnableCatchAllCheck() *TestVerifier{
	return t
}

func (t *TestVerifier) Proxy(s string) *TestVerifier {
	return t
}

func randBool() bool {
	return rand.Intn(2) == 1
}

func (t *TestVerifier) Verify(s string) (*emailverifier.Result, error) {
	errs := []string{
		"no such host", "has timed out", "i/o timeout", "temporarily unavailable", "you're blocked", "catastrophic damage",
	}

	var e error = nil

	if randBool() {
		e = errors.New(errs[rand.Intn(len(errs))])
	}

	MAX_MS := 1000
	MIN_MS := 100
	time.Sleep(time.Duration(rand.Intn(MAX_MS - MIN_MS + 1) + MIN_MS) * time.Millisecond)

	res := emailverifier.Result{
		Email: s,
		Reachable: "yes",
		Syntax: emailverifier.Syntax{ Username: "", Domain: "", Valid: randBool() },
		SMTP: &emailverifier.SMTP{
			HostExists: true,
			FullInbox: false,
			CatchAll: true,
			Deliverable: true,
			Disabled: false,
		},
		Gravatar: nil,
		Suggestion: "",
		Disposable: false,
		RoleAccount: false,
		Free: true,
		HasMxRecords: true,
	}

	return &res, e
}

