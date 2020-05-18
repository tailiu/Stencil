package auxiliary

import (
	"time"
	"math"
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

func RandomNumber(min, max int) int {
	// Init()
	if max-min <= 0 {
		return min
	}
	return rand.Intn(max-min) + min
}

func RandomNonnegativeInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(math.MaxInt64)
}

func RandomNonnegativeIntWithUpperBound(upperBound int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(upperBound)
}
