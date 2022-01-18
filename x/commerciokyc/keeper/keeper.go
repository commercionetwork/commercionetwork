package keeper

import (
	"fmt"
	"time"

	mint "github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

const (
	SecondsPerYear time.Duration = time.Hour * 24 * 365
)

var membershipCosts = map[string]int64{
	types.MembershipTypeGreen:  1,
	types.MembershipTypeBronze: 25,
	types.MembershipTypeSilver: 250,
	types.MembershipTypeGold:   2500,
	types.MembershipTypeBlack:  50000,
}

type Keeper struct {
	Cdc           codec.Marshaler
	StoreKey      sdk.StoreKey
	memKey        sdk.StoreKey
	bankKeeper    bank.Keeper
	GovKeeper     government.Keeper
	accountKeeper auth.AccountKeeper
	MintKeeper    mint.Keeper
	paramSpace    paramtypes.Subspace
}

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	bankKeeper bank.Keeper,
	govKeeper government.Keeper,
	accountKeeper auth.AccountKeeper,
	mintKeeper mint.Keeper,
	paramSpace paramtypes.Subspace,
) *Keeper {
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		Cdc:           cdc,
		StoreKey:      storeKey,
		memKey:        memKey,
		bankKeeper:    bankKeeper,
		GovKeeper:     govKeeper,
		accountKeeper: accountKeeper,
		MintKeeper:    mintKeeper,
		paramSpace:    paramSpace,
	}
}

func (k Keeper) GetAccountKeeper() auth.AccountKeeper {
	return k.accountKeeper
}

func (k Keeper) GetBankKeeper() bank.Keeper {
	return k.bankKeeper
}
