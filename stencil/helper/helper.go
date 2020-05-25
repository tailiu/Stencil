/*
 * Helping Functions, because Go is shit and Python is <3
 */

package helper

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	logg "github.com/withmandala/go-log"
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
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}

func ContainsIdx(list []string, str string) (int, bool) {
	for i, v := range list {
		if strings.EqualFold(v, str) {
			return i, true
		}
	}
	return -1, false
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

func Sublist(first, second []string) bool {

	set := make(map[string]int)

	for _, value := range first {
		value = strings.ToLower(value)
		set[value] = 1
	}

	for _, value := range second {
		value = strings.ToLower(value)
		if _, found := set[value]; !found {
			fmt.Println("value not found", value)
			fmt.Println("1:", first)
			fmt.Println("2:", second)
			log.Fatal("check: stencil.helper.Sublist()")
			return false
		}
	}
	// fmt.Println("IS SUBSET!!!!")
	return true
}

func GetKeysOfPhyTabMap(m map[string][][]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func Shuffle(vals []int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

func ConcatMaps(a map[string]string, b map[string]string) {
	if b == nil || a == nil {
		return
	}
	for k, v := range b {
		a[k] = v
	}
}
func Trace() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	return fn.Name()
}

func ParseFloat(str string) (float64, error) {

	//Some number is specifed in scientific notation
	pos := strings.IndexAny(str, "eE")
	if pos < 0 {
		return strconv.ParseFloat(str, 64)
	}

	baseStr := str[0:pos]
	baseVal, err := strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return 0, err
	}

	expStr := str[(pos + 1):]
	expVal, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return baseVal * math.Pow10(int(expVal)), nil
}

func GetInt64(val interface{}) int64 {
	if valInt, ok := val.(int64); ok {
		return valInt
	} else if valFloat, ok := val.(float64); ok {
		valInt := int64(math.Ceil(valFloat))
		return valInt
	} else if valString, ok := val.(string); ok {
		if v, err := strconv.ParseInt(valString, 10, 64); err != nil {
			log.Fatalf("@GetInt64: Failed string conversion to int64 | %T | %v | %v", val, val, err)
		} else {
			return v
		}
	}
	log.Fatalf("@GetInt64: Neither int64 nor float64 nor valid stringInt | %T | %v", val, val)
	return 0
}

func ConvertScientificNotationToString(val interface{}) string {
	if valInt, ok := val.(int64); ok {
		return fmt.Sprint(valInt)
	} else if valFloat, ok := val.(float64); ok {
		valInt := int64(math.Ceil(valFloat))
		return fmt.Sprint(valInt)
	}
	return fmt.Sprint(val)
}

func ConvertIntFloatToString(val interface{}) interface{} {
	if valInt, ok := val.(int64); ok {
		return fmt.Sprint(valInt)
	} else if valFloat, ok := val.(float64); ok {
		valInt := int64(math.Ceil(valFloat))
		return fmt.Sprint(valInt)
	}
	return val
}

func CreateLogger(debug bool) *logg.Logger {
	logger := logg.New(os.Stderr)

	logger.WithTimestamp()
	logger.WithColor()
	if debug {
		logger.WithDebug()
	}
	return logger
}
