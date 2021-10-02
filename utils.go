package authentify

import (
	"crypto/rand"
	"errors"
	"github.com/itsabgr/go-handy"
	"io"
	"math"
	"strings"
)

func RandString(chars []rune, length uint) string {

	if len(chars) == 0 {
		panic(errors.New("empty chars"))
	}
	if len(chars) > math.MaxUint8 {
		panic(errors.New("too many chars"))
	}
	if length <= 0 {
		panic(errors.New("zero length"))
	}
	charsLen := byte(len(chars))
	nums := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, nums)
	if err != nil {
		panic(err)
	}
	str := strings.Builder{}
	str.Grow(int(length))
	for i := range handy.N(length) {
		char := chars[nums[i]%charsLen]
		str.WriteRune(char)
	}
	return str.String()
}
