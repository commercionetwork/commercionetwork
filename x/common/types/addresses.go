package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Addresses is an alias for a list of sdk.AccAddress that
// enables custom operations
type Addresses []sdk.AccAddress

// AppendIfMissing returns a new Addresses instance containing the given
// address if it wasn't already present
func (addresses Addresses) AppendIfMissing(address sdk.AccAddress) Addresses {
	if addresses.Contains(address) {
		return addresses
	} else {
		return append(addresses, address)
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
