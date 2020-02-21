package memberships

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(_ sdk.Context, _ Keeper, _ auth.AccountKeeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
