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
	}
	return append(addresses, address), true
}

// RemoveIfExisting returns a new Addresses instance that does not contain the
// given address.
func (addresses Addresses) RemoveIfExisting(address sdk.AccAddress) (Addresses, bool) {
	indexOf := addresses.IndexOf(address)
	if indexOf > -1 {
		return append(addresses[:indexOf], addresses[indexOf+1:]...), true
	}
	return addresses, false
}

// IndexOf returns the index of the given address inside the addresses array,
// or -1 if such an address was not found
func (addresses Addresses) IndexOf(address sdk.Address) int {
	for i, a := range addresses {
		if a.Equals(address) {
			return i
		}
	}
	return -1
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

// Empty returns true if this slice does not contain any address
func (addresses Addresses) Empty() bool {
	return len(addresses) == 0
}
