package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestCdpRequest = CdpRequest{
	Signer:          TestOwner,
	DepositedAmount: TestDepositedAmount,
	Timestamp:       TestTimestamp,
}

func TestCdpRequest_Equals_true(t *testing.T) {
	cdpReq := TestCdpRequest
	actual := TestCdpRequest.Equals(cdpReq)
	assert.True(t, actual)
}

func TestCdpRequest_Equals_false(t *testing.T) {
	cdpReq := CdpRequest{
		Signer:          nil,
		DepositedAmount: nil,
		Timestamp:       "",
	}

	actual := TestCdpRequest.Equals(cdpReq)
	assert.False(t, actual)
}
