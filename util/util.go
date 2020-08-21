package util

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/vocdoni/go-dvote/log"
)

// ErrPrint prints an error to stdout. If err is nil, return false. If err is not nil, return true
func ErrPrint(err error) bool {
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}

// ErrFatal calls log.Fatal if err is not nil
func ErrFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// MsToString turns a milliseconds int32 to a readable string
func MsToString(ms int32) string {
	seconds, ms := ms/1000, ms%1000
	minutes, seconds := seconds/60, seconds%60
	return fmt.Sprintf("%02d:%02d:%04d", minutes, seconds, ms)
}

// MsToSecondsString only returns the second part of a ms unit, as a string
func MsToSecondsString(ms int32) string {
	seconds := ms / 1000
	seconds = seconds % 60
	return fmt.Sprintf("%02d", seconds)
}

// IntToString takes an int32, int64, or int, and returns a string
func IntToString(val interface{}) string {
	if i, ok := val.(int); ok {
		return strconv.Itoa(i)
	}
	if i, ok := val.(int32); ok {
		return strconv.Itoa(int(i))
	}
	if i, ok := val.(int64); ok {
		return strconv.Itoa(int(i))
	}
	if i, ok := val.(uint32); ok {
		return strconv.Itoa(int(i))
	}
	if i, ok := val.(uint64); ok {
		return strconv.Itoa(int(i))
	}
	return ""
}

// Min returns the min of two ints
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the max of two ints
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SearchSlice searches a slice of strings and returns all strings who have substrings matching the search term
func SearchSlice(source []string, search string) []string {
	var results []string
	for _, str := range source {
		// for i, str := range source {
		if strings.Contains(str, search) {
			results = append(results, str)
			// results = append(results, strconv.Itoa(i))
		}
	}
	return results
}

// TrimSlice trims a slice of strings to lim elements. If rev is set to true, trims from beginning rather than end.
func TrimSlice(slice []string, lim int, page *int) []string {
	if *page < 0 {
		*page = 0
		fmt.Println("Invalid page number")
	}
	len := len(slice)
	if (*page+1)*lim > len+lim-1 {
		*page = (len - 1) / lim
	}
	start := Min(0+(*page*lim), len)
	end := Max(Min(len, (*page+1)*lim), 0)
	return slice[start:end]
}

//StripHexString removes the hex prefix from a string
func StripHexString(str string) string {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		return str[2:]
	}
	return str
}

//HexToString converts an array of hexbytes to a string
func HexToString(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}
