package randstring

import (
	"math/rand"
	"time"
)

func Rune(letterRuneList []rune, n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRuneList[rand.Intn(len(letterRuneList))]
	}
	return string(b)
}
