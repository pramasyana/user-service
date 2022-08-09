package model

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// IterationsCount to set the iteration
	IterationsCount = 15000
	// SaltSize to set salt size
	SaltSize = 64
)

// PasswordHasher interface abstraction
type PasswordHasher interface {
	Hash(data []byte) []byte
	ParseSalt(salt string) error
	GenerateSalt() string
}

// PBKDF2Hasher data structure
type PBKDF2Hasher struct {
	salt      []byte
	keyLen    int
	hash      func() hash.Hash
	iteration int
}

// NewPBKDF2Hasher function for creating hashing service
func NewPBKDF2Hasher(saltSize, KeyLen, iteration int, hash func() hash.Hash) *PBKDF2Hasher {
	return &PBKDF2Hasher{
		salt:      make([]byte, saltSize),
		keyLen:    KeyLen,
		iteration: iteration,
		hash:      hash,
	}
}

// GenerateSalt function
func (h *PBKDF2Hasher) GenerateSalt() string {
	// initiate new salt to prevent replacing data
	salt := h.salt
	rand.Read(salt)
	saltString := base64.StdEncoding.EncodeToString(salt)

	combinedSalt := fmt.Sprintf("%v.%v", h.iteration, saltString)
	return combinedSalt
}

// Hash PBKDF2Hasher function for hashing string
func (h *PBKDF2Hasher) Hash(data []byte) []byte {
	hashed := pbkdf2.Key(data, h.salt, h.iteration, h.keyLen, h.hash)

	// force replacing salt to prevent bigger byte
	h.salt = make([]byte, h.keyLen)

	return hashed
}

// ParseSalt PBKDF2Hasher function for parsing salt
func (h *PBKDF2Hasher) ParseSalt(salt string) error {
	var iteration int
	var saltOnly string

	_, err := fmt.Sscanf(salt, "%d.%s", &iteration, &saltOnly)

	if err != nil {
		return nil
	}

	h.iteration = iteration

	if cap(h.salt) < len(salt) {
		h.salt = make([]byte, len(salt))
	}

	copy(h.salt, []byte(salt))

	return nil
}
