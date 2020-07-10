package util

import (
	"fmt"
	"log"
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
