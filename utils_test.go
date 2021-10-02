package authentify

import (
	"errors"
	"fmt"
	"github.com/itsabgr/go-handy"
	rand2 "math/rand"
	"testing"
	"unicode/utf8"
)

func TestRandString(t *testing.T) {
	defer handy.Catch(func(recovered interface{}) {
		t.Fatal(recovered)
	})
	for range handy.N(10) {
		chars := []rune("1234567Y UIOPقعغفبسدهبغ۴قلا.یبنملپیبمن😃😃😃😃😃😃😃😃😃😃")
		length := uint(rand2.Int31n(100) + 1)
		str := RandString(chars, length)
		strLen := uint(utf8.RuneCountInString(str))
		handy.Assert(strLen == length, fmt.Errorf("expected %d got %d", length, strLen))
	loop:
		for _, got := range str {
			for _, allowed := range chars {
				if allowed == got {
					continue loop
				}
			}
			panic(errors.New("invalid char"))
		}
	}
}
