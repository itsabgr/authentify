package authentify

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Code interface {
	fmt.Stringer
	GetCode() string
	GetPrefix() string
	Hash() (HashedCode, error)
}
type HashedCode interface {
	Bytes() []byte
	Check(a Code) error
}

type hashedCode struct {
	bytes []byte
}

func (h *hashedCode) Bytes() []byte {
	return h.bytes
}

func (h *hashedCode) Check(code Code) error {
	return bcrypt.CompareHashAndPassword(h.Bytes(), []byte(code.String()))
}

func AsHashedCode(hash []byte) (HashedCode, error) {
	return &hashedCode{bytes: hash}, nil
}

type code struct {
	prefix, code string
}

func (c *code) Hash() (HashedCode, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(c.String()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return AsHashedCode(hash)
}

func (c *code) GetCode() string {
	return c.code
}

func (c *code) GetPrefix() string {
	return c.prefix
}

func (c *code) String() string {
	return c.GetPrefix() + "-" + c.GetCode()
}

const (
	DefaultPrefixChars = "ABC"
	DefaultCodeLength  = 6
	DefaultCodeChars   = "1234567890"
)

func NewCode(prefix, c string) (Code, error) {
	return &code{
		prefix: prefix,
		code:   c,
	}, nil
}
