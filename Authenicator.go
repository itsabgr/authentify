package authentify

import (
	"github.com/itsabgr/go-handy"
	"time"
)

type authenticator struct {
	opt Options
}

type Authenticator interface {
	SendCode(sender Sender, to string) (_ Token, prefix string, _ error)
	RetrieveToken(receiver string) (Token, error)
}

func NewAuthenticator(opt Options) (Authenticator, error) {
	return &authenticator{
		opt: opt,
	}, nil
}

func (a *authenticator) SendCode(sender Sender, to string) (_ Token, prefix string, _ error) {
	deadline := time.Now().Add(a.opt.TTL)
	code, err := NewCode(RandString(a.opt.PrefixChars, 1), RandString(a.opt.CodeChars, a.opt.CodeLength))
	handy.Throw(err)
	hash, err := code.Hash()
	handy.Throw(err)
	salt, err := GenSalt(a.opt.SaltLength)
	handy.Throw(err)
	token, err := NewToken(sender.Name(), to, deadline, hash, salt)
	handy.Throw(err)
	tokenBin, err := token.Encode()
	handy.Throw(err)
	err = sender.Send(to, code, deadline)
	if err != nil {
		return nil, "", err
	}

	err = a.opt.Repo.Store(to, tokenBin, deadline)
	handy.Throw(err)
	return token, code.GetPrefix(), nil
}

func (a *authenticator) RetrieveToken(receiver string) (Token, error) {
	val, err := a.opt.Repo.FindByID(receiver)
	if err != nil {
		return nil, err
	}
	return DecodeToken(val)
}
