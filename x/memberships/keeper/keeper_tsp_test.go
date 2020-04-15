package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/memberships/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_AddTrustedServiceProvider(t *testing.T) {
	tests := []struct {
		name         string
		existingTsps ctypes.Addresses
		newTsp       sdk.AccAddress
		expected     ctypes.Addresses
	}{
		{
			name:         "Empty list is updated properly",
			existingTsps: ctypes.Addresses{},
			newTsp:       testTsp,
			expected:     ctypes.Addresses{testTsp},
		},
		{
			name:         "Existing list is updated properly",
			existingTsps: ctypes.Addresses{testTsp},
			newTsp:       testUser,
			expected:     ctypes.Addresses{testTsp, testUser},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, t := range test.existingTsps {
				k.AddTrustedServiceProvider(ctx, t)
			}

			k.AddTrustedServiceProvider(ctx, test.newTsp)

			stored := k.GetTrustedServiceProviders(ctx)
			require.Equal(t, test.expected, stored)
		})
	}
}

func TestKeeper_GetTrustedServiceProviders(t *testing.T) {
	tests := []struct {
		name string
		tsps ctypes.Addresses
	}{
		{
			name: "Empty list is returned properly",
			tsps: ctypes.Addresses{},
		},
		{
			name: "Existing list is returned properly",
			tsps: ctypes.Addresses{testTsp, testUser, testUser2},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, t := range test.tsps {
				k.AddTrustedServiceProvider(ctx, t)
			}

			signers := k.GetTrustedServiceProviders(ctx)

			for _, s := range test.tsps {
				require.Contains(t, signers, s)
			}
		})
	}
}

func TestKeeper_IsTrustedServiceProvider(t *testing.T) {
	tests := []struct {
		name     string
		tsps     ctypes.Addresses
		address  sdk.AccAddress
		expected bool
	}{
		{
			name:     "Empty list returns false",
			tsps:     ctypes.Addresses{},
			address:  testTsp,
			expected: false,
		},
		{
			name:     "Not present TSP returns false",
			tsps:     ctypes.Addresses{testTsp},
			address:  testUser,
			expected: false,
		},
		{
			name:     "Present TSP returns true",
			tsps:     ctypes.Addresses{testTsp, testUser},
			address:  testUser,
			expected: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			if !test.tsps.Empty() {
				store.Set([]byte(types.TrustedSignersStoreKey), k.Cdc.MustMarshalBinaryBare(&test.tsps))
			}

			require.Equal(t, test.expected, k.IsTrustedServiceProvider(ctx, test.address))
		})
	}
}
