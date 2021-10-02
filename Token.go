package authentify

import (
	"bytes"
	"crypto/subtle"
	"encoding/gob"
	"github.com/google/uuid"
	"io"
	"time"
)

type Token interface {
	Deadline() time.Time
	Validate(sender string, code Code, salt Salt) bool
	Receiver() string
	EncodeTo(dst io.Writer) error
	Encode() ([]byte, error)
	Sender() string
	Salt() Salt
}
type token struct {
	salt       Salt
	receiver   string
	senderName string
	code       HashedCode
	deadline   time.Time
}

func (t *token) Sender() string {
	return t.senderName
}

func (t *token) Salt() Salt {
	return t.salt
}

func (t *token) Encode() ([]byte, error) {
	b := &bytes.Buffer{}
	err := t.EncodeTo(b)
	return b.Bytes(), err
}

type TokenID = uuid.UUID

func (t *token) Deadline() time.Time {
	return t.deadline
}

func (t *token) Validate(sender string, code Code, salt Salt) bool {
	ok := subtle.ConstantTimeCompare([]byte(sender), []byte(t.Sender())) == 1
	ok2 := t.Salt().Equal(salt.Bytes())
	err := t.code.Check(code)
	return ok2 && ok && err == nil
}

func (t *token) Receiver() string {
	return t.receiver
}

type tokenStruct struct {
	Receiver string
	Hash     []byte
	Deadline time.Time
	Salt     []byte
	Sender   string
}

func (t *token) EncodeTo(dst io.Writer) error {
	return gob.NewEncoder(dst).Encode(&tokenStruct{
		Receiver: t.Receiver(),
		Hash:     t.code.Bytes(),
		Deadline: t.Deadline(),
		Salt:     t.salt.Bytes(),
		Sender:   t.Sender(),
	})
}
func DecodeToken(b []byte) (Token, error) {
	structure := &tokenStruct{}
	err := gob.NewDecoder(bytes.NewReader(b)).Decode(structure)
	if err != nil {
		return nil, err
	}
	hash, err := AsHashedCode(structure.Hash)
	if err != nil {
		return nil, err
	}
	salt, err := AsSalt(structure.Salt)
	if err != nil {
		return nil, err
	}
	return NewToken(structure.Sender, structure.Receiver, structure.Deadline, hash, salt)
}
func NewToken(senderName, to string, deadline time.Time, code HashedCode, salt Salt) (Token, error) {
	token := &token{}
	token.receiver = to
	token.code = code
	token.deadline = deadline
	token.salt = salt
	token.senderName = senderName
	return token, nil
}
