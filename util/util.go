package util

import (
	"fmt"
	"log"
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
