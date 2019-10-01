package mint

import (
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type UsersCDP struct {
	Owner sdk.AccAddress `json:"Owner"`
	CDPs  types.CDPs     `json:"cdpS"`
}

// GenesisState - docs genesis state
type GenesisState struct {
	Users []UsersCDP
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	users := keeper.GetUsersSet(ctx)

	var usersCDPs []UsersCDP
	for _, user := range users {
		CDPs := keeper.GetCDPs(ctx, user)
		userCDPs := UsersCDP{
			Owner: user,
			CDPs:  CDPs,
		}
		usersCDPs = append(usersCDPs, userCDPs)
	}

	return GenesisState{
		Users: usersCDPs,
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(state GenesisState) error {
	for _, userCDPs := range state.Users {
		if userCDPs.Owner.Empty() {
			return sdk.ErrInvalidAddress(userCDPs.Owner.String())
		}
		for _, cdp := range userCDPs.CDPs {
			err := cdp.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
