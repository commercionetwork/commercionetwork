package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Accreditation struct {
	Accrediter sdk.AccAddress `json:"accrediter"`
	User       sdk.AccAddress `json:"user"`
	Rewarded   bool           `json:"rewarded"` // Tells if the accrediter has already been rewarded for this accreditation
}
