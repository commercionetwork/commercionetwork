package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const attributeKeySender = "sender"

func EmitCommonEvents(ctx sdk.Context, sender sdk.AccAddress) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(attributeKeySender, sender.String())))
}
