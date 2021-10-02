package authentify

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"github.com/itsabgr/go-handy"
	"io"
)

type Salt interface {
	Bytes() []byte
	Equal([]byte) bool
	Hex() string
}

type salt struct {
	b []byte
}

func (s *salt) Equal(bytes []byte) bool {
	return subtle.ConstantTimeCompare(bytes, s.Bytes()) == 1
}

func (s *salt) Bytes() []byte {
	return s.b
}

func (s *salt) Hex() string {
	return hex.EncodeToString(s.Bytes())
}
func GenSalt(len uint) (Salt, error) {
	b := make([]byte, len)
	_, err := io.ReadFull(rand.Reader, b)
	handy.Throw(err)
	return AsSalt(b)
}

func AsSalt(b []byte) (Salt, error) {
	return &salt{b}, nil
}
