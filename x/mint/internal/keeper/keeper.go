package keeper

import (
	"fmt"

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

func (keeper Keeper) SetCreditsDenom(ctx sdk.Context, denom string) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.CreditsDenomStoreKey), []byte(denom))
}

func (keeper Keeper) GetCreditsDenom(ctx sdk.Context) string {
	store := ctx.KVStore(keeper.StoreKey)
	return string(store.Get([]byte(types.CreditsDenomStoreKey)))
}

func (keeper Keeper) getCDPkey(address sdk.AccAddress) []byte {
	return []byte(types.CDPStoreKey + address.String())
}

func (keeper Keeper) AddCDP(ctx sdk.Context, cdp types.CDP) {
	var cdps types.CDPs
	store := ctx.KVStore(keeper.StoreKey)
	cdpsBz := store.Get(keeper.getCDPkey(cdp.Owner))
	keeper.cdc.MustUnmarshalBinaryBare(cdpsBz, &cdps)
	cdps, found := cdps.AppendIfMissing(cdp)
	if !found {
		store.Set(keeper.getCDPkey(cdp.Owner), keeper.cdc.MustMarshalBinaryBare(cdps))
	}
}

func (keeper Keeper) DeleteCDP(ctx, cdp types.CDP) bool {
	return true
}

// OpenCDP subtract the given token amount from user's wallet and deposit it into the liquidity pool then,
// send him the corresponding commercio cash credits amount.
// Errors occurs if:
// 1) deposited tokens haven't been priced yet, or are negatives or invalid
// 2) signer's funds are not enough
// 3)
func (keeper Keeper) OpenCDP(ctx sdk.Context, cdpRequest types.CDPRequest) (sdk.Coins, sdk.Error) {

	if !cdpRequest.DepositedAmount.IsValid() || cdpRequest.DepositedAmount.IsAnyNegative() {
		return nil, sdk.ErrInvalidCoins(cdpRequest.DepositedAmount.String())
	}

	store := ctx.KVStore(keeper.StoreKey)
	fiatValue := sdk.NewInt(0)

	//Check if all tokens in deposit amount have a price and calculate the total FIAT value of them
	for _, token := range cdpRequest.DepositedAmount {
		assetPrice, found := keeper.PriceFeedKeeper.GetCurrentPrice(ctx, token.Denom)
		if found == false {
			return nil, sdk.ErrInvalidCoins(fmt.Sprintf("no current price for given token: %s", token.Denom))
		}
		fiatValue = fiatValue.Add(token.Amount.Mul(assetPrice.Price.RoundInt()))
	}

	//Subtract the given deposit amount from user's wallet
	_, err := keeper.BankKeeper.SubtractCoins(ctx, cdpRequest.Signer, cdpRequest.DepositedAmount)
	if err != nil {
		return nil, err
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
	credits, err = keeper.BankKeeper.AddCoins(ctx, cdpRequest.Signer, credits)
	if err != nil {
		return nil, err
	}

	cdp := types.NewCDP(cdpRequest, credits)
	keeper.AddCDP(ctx, cdp)

	return credits, nil
}

//Withdraw token subtract the given credits amount from user's wallet and send to it the corresponding value in commercio tokens.
//If user's wallet hasn't enought credits, an error will occur.
func (keeper Keeper) CloseCDP(ctx sdk.Context, user sdk.AccAddress, credits sdk.Coins) (sdk.Coins, error) {
	return nil, nil
}
