package util

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"go.vocdoni.io/dvote/types"
	"go.vocdoni.io/proto/build/go/models"
)

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

// EncodeInt encodes an integer to a byte array
func EncodeInt(val interface{}) []byte {
	var val64 int64
	buf := make([]byte, binary.MaxVarintLen64)
	if i, ok := val.(int64); ok {
		val64 = int64(i)
		goto encode
	}
	if i, ok := val.(int32); ok {
		val64 = int64(i)
		goto encode
	}
	if i, ok := val.(uint64); ok {
		val64 = int64(i)
		goto encode
	}
	if i, ok := val.(uint32); ok {
		val64 = int64(i)
		goto encode
	}
	if i, ok := val.(int); ok {
		val64 = int64(i)
		goto encode
	} else {
		logger.Error(fmt.Errorf("cannot encode value: type is not integer"))
	}
encode:
	binary.PutVarint(buf, val64)
	return buf
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

// StringInSlice returns true if given string is in the list of strings
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// TrimSlice trims a slice of strings to lim elements. If rev is set to true, trims from beginning rather than end.
func TrimSlice(slice []string, lim int, page *int) []string {
	if *page < 0 {
		*page = 0
		logger.Info("Invalid page number")
	}
	len := len(slice)
	if (*page+1)*lim > len+lim-1 {
		*page = (len - 1) / lim
	}
	start := Min(0+(*page*lim), len)
	end := Max(Min(len, (*page+1)*lim), 0)
	return slice[start:end]
}

//TrimHex removes the hex prefix from a string
func TrimHex(str string) string {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		return str[2:]
	}
	return str
}

//HexToString converts an array of hexbytes to a string
func HexToString(bytes []byte) string {
	return strings.ToLower(hex.EncodeToString(bytes))
}

// GetTransactionType translates a raw transaction to a type string
func GetTransactionType(raw *models.Tx) string {
	switch raw.Payload.(type) {
	case *models.Tx_Vote:
		return types.TxVote
	case *models.Tx_NewProcess:
		return types.TxNewProcess
	case *models.Tx_CancelProcess:
		return types.TxCancelProcess
	case *models.Tx_Admin:
		return "admin"
	case *models.Tx_SetProcess:
		return "setProcess"
	}
	return "unknown"
}

// GetTransactionName translates a raw transaction type to a name
func GetTransactionName(raw string) string {
	name, ok := config.TransactionTypeMap[raw]
	if ok {
		return name
	}
	return raw
}

// GetProcessName translates a raw process type to a name
func GetProcessName(raw string) string {
	name, ok := config.ProcessTypeMap[raw]
	if ok {
		return name
	}
	return raw
}

// GetEnvelopeName translates a raw envelope type to a name
func GetEnvelopeName(raw string) string {
	name, ok := config.EnvelopeTypeMap[raw]
	if ok {
		return name
	}
	return raw
}
