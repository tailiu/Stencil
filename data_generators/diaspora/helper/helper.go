package helper

import (
	"math/rand"
	"strings"
	"time"
)

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomNumber(min, max int) int {
	Init()
	if max-min <= 0 {
		return min
	}
	return rand.Intn(max-min) + min
}

func RandomString(n int) string {
	Init()
	var letters = []rune("_9876543210zyxwvutsrqponmlkjihgfedcba0123456789_ZYXWVUTSRQPONMLKJIHGFEDCBA0123456789_abcdefghijklmnopqrstuvwxyz_9876543210ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomText(n int) string {
	Init()
	var letters = []rune(" 9876543210 zyxwvutsrqponmlkjihgfedcba 0123456789 ZYXWVUTSRQPONMLKJIHGFEDCBA 0123456789 abcdefghijklmnopqrstuvwxyz 9876543210ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func Contains(list []string, str string) bool {
	for _, v := range list {
		// fmt.Println(v, "==", str)
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}
