package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var funderAddr, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var validAmount = sdk.NewCoins(sdk.Coin{
	Denom:  BondDenom,
	Amount: sdk.NewInt(100),
})

var ValidMsgIncrementBlockRewardsPool = *NewMsgIncrementBlockRewardsPool(
	funderAddr.String(),
	validAmount,
)
