package vbr

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/vbr/keeper"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

type GenesisState struct {
	PoolAmount sdk.DecCoins `json:"pool_amount"`
	RewardRate sdk.Dec      `json:"reward_rate"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets the initial Block Reward Pool amount for genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data GenesisState) {
	// Set the reward pool - Should never be nil as its validated inside the ValidateGenesis method
	keeper.SetTotalRewardPool(ctx, data.PoolAmount)

	moduleAcc := keeper.VbrAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if moduleAcc.GetCoins().Empty() {
		amount, _ := data.PoolAmount.TruncateDecimal()
		err := keeper.MintVBRTokens(ctx, sdk.NewCoins(amount...))
		if err != nil {
			panic(err) // could not mint tokens on chain start, fatal!
		}
	}
	keeper.SetRewardRate(ctx, data.RewardRate)

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	return GenesisState{
		PoolAmount: keeper.GetTotalRewardPool(ctx),
		RewardRate: keeper.GetRewardRate(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.PoolAmount == nil || data.PoolAmount.Empty() {
		return errors.New("validator block reward pool cannot be empty")
	}

	if !data.PoolAmount.IsValid() {
		return fmt.Errorf("invalid validator block reward pool: %s", data.PoolAmount.String())
	}

	return types.ValidateRewardRate(data.RewardRate)
}
