package helpmate

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
)

func SHA256Encode(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func SHA256HMACHash(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func RandSecret(size uint) string {
	if size == 0 {
		size = 20
	}
	secret := make([]byte, size)
	_, _ = rand.Reader.Read(secret)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret)
}

func FormatBDMobile(m string) string {
	return "+88" + m[len(m)-11:]
}
