package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	funderAddr, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	validAmount   = sdk.NewCoins(sdk.Coin{
		Denom:  BondDenom,
		Amount: sdk.NewInt(100),
	})

	ValidMsgIncrementBlockRewardsPool = *NewMsgIncrementBlockRewardsPool(
		funderAddr.String(),
		validAmount,
	)
)

var (
	validDistrEpochIdentifier = EpochDay

	validEarnRate   = sdk.NewDecWithPrec(5, 1)
	invalidEarnRate = sdk.NewDecWithPrec(-5, 1)

	validParams = DefaultParams()

	governmentAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

	ValidMsgSetParams = *NewMsgSetParams(
		governmentAddress.String(),
		validDistrEpochIdentifier,
		validEarnRate,
	)
)
