package authentify

import "time"

type Options struct {
	PrefixChars, CodeChars []rune
	CodeLength, SaltLength uint
	Repo                   Repo
	TTL                    time.Duration
}
