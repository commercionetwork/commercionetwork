package utilities

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------------------
// --- Utility functions
// ----------------------------------

func AppendStringIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func AppendAddressIfMissing(slice []sdk.AccAddress, i sdk.AccAddress) []sdk.AccAddress {
	for _, ele := range slice {
		if ele.Equals(i) {
			return slice
		}
	}
	return append(slice, i)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func DidInSlice(a sdk.AccAddress, list []sdk.AccAddress) bool {
	for _, b := range list {
		if b.Equals(a) {
			return true
		}
	}
	return false
}
