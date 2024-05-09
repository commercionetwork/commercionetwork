package keeper

import (
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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
			var expected sdk.AccAddress
			if tt.wantErr {
				expected = tt.governmentOld
			} else {
				expected = tt.governmentToSet
			}

			require.Equal(t, expected, k.GetGovernmentAddress(ctx))
		})
	}
}

func TestKeeper_GetGovernmentAddress(t *testing.T) {

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

func TestKeeper_GetGovernment300Address(t *testing.T) {
	type fields struct {
		cdc      codec.Codec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey
	}
	type args struct {
		ctx sdk.Context
	}
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeperWithV300Government(t, governmentTestAddress)

			if got := k.GetGovernment300Address(ctx); !reflect.DeepEqual(got, governmentTestAddress) {
				t.Errorf("Keeper.GetGovernment300Address() = %v, want %v", got, governmentTestAddress)
			}
		})
	}
}
