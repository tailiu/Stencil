/*
 * Helping Functions, because Go is shit and Python is <3
 */

package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Linebreak(ch string, times ...int) {

	t := 1

	if len(times) > 0 {
		t = times[0]
	}

	for i := 0; i < t; i++ {
		fmt.Printf(ch)
	}
	fmt.Println()
}

func Contains(list []string, str string) bool {
	for _, v := range list {
		// fmt.Println(v, "==", str)
		if strings.ToLower(v) == strings.ToLower(str) {
			return true
		}
	}
	return false
}

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomNumber(min, max int) int {
	Init()
	return rand.Intn(max-min) + min
}

func RandomChars(n int) string {
	Init()
	var letters = []rune("zyxwvutsrqponmlkjihgfedcbaZYXWVUTSRQPONMLKJIHGFEDCBAabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomString(n int) string {
	Init()
	var letters = []rune("_9876543210zyxwvutsrqponmlkjihgfedcba0123456789_ZYXWVUTSRQPONMLKJIHGFEDCBA0123456789_abcdefghijklmnopqrstuvwxyz_9876543210ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
	b := make([]rune, n-1)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return RandomChars(1) + string(b)
}
