package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

var upgrades = map[string]func(ctx sdk.Context, plan upgrade.Plan){
	/*"testUpgrade": func(ctx sdk.Context, plan upgrade.Plan) {
		ctx.Logger().Info("testUpgrade done!")
	},*/
}
