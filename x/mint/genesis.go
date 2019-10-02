package mint

import (
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - docs genesis state
type GenesisState struct {
	UsersCDPs []types.CDPs
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {

	for _, userCDPs := range data.UsersCDPs {
		for _, cdp := range userCDPs {
			keeper.AddCDP(ctx, cdp)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	users := keeper.GetUsersSet(ctx)

	var usersCDPs = make([]types.CDPs, 0)

	for _, user := range users {
		CDPs := keeper.GetCDPs(ctx, user)
		usersCDPs = append(usersCDPs, CDPs)
	}

	return GenesisState{
		UsersCDPs: usersCDPs,
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(state GenesisState) error {
	for _, userCDPs := range state.UsersCDPs {
		for _, cdp := range userCDPs {
			err := cdp.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
