package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Identity struct {
	Owner       sdk.AccAddress `json:"owner"`
	DidDocument string         `json:"did_document"`
}
