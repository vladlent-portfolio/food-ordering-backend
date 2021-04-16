package common

import (
	"math/rand"
	"strings"
	"time"
)

func IsDuplicateKeyErr(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}

func RandomInt(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}
