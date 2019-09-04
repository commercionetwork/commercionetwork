package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Identity contains the data of an identity that can be specified inside the genesis and are
// exported during the current state exportation.
type Identity struct {
	Owner       sdk.AccAddress `json:"owner"`
	DidDocument string         `json:"did_document"`
}
