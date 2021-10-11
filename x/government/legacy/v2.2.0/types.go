package v2_2_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "government"
)

// GenesisState - government genesis state
type GenesisState struct {
	GovernmentAddress sdk.AccAddress `json:"government_address"`
	TumblerAddress    sdk.AccAddress `json:"tumbler_address"`
}
