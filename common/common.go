package common

import "strings"

func IsDuplicateKeyErr(err error) bool {
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}
