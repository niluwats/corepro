package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateCode(usage string) string {
	rand.Seed(time.Now().UnixNano())
	strCode := strconv.Itoa(rand.Intn(100000))
	currLength := len(strCode)
	if usage == "verify_mobile" {
		if currLength < 5 {
			return strCode + strings.Repeat("0", (5-currLength))
		}
		return strCode
	} else {
		return Encode(strCode)
	}
}
