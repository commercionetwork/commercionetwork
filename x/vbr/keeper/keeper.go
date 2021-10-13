package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	distKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	govKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc      	codec.Marshaler
		storeKey	sdk.StoreKey
		//memKey   sdk.StoreKey
		// this line is used by starport scaffolding # ibc/keeper/attribute
		distKeeper  distKeeper.Keeper
		bankKeeper	bankKeeper.Keeper
		accountKeeper accountKeeper.AccountKeeper
		govKeeper   govKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey sdk.StoreKey,
	memKey sdk.StoreKey,
	// this line is used by starport scaffolding # ibc/keeper/parameter
	distKeeper   distKeeper.Keeper,
	bankKeeper	bankKeeper.Keeper,
	accountKeeper accountKeeper.AccountKeeper,
	govKeeper    govKeeper.Keeper,

) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		//memKey:   memKey,
		// this line is used by starport scaffolding # ibc/keeper/return
		distKeeper: distKeeper,
		bankKeeper: bankKeeper,
		accountKeeper: accountKeeper,
		govKeeper: govKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}


// -------------
// --- Pool
// -------------

// SetTotalRewardPool allows to set the value of the total rewards pool that has left
func (k Keeper) SetTotalRewardPool(ctx sdk.Context, updatedPool sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	if !updatedPool.Empty() {
		store.Set([]byte(types.PoolStoreKey), k.cdc.MustMarshalBinaryBare(&updatedPool))
	} else {
		store.Delete([]byte(types.PoolStoreKey))
	}
}

// GetTotalRewardPool returns the current total rewards pool amount
func (k Keeper) GetTotalRewardPool(ctx sdk.Context) sdk.DecCoins {
	macc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	mcoins := macc.GetCoins()

	return sdk.NewDecCoinsFromCoins(mcoins...)
}

// SetRewardRate store the vbr reward rate.
func (k Keeper) SetRewardRateKeeper(ctx sdk.Context, rate sdk.Dec) error {
	if err := types.ValidateRewardRate(rate); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.RewardRateKey), k.cdc.MustMarshalBinaryBare(rate))
	return nil
}

// SetAutomaticWithdraw store the automatic withdraw flag.
func (k Keeper) SetAutomaticWithdrawKeeper(ctx sdk.Context, autoW bool) error {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.AutomaticWithdraw), k.cdc.MustMarshalBinaryBare(autoW))
	return nil
}
