package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Utility type
type Minters []sdk.AccAddress

// Contains returns true iff the given minter is contained inside the minters list
func (minters Minters) Contains(minter sdk.AccAddress) bool {
	for _, mint := range minters {
		if mint.Equals(minter) {
			return true
		}
	}

	return false
}
