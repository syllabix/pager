package pager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResult(t *testing.T) {
	page := PageDetails{10, 1}

	assert.Equal(t, 0, page.Result(0).TotalPageCount)
	assert.Equal(t, 1, page.Result(1).TotalPageCount)
	assert.Equal(t, 1, page.Result(10).TotalPageCount)
	assert.Equal(t, 2, page.Result(11).TotalPageCount)
}
