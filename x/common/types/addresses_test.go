package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAddresses_AppendIfMissing(t *testing.T) {
	addr1, _ := sdk.AccAddressFromBech32("cosmos16cx5gezkp79wkeynt9vduqrz55gdq3mtj4cmuc")
	addr2, _ := sdk.AccAddressFromBech32("cosmos1vt9vnyhukw65vvqxzp0vdnvddqlc9k88x2ajrm")

	tests := []struct {
		name             string
		addresses        types.Addresses
		address          sdk.AccAddress
		shouldBeAppended bool
	}{
		{
			name:             "Existing address is not appended",
			addresses:        types.Addresses{addr1, addr2},
			address:          addr1,
			shouldBeAppended: false,
		},
		{
			name:             "New address is appended into existing list",
			addresses:        types.Addresses{addr1},
			address:          addr2,
			shouldBeAppended: true,
		},
		{
			name:             "New address is appended into empty list",
			addresses:        types.Addresses{},
			address:          addr1,
			shouldBeAppended: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, appended := test.addresses.AppendIfMissing(test.address)
			require.Equal(t, test.shouldBeAppended, appended)
			require.Contains(t, result, test.address)
		})
	}
}

func TestAddresses_RemoveIfExisting(t *testing.T) {
	addr1, _ := sdk.AccAddressFromBech32("cosmos16cx5gezkp79wkeynt9vduqrz55gdq3mtj4cmuc")
	addr2, _ := sdk.AccAddressFromBech32("cosmos1vt9vnyhukw65vvqxzp0vdnvddqlc9k88x2ajrm")
	addr3, _ := sdk.AccAddressFromBech32("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc")

	tests := []struct {
		name            string
		addresses       types.Addresses
		address         sdk.AccAddress
		shouldBeRemoved bool
	}{
		{
			name:            "Cannot remove from empty list",
			addresses:       types.Addresses{},
			address:         addr1,
			shouldBeRemoved: false,
		},
		{
			name:            "Not found address is not removed",
			addresses:       types.Addresses{addr1, addr2},
			address:         addr3,
			shouldBeRemoved: false,
		},
		{
			name:            "Starting address is removed properly",
			addresses:       types.Addresses{addr1, addr2, addr3},
			address:         addr1,
			shouldBeRemoved: true,
		},
		{
			name:            "Middle address is removed properly",
			addresses:       types.Addresses{addr1, addr2, addr3},
			address:         addr2,
			shouldBeRemoved: true,
		},
		{
			name:            "Ending address is removed properly",
			addresses:       types.Addresses{addr1, addr2, addr3},
			address:         addr3,
			shouldBeRemoved: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, removed := test.addresses.RemoveIfExisting(test.address)
			require.Equal(t, test.shouldBeRemoved, removed)
			require.NotContains(t, result, test.address)
		})
	}
}

func TestAddresses_IndexOf(t *testing.T) {
	addr1, _ := sdk.AccAddressFromBech32("cosmos16cx5gezkp79wkeynt9vduqrz55gdq3mtj4cmuc")
	addr2, _ := sdk.AccAddressFromBech32("cosmos1vt9vnyhukw65vvqxzp0vdnvddqlc9k88x2ajrm")
	addr3, _ := sdk.AccAddressFromBech32("cosmos18q5k63dkyazl88hzvcyx26lqas7al62hqaxlyc")

	require.Equal(t, -1, types.Addresses{addr1}.IndexOf(addr2))
	require.Equal(t, -1, types.Addresses{}.IndexOf(addr1))
	require.Equal(t, 0, types.Addresses{addr1, addr2, addr3}.IndexOf(addr1))
	require.Equal(t, 1, types.Addresses{addr1, addr2, addr3}.IndexOf(addr2))
	require.Equal(t, 2, types.Addresses{addr1, addr2, addr3}.IndexOf(addr3))
}

func TestAddresses_Contains(t *testing.T) {
	addr1, _ := sdk.AccAddressFromBech32("cosmos16cx5gezkp79wkeynt9vduqrz55gdq3mtj4cmuc")
	addr2, _ := sdk.AccAddressFromBech32("cosmos1vt9vnyhukw65vvqxzp0vdnvddqlc9k88x2ajrm")

	require.False(t, types.Addresses{}.Contains(addr1))
	require.False(t, types.Addresses{addr1}.Contains(addr2))
	require.True(t, types.Addresses{addr1, addr2}.Contains(addr1))
}

func TestAddresses_Empty(t *testing.T) {
	addr1, _ := sdk.AccAddressFromBech32("cosmos16cx5gezkp79wkeynt9vduqrz55gdq3mtj4cmuc")

	require.True(t, types.Addresses{}.Empty())
	require.False(t, types.Addresses{addr1}.Empty())
}
