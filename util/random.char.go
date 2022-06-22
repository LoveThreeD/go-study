package util

import (
	"math/rand"
	"time"
)

var (
	_lettersAccount = []rune("ABCDEFGHIJKMNPQRSTUVWXYZ")
	_lettersPasswd  = []rune("abcdefghijkmnpqrstuvwxyz123456789")
)

func RandNCharAccount(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = _lettersAccount[rand.Intn(len(_lettersAccount)-1)]
	}

	return string(b)

}

func RandNCharPasswd(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = _lettersPasswd[rand.Intn(len(_lettersPasswd)-1)]
	}

	return string(b)
}
