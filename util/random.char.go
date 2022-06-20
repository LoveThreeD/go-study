package util

import (
	"math/rand"
	"time"
	"unsafe"
)

const (
	letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	_lettersAccount = "ABCDEFGHIJKMNPQRSTUVWXYZ"
	_lettersPasswd  = "abcdefghijkmnpqrstuvwxyz123456789"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIDBits = 6
	// All 1-bits as many as letterIdBits
	letterIDMask = 1<<letterIDBits - 1
	letterIDMax  = 63 / letterIDBits
)

func RandNChar(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIDMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIDMax
		}
		if idx := int(cache & letterIDMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIDBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func RandNCharAccount(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIDMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIDMax
		}
		if idx := int(cache & letterIDMask); idx < len(_lettersAccount) {
			b[i] = _lettersAccount[idx]
			i--
		}
		cache >>= letterIDBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func RandNCharPasswd(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIDMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIDMax
		}
		if idx := int(cache & letterIDMask); idx < len(_lettersPasswd) {
			b[i] = _lettersPasswd[idx]
			i--
		}
		cache >>= letterIDBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
