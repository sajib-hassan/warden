package encryptor

import (
	"github.com/alexedwards/argon2id"
)

var argon2Params *argon2id.Params

func init() {
	// DefaultParams provides some sane default parameters for hashing passwords.
	//
	// Follows recommendations given by the Argon2 RFC:
	// "The Argon2id variant with t=1 and maximum available memory is RECOMMENDED as a
	// default setting for all environments. This setting is secure against side-channel
	// attacks and maximizes adversarial costs on dedicated bruteforce hardware.""
	//
	// The default parameters should generally be used for development/testing purposes
	// only. Custom parameters should be set for production applications depending on
	// available memory/CPU resources and business requirements.
	argon2Params = &argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  1,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func GenerateFromPassword(password string) (encodedHash string, err error) {
	// CreateHash returns a Argon2id hash of a plain-text password using the
	// provided algorithm parameters. The returned hash follows the format used
	// by the Argon2 reference C implementation and looks like this:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	return argon2id.CreateHash(password, argon2Params)
}

func ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(password, encodedHash)
}
