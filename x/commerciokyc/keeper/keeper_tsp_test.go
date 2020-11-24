package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"

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
		{
			name:         "Existing tsp is updated properly",
			existingTsps: ctypes.Addresses{testTsp},
			newTsp:       testTsp,
			expected:     ctypes.Addresses{testTsp},
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
		govAddr  sdk.AccAddress
		expected bool
	}{
		{
			name:     "Empty list returns false",
			tsps:     ctypes.Addresses{},
			address:  testTsp,
			govAddr:  testUser3,
			expected: false,
		},
		{
			name:     "Not present TSP returns false",
			tsps:     ctypes.Addresses{testTsp},
			address:  testUser,
			govAddr:  testUser3,
			expected: false,
		},
		{
			name:     "Present TSP returns true",
			tsps:     ctypes.Addresses{testTsp, testUser},
			address:  testUser,
			govAddr:  testUser3,
			expected: true,
		},
		{
			name:     "Government address returns true",
			tsps:     ctypes.Addresses{testTsp, testUser},
			address:  testUser3,
			govAddr:  testUser3,
			expected: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, gk, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			if !test.tsps.Empty() {
				store.Set([]byte(types.TrustedSignersStoreKey), k.Cdc.MustMarshalBinaryBare(&test.tsps))
			}
			_ = gk.SetGovernmentAddress(ctx, test.govAddr)

			require.Equal(t, test.expected, k.IsTrustedServiceProvider(ctx, test.address))
		})
	}
}

func TestKeeper_RemoveTrustedServiceProviders(t *testing.T) {
	tests := []struct {
		name         string
		existingTsps ctypes.Addresses
		expected     ctypes.Addresses
		deleteTsp    sdk.AccAddress
	}{
		{
			name:         "Not present TSP returns same list",
			existingTsps: ctypes.Addresses{testTsp},
			expected:     ctypes.Addresses{testTsp},
			deleteTsp:    testUser,
		},
		{
			name:         "Present TSP returns tsp's list without it",
			existingTsps: ctypes.Addresses{testTsp, testUser},
			expected:     ctypes.Addresses{testTsp},
			deleteTsp:    testUser,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, t := range test.existingTsps {
				k.AddTrustedServiceProvider(ctx, t)
			}

			k.RemoveTrustedServiceProvider(ctx, test.deleteTsp)
			stored := k.GetTrustedServiceProviders(ctx)
			require.Equal(t, test.expected, stored)

		})
	}
}
