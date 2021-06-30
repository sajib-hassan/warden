package mfa

import "net/http"

const (
	KeyPathaoChallengeToken = "Pathao-Challenge-Token"
)

func SetChallengeHeader(w http.ResponseWriter, t *TwoFa) {
	w.Header().Set(KeyPathaoChallengeToken, t.Token)
}

func GetChallengeHeader(r *http.Request) string {
	return r.Header.Get(KeyPathaoChallengeToken)
}
