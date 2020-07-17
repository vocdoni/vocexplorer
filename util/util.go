package util

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// ErrPrint prints an error to stdout. If err is nil, return false. If err is not nil, return true
func ErrPrint(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
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

//SToTime turns a seconds  inf32 into a readable datetime string
func SToTime(seconds int32) string {
	return time.Unix(int64(seconds), 0).String()
}

// GetAPIStatus returns true if the APIList contains the given target api
func GetAPIStatus(target string, APIList []string) bool {
	for _, api := range APIList {
		if api == target {
			return true
		}
	}
	return false
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
