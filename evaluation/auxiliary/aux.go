package auxiliary

import (
	"time"
	"math/rand"
)

func RandStrSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func RandomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2147483647)
}
