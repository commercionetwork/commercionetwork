package pricefeed

import (
	"fmt"
	"strings"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - docs genesis state
type GenesisState struct {
	Oracles   ctypes.Addresses `json:"oracles"`
	Assets    ctypes.Strings   `json:"assets"`
	RawPrices RawPrices        `json:"raw_prices"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, genState GenesisState) {
	for _, oracle := range genState.Oracles {
		keeper.AddOracle(ctx, oracle)
	}

	for _, asset := range genState.Assets {
		keeper.AddAsset(ctx, asset)
	}

	for _, rawPrice := range genState.RawPrices {
		if err := keeper.AddRawPrice(ctx, rawPrice.Oracle, rawPrice.Price); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Oracles:   keeper.GetOracles(ctx),
		Assets:    keeper.GetAssets(ctx),
		RawPrices: keeper.GetRawPrices(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(state GenesisState) error {
	for _, asset := range state.Assets {
		if len(strings.TrimSpace(asset)) == 0 {
			return fmt.Errorf("%s, is empty", asset)
		}
	}
	for _, oracle := range state.Oracles {
		if oracle.Empty() {
			return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Found Empty oracle address")
		}
	}
	return nil
}
