package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BlockRewardsPool struct {
	Funds sdk.DecCoins `json:"funds"`
}

func InitBlockRewardsPool() BlockRewardsPool {
	return BlockRewardsPool{
		Funds: sdk.DecCoins{},
	}
}

func (brp BlockRewardsPool) ValidateGenesis() error {
	if brp.Funds.IsAnyNegative() {
		return fmt.Errorf("negative Funds in block reward pool, is %v", brp.Funds)
	}

	return nil
}
