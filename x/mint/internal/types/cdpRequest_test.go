package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestCdpRequest = CDPRequest{
	Signer:          TestOwner,
	DepositedAmount: TestDepositedAmount,
	Timestamp:       TestTimestamp,
}

func TestCDPRequest_Equals_true(t *testing.T) {
	cdpReq := TestCdpRequest
	actual := TestCdpRequest.Equals(cdpReq)
	assert.True(t, actual)
}

func TestCDPRequest_Equals_false(t *testing.T) {
	cdpReq := CDPRequest{
		Signer:          nil,
		DepositedAmount: nil,
		Timestamp:       "",
	}

	actual := TestCdpRequest.Equals(cdpReq)
	assert.False(t, actual)
}
