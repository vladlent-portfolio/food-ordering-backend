package common

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsDuplicateKeyErr(t *testing.T) {
	assert.True(t, IsDuplicateKeyErr(errors.New("ERROR: duplicate key value violates unique constraint \"categories_pkey\" (SQLSTATE 23505)")))
	assert.False(t, IsDuplicateKeyErr(errors.New("92374283uasdfj")))
}
