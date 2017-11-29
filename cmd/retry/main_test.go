package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_Exec_Fails(t *testing.T) {
	var status int
	application{
		Args:     []string{"cmd", "unknown"},
		Shutdown: func(code int) { status = code },
	}.Run()

	assert.Equal(t, 1, status)
}
