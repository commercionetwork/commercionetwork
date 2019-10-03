package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var requestSender, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")
var requestRecipient, _ = sdk.AccAddressFromBech32("cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa")

func TestDidDepositRequest_Validate_EmptyStatus(t *testing.T) {
	request := DidDepositRequest{
		FromAddress:   requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
		Recipient:     requestRecipient,
	}
	assert.Nil(t, request.Validate())
}

func TestDidDepositRequest_Validate_InvalidStatusType(t *testing.T) {
	request := DidDepositRequest{
		Status: &RequestStatus{
			Type:    "",
			Message: "message",
		},
		FromAddress:   requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
		Recipient:     requestRecipient,
	}
	err := request.Validate()
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
}
