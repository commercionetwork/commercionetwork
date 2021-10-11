package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
)

const (
	eventNewPosition       = "new_position"
	eventBurnCCC           = "burned_ccc"
	eventSetConversionRate = "new_conversion_rate"
	eventSetFreezePeriod   = "new_freeze_period"
)

type Keeper struct {
	cdc           codec.Marshaler
	storeKey      sdk.StoreKey
	memKey        sdk.StoreKey
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	govKeeper     government.Keeper
}

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	bankKeeper bank.Keeper,
	accountKeeper auth.AccountKeeper,
	govKeeper government.Keeper,
) *Keeper {
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		govKeeper:     govKeeper,
	}
}
