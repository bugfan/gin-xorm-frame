package models

import (
	"encoding/base32"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandByte(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func OtpSecret(size int) string {
	return strings.TrimRight(base32.StdEncoding.EncodeToString(RandByte(size)), "=")
}
