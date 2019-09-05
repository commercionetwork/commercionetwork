package accreditation

import (
	"github.com/commercionetwork/commercionetwork/x/accreditation/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - docs genesis state
type GenesisState struct {
	Accreditations     []types.Accreditation `json:"users_data"`
	TrustworthySigners []sdk.AccAddress      `json:"trustworthy_signers"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// Import the signers
	for _, signer := range data.TrustworthySigners {
		keeper.AddTrustworthySigner(ctx, signer)
	}

	// Import all the accreditations
	for _, accreditation := range data.Accreditations {
		keeper.SetAccrediter(ctx, accreditation.Accrediter, accreditation.User)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Accreditations:     keeper.GetAccreditations(ctx),
		TrustworthySigners: keeper.GetTrustworthySigners(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}
