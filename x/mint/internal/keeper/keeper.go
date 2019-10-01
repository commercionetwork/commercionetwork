package keeper

import (
	"fmt"
	"strings"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type Keeper struct {
	StoreKey        sdk.StoreKey
	BankKeeper      bank.BaseKeeper
	PriceFeedKeeper pricefeed.Keeper
	StakingKeeper   staking.Keeper

	cdc *codec.Codec
}

func NewKeeper(sk sdk.StoreKey, bk bank.BaseKeeper, pk pricefeed.Keeper, stk staking.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey:        sk,
		BankKeeper:      bk,
		PriceFeedKeeper: pk,
		StakingKeeper:   stk,
		cdc:             cdc,
	}
}

func (keeper Keeper) SetCreditsDenom(ctx sdk.Context, den string) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.CreditsDenomStoreKey), []byte(den))
}

func (keeper Keeper) GetCreditsDenom(ctx sdk.Context) string {
	store := ctx.KVStore(keeper.StoreKey)
	return string(store.Get([]byte(types.CreditsDenomStoreKey)))
}

func (keeper Keeper) getCDPkey(address sdk.AccAddress) []byte {
	return []byte(types.CDPStoreKey + address.String())
}

// GetUsersSet returns the list of all the users that open at least one CDP.
func (keeper Keeper) GetUsersSet(ctx sdk.Context) ctypes.Addresses {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CDPStoreKey))
	var users = ctypes.Addresses{}
	for ; iterator.Valid(); iterator.Next() {
		addressStr := strings.ReplaceAll(string(iterator.Key()), types.CDPStoreKey, "")
		address, _ := sdk.AccAddressFromBech32(addressStr)
		users, _ = users.AppendIfMissing(address)
	}
	return users
}

//add a CDPs to user's CDPs list
func (keeper Keeper) AddCDP(ctx sdk.Context, cdp types.CDP) {
	var cdpS types.CDPs
	store := ctx.KVStore(keeper.StoreKey)
	cdpSBz := store.Get(keeper.getCDPkey(cdp.Owner))
	keeper.cdc.MustUnmarshalBinaryBare(cdpSBz, &cdpS)
	cdpS, found := cdpS.AppendIfMissing(cdp)
	if !found {
		store.Set(keeper.getCDPkey(cdp.Owner), keeper.cdc.MustMarshalBinaryBare(cdpS))
	}
}

func (keeper Keeper) GetCDPs(ctx sdk.Context, owner sdk.AccAddress) types.CDPs {
	var cdpS types.CDPs
	store := ctx.KVStore(keeper.StoreKey)
	cdpSBz := store.Get(keeper.getCDPkey(owner))
	keeper.cdc.MustUnmarshalBinaryBare(cdpSBz, &cdpS)
	return cdpS
}

func (keeper Keeper) GetCDP(ctx sdk.Context, owner sdk.AccAddress, timestamp string) *types.CDP {
	cdpS := keeper.GetCDPs(ctx, owner)
	cdp, found := cdpS.GetCdpFromTimestamp(timestamp)
	if !found {
		return nil
	}
	return cdp
}

func (keeper Keeper) DeleteCDP(ctx sdk.Context, owner sdk.AccAddress, timestamp string) bool {
	var cdpS types.CDPs
	store := ctx.KVStore(keeper.StoreKey)
	cdpSBz := store.Get(keeper.getCDPkey(owner))
	keeper.cdc.MustUnmarshalBinaryBare(cdpSBz, &cdpS)
	cdpS, found := cdpS.RemoveWhenFound(timestamp)
	if found {
		store.Set(keeper.getCDPkey(owner), keeper.cdc.MustMarshalBinaryBare(cdpS))
		return true
	}
	return false
}

// OpenCDP subtract the given token's amount from user's wallet and deposit it into the liquidity pool then,
// sending him the corresponding commercio cash credits amount.
// If all these operations are done correctly, a Collateralized Debt Position is opened.
// Errors occurs if:
// 1) deposited tokens haven't been priced yet, or are negatives or invalid;
// 2) signer's funds are not enough
func (keeper Keeper) OpenCDP(ctx sdk.Context, cdpRequest types.CDPRequest) sdk.Error {

	if !cdpRequest.DepositedAmount.IsValid() || cdpRequest.DepositedAmount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(cdpRequest.DepositedAmount.String())
	}

	store := ctx.KVStore(keeper.StoreKey)
	fiatValue := sdk.NewInt(0)

	//Check if all tokens in deposit amount have a price and calculate the total FIAT value of them
	for _, token := range cdpRequest.DepositedAmount {
		assetPrice, found := keeper.PriceFeedKeeper.GetCurrentPrice(ctx, token.Denom)
		if found == false {
			return sdk.ErrInvalidCoins(fmt.Sprintf("no current price for given token: %s", token.Denom))
		}
		fiatValue = fiatValue.Add(token.Amount.Mul(assetPrice.Price.RoundInt()))
	}

	//Subtract the given deposit amount from user's wallet
	_, err := keeper.BankKeeper.SubtractCoins(ctx, cdpRequest.Signer, cdpRequest.DepositedAmount)
	if err != nil {
		return err
	}

	poolBz := store.Get([]byte(types.LiquidityPoolStoreKey))
	var liquidityPool = sdk.Coins{}
	keeper.cdc.MustUnmarshalBinaryBare(poolBz, &liquidityPool)

	//depositing the amount to the liquidity pool
	liquidityPool = liquidityPool.Add(cdpRequest.DepositedAmount)
	store.Set([]byte(types.LiquidityPoolStoreKey), keeper.cdc.MustMarshalBinaryBare(liquidityPool))

	//get credits' amount = DepositAmount value / credits price (always 1 euro) / 2 is the power of collateral which is 2:1 (comm -> ccc)
	creditsAmount := fiatValue.Quo(sdk.NewInt(2))

	//add credits to users wallet
	credits := sdk.NewCoins(sdk.NewCoin(keeper.GetCreditsDenom(ctx), creditsAmount))
	credits, _ = keeper.BankKeeper.AddCoins(ctx, cdpRequest.Signer, credits)

	//Creating CDP and adding to the store
	cdp := types.NewCDP(cdpRequest, credits)
	keeper.AddCDP(ctx, cdp)

	return nil
}

//CloseCDP subtract the CDP's liquidity amount (commercio cash credits) from user's wallet, after that sends the
//deposited amount back to it. If these two operations ends without errors, the CDP get closed.
//Errors occurs if:
//- cdp doesnt exist
//- subtracting or adding fund to account don't end well
func (keeper Keeper) CloseCDP(ctx sdk.Context, user sdk.AccAddress, timestamp string) sdk.Error {
	cdp := keeper.GetCDP(ctx, user, timestamp)
	if cdp == nil {
		return sdk.ErrInternal("cannot close an inexistent cdp")
	}

	//subtracting liquidity amount from user's wallet
	_, err := keeper.BankKeeper.SubtractCoins(ctx, user, cdp.LiquidityAmount)
	if err != nil {
		return err
	}
	//adding back the deposited amount to user's wallet
	_, err = keeper.BankKeeper.AddCoins(ctx, user, cdp.DepositedAmount)
	if err != nil {
		return err
	}

	_ = keeper.DeleteCDP(ctx, user, timestamp)

	return nil
}
