package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func GenerateHashedPassword(pwd string, p *ArgonParams) (hash string, err error) {
	// create random number in byte format
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return
	}

	// generate hash password
	rawHash := argon2.IDKey([]byte(pwd), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(rawHash)

	hash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return
}

func ComparePasswordAndHash(pwd, hashed string) (match bool, err error) {
	// extract parameter from hashed password from database
	p, salt, hash, err := decodeHashedPassword(hashed)
	if err != nil {
		return false, err
	}

	// get hashed from incoming password
	incomingHashedPassword := argon2.IDKey([]byte(pwd), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// return match
	if subtle.ConstantTimeCompare(hash, incomingHashedPassword) == 1 {
		return true, nil
	}

	// return not match but no error
	return false, nil
}

// generate random bytes to salt
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// decode hashed password
func decodeHashedPassword(hashed string) (p *ArgonParams, salt, hash []byte, err error) {
	v := strings.Split(hashed, "$")
	if len(v) != 6 {
		err = ErrInvalidHash
		return
	}

	// check version argon
	var version int
	_, err = fmt.Sscanf(v[2], "v=%d", &version)
	if err != nil {
		return
	}

	if version != argon2.Version {
		err = ErrIncompatibleVersion
		return
	}

	// get all params from hashed password
	p = &ArgonParams{}
	_, err = fmt.Sscanf(v[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return
	}

	// get salt
	salt, err = base64.RawStdEncoding.Strict().DecodeString(v[4])
	if err != nil {
		return
	}

	// get keyLength
	p.SaltLength = uint32(len(salt))

	// get hash
	hash, err = base64.RawStdEncoding.Strict().DecodeString(v[5])
	if err != nil {
		return
	}

	// get keyLength
	p.KeyLength = uint32(len(hash))

	return
}
