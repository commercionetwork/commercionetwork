package customstaking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

var MinimumDeposit = sdk.NewCoin("ucommercio", sdk.TokensFromConsensusPower(50000))

func NewHandler(k keeper.Keeper, stakingHandler sdk.Handler) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateValidator:
			return handleMsgCreateValidator(ctx, stakingHandler, msg)

		default:
			return stakingHandler(ctx, msg)
		}
	}
}

// These functions assume everything has been authenticated,
// now we just perform action and save

func handleMsgCreateValidator(ctx sdk.Context, handler sdk.Handler, msg types.MsgCreateValidator) (*sdk.Result, error) {
	if msg.Value.IsLT(MinimumDeposit) {
		return nil, ErrMinimumStake
	}

	return handler(ctx, msg)
}
