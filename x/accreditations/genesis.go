package accreditations

import (
	"errors"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - accreditations genesis state
type GenesisState struct {
	LiquidityPoolAmount     sdk.Coins         `json:"liquidity_pool_amount"`
	Invites                 []types.Invite    `json:"invites"`
	TrustedServiceProviders ctypes.Addresses  `json:"trusted_service_providers"`
	Credentials             types.Credentials `json:"credentials"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// Set the liquidity pool
	if data.LiquidityPoolAmount != nil {
		keeper.SetPoolFunds(ctx, data.LiquidityPoolAmount)
	}

	// Import the signers
	for _, signer := range data.TrustedServiceProviders {
		keeper.AddTrustedServiceProvider(ctx, signer)
	}

	// Import all the invites
	for _, invite := range data.Invites {
		keeper.SaveInvite(ctx, invite)
	}

	// Import the credentials
	for _, credential := range data.Credentials {
		keeper.SaveCredential(ctx, credential)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		LiquidityPoolAmount:     keeper.GetPoolFunds(ctx),
		Invites:                 keeper.GetInvites(ctx),
		TrustedServiceProviders: keeper.GetTrustedServiceProviders(ctx),
		Credentials:             keeper.GetCredentials(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.LiquidityPoolAmount.IsAnyNegative() {
		return errors.New("liquidity pool amount cannot contain negative values")
	}

	return nil
}
