package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Accreditation struct {
	Accrediter sdk.AccAddress `json:"accrediter"`
	User       sdk.AccAddress `json:"user"`
}
