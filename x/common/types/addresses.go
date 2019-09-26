package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Addresses is an alias for a list of sdk.AccAddress that
// enables custom operations
type Addresses []sdk.AccAddress

// AppendIfMissing returns a new Addresses instance containing the given
// address if it wasn't already present
func (addresses Addresses) AppendIfMissing(address sdk.AccAddress) (Addresses, bool) {
	if addresses.Contains(address) {
		return addresses, false
	} else {
		return append(addresses, address), true
	}
}

// Contains returns true iff the addresses list contains the given address
func (addresses Addresses) Contains(address sdk.Address) bool {
	for _, ele := range addresses {
		if ele.Equals(address) {
			return true
		}
	}
	return false
}

func (addresses Addresses) Empty() bool {
	return len(addresses) == 0
}
