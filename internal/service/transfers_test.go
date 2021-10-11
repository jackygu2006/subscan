package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_GetTranferList(t *testing.T) {
	_, count := testSrv.GetTransferList(0, 10, map[string]string{})
	assert.Equal(t, 1, count)
}
