package keeper

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_getDepositRequestStoreKey(t *testing.T) {
	tests := []struct {
		name    string
		proof   string
		wantB64 string
	}{
		{
			"get deposit request key from a proof",
			"proof",
			"aWRkZXBvc2l0UmVxdWVzdHByb29m",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantB64, base64.StdEncoding.EncodeToString(getDepositRequestStoreKey(tt.proof)))
		})
	}
}

func TestKeeper_getDidPowerUpRequestStoreKey(t *testing.T) {
	tests := []struct {
		name    string
		proof   string
		wantB64 string
	}{
		{
			"get power up request key request from a proof",
			"proof",
			"aWRwb3dlclVwUmVxdWVzdHByb29m",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantB64, base64.StdEncoding.EncodeToString(getDidPowerUpRequestStoreKey(tt.proof)))
		})
	}
}

func TestKeeper_getIdentityStoreKey(t *testing.T) {
	SetupTestInput()
	tests := []struct {
		name    string
		owner   sdk.AccAddress
		wantB64 string
	}{
		{
			"get identity store key request from a proof",
			TestGovernment,
			"aWQ6aWRlbnRpdGllczpQVVgsvYJiWsmPrk8DD+XET1lIDg==",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantB64, base64.StdEncoding.EncodeToString(getIdentityStoreKey(tt.owner)))
		})
	}
}

func Test_getHandledPowerUpRequestsReferenceStoreKey(t *testing.T) {
	SetupTestInput()
	tests := []struct {
		name      string
		reference string
		want      string
	}{
		{
			"get handled power up request store key request from a proof",
			"reference",
			"aWRoYW5kbGVkUG93ZXJVcFJlcXVlc3RzUmVmZXJlbmNlcmVmZXJlbmNl",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, base64.StdEncoding.EncodeToString(getHandledPowerUpRequestsReferenceStoreKey(tt.reference)))
		})
	}
}
