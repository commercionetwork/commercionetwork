package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
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

// DepositToken subtract the given token amount from user's wallet and deposit it into the liquidity pool then,
// send him the corresponding commercio cash credits amount.
// If user's wallet hasn't got enough funds or if there will be some problems while adding credits to the wallet,
// an error will occur.
func (keeper Keeper) DepositToken(ctx sdk.Context, user sdk.AccAddress, depositAmount sdk.Coins) (sdk.Coins, sdk.Error) {

	store := ctx.KVStore(keeper.StoreKey)

	//get deposit amount's current price
	assetPrice, found := keeper.PriceFeedKeeper.GetCurrentPrice(ctx, keeper.StakingKeeper.BondDenom(ctx))
	if found == false {
		return nil, sdk.ErrInvalidCoins("no current price for given tokens")
	}

	//get credits' current price
	creditsPrice, found := keeper.PriceFeedKeeper.GetCurrentPrice(ctx, types.CreditsDenom)
	if found == false {
		return nil, sdk.ErrInvalidCoins("no current price for given tokens")
	}

	//Subtract the given deposit amount from user's wallet
	_, err := keeper.BankKeeper.SubtractCoins(ctx, user, depositAmount)
	if err != nil {
		return nil, err
	}

	poolBz := store.Get([]byte(types.LiquidityPoolKey))
	var liquidityPool = sdk.DecCoins{}
	keeper.cdc.MustUnmarshalBinaryBare(poolBz, &liquidityPool)

	//get the depositAmount value = deposit amount * asset price
	da := sdk.NewDecCoins(depositAmount)
	//depositing the amount to pool
	liquidityPool = liquidityPool.Add(da)
	store.Set([]byte(types.LiquidityPoolKey), keeper.cdc.MustMarshalBinaryBare(liquidityPool))

	//calculating deposit amount value
	depositAmountValue := da.MulDec(assetPrice.Price)

	//get credits' amount = depositAmount value / credits price
	creditsAmount := depositAmountValue.QuoDec(creditsPrice.Price)

	//add credits to users wallet
	credits := sdk.NewCoins(sdk.NewCoin(types.CreditsDenom, creditsAmount.TruncateDecimal()))
	credits, err = keeper.BankKeeper.AddCoins(ctx, user, credits)
	if err != nil {
		return nil, err
	}

	return credits, nil
}

//Withdraw token subtract the given credits amount from user's wallet and send to it the corresponding value in commercio tokens.
//If user's wallet hasn't enought credits, an error will occur.
func (keeper Keeper) WithdrawToken(ctx sdk.Context, user sdk.AccAddress, credits sdk.Coins) (sdk.Coins, error) {

	//Subtract credits amount from user's wallet
	_, err := keeper.BankKeeper.SubtractCoins(ctx, user, credits)
	if err != nil {
		return nil, err
	}

	//get credit's current price
	creditsPrice := keeper.PriceFeedKeeper.GetPrice(ctx, ctypes.DefaultCreditsDenom)

	//get the credits' value = credits amount * credits price
	creditsValue := creditsPrice.Mul(credits.AmountOf(ctypes.DefaultCreditsDenom))

	//get token's current price
	tokenPrice := keeper.PriceFeedKeeper.GetPrice(ctx, ctypes.DefaultBondDenom)

	//get tokens' amount = credits value / token price
	tokensAmount := creditsValue.Quo(tokenPrice)
	tokens := sdk.NewCoins(sdk.NewCoin(ctypes.DefaultBondDenom, tokensAmount))
	credits, err = keeper.BankKeeper.AddCoins(ctx, user, tokens)
	if err != nil {
		return nil, err
	}

	return tokens, err
}
