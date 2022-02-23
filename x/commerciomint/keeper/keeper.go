package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramTypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
)

const (
	eventNewPosition = "new_position"
	eventBurnCCC     = "burned_ccc"
	eventSetParams   = "new_params"
)

type Keeper struct {
	cdc           codec.Marshaler
	storeKey      sdk.StoreKey
	memKey        sdk.StoreKey
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	govKeeper     governmentKeeper.Keeper
	paramSpace    paramTypes.Subspace
}

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	bankKeeper bank.Keeper,
	accountKeeper auth.AccountKeeper,
	govKeeper governmentKeeper.Keeper,
	paramSpace paramTypes.Subspace,
) *Keeper {
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		govKeeper:     govKeeper,
		paramSpace:    paramSpace,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
