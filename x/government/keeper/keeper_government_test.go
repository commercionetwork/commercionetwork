package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SetGovernmentAddress(t *testing.T) {

	tests := []struct {
		name            string
		governmentToSet sdk.AccAddress
		governmentOld   sdk.AccAddress
		wantErr         bool
	}{
		{
			"empty store",
			governmentTestAddress,
			nil,
			false,
		},
		{
			"empty address",
			nil,
			nil,
			true,
		},
		{
			"same government already set",
			governmentTestAddress,
			governmentTestAddress,
			true,
		},
		{
			"new government with government already set",
			notGovernmentAddress,
			governmentTestAddress,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeperWithGovernmentAddress(t, tt.governmentOld)

			if err := k.SetGovernmentAddress(ctx, tt.governmentToSet); (err != nil) != tt.wantErr {
				t.Errorf("Keeper.SetGovernmentAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				require.Equal(t, tt.governmentOld, k.GetGovernmentAddress(ctx))
			} else {
				require.Equal(t, tt.governmentToSet, k.GetGovernmentAddress(ctx))
			}
		})
	}
}

func TestKeeper_GetGovernmentAddress(t *testing.T) {

	governmentTestAddress, err := sdk.AccAddressFromBech32("did:com:1zg4jreq2g57s4efrl7wnh2swtrz3jt9nfaumcm")
	require.NoError(t, err)

	tests := []struct {
		name    string
		address sdk.AccAddress
	}{
		{
			"expected government",
			governmentTestAddress,
		},
		{
			"empty government",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeperWithGovernmentAddress(t, tt.address)
			require.Equal(t, tt.address, k.GetGovernmentAddress(ctx))
		})
	}
}
