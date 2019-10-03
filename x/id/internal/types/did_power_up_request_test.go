package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestDidPowerUpRequest_Validate_EmptyStatus(t *testing.T) {
	request := DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
	}
	assert.Nil(t, request.Validate())
}

func TestDidPowerUpRequest_Validate_InvalidStatusType(t *testing.T) {
	request := DidPowerUpRequest{
		Status: &RequestStatus{
			Type:    "",
			Message: "message",
		},
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
	}
	err := request.Validate()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
}
