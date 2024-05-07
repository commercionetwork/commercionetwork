package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

const (
	eventSavedDocument = "new_saved_document"
	eventSavedReceipt  = "new_saved_receipt"
)

type (
	Keeper struct {
		cdc      codec.Codec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
