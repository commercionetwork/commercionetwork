package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const commonEvents = "commerciotx"

func EmitCommonEvents(ctx sdk.Context, sender sdk.AccAddress) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		commonEvents,
		sdk.NewAttribute("account", sender.String())))
}
