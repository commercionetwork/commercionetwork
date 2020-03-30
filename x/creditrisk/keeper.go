package creditrisk

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/commercionetwork/commercionetwork/x/creditrisk/types"
)

// Keeper
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	supplyKeeper supply.Keeper
}

// NewKeeper creates new instances of the module Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, supplyKeeper supply.Keeper) Keeper {
	// ensure creditrisk module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		supplyKeeper: supplyKeeper,
	}
}

// GetPoolFunds returns the funds currently present inside the rewards pool
func (k Keeper) GetPoolFunds(ctx sdk.Context) sdk.Coins { return k.getModuleAccount(ctx).GetCoins() }

// getModuleAccount returns the module account for the accreditations module
func (k Keeper) getModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	modAcc := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	if modAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	return modAcc
}

func (k Keeper) setModuleAccount(ctx sdk.Context, modAcc exported.ModuleAccountI) {
	k.supplyKeeper.SetModuleAccount(ctx, modAcc)
}
