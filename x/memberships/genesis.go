package memberships

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState state at genesis
type GenesisState struct {
	Minters     Minters            `json:"minters"`     // List of users allowed to sign a membership giving tx
	Memberships []types.Membership `json:"memberships"` // List of all the existing memberships
}

// InitGenesis sets membership information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, genState GenesisState) {
	// Import the minters
	for _, minter := range genState.Minters {
		keeper.AddTrustedMinter(ctx, minter)
	}

	// Import the memberships
	for _, membership := range genState.Memberships {
		_, err := keeper.AssignMembership(ctx, membership.Owner, membership.MembershipType)
		if err != nil {
			panic(err)
		}
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Minters:     keeper.GetTrustedMinters(ctx),
		Memberships: keeper.GetMembershipsSet(ctx),
	}
}
