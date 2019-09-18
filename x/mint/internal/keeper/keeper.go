package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	pricefeed "github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	BankKeeper      bank.BaseKeeper
	PriceFeedKeeper pricefeed.Keeper

	cdc *codec.Codec
}

func NewKeeper(bk bank.BaseKeeper, pk pricefeed.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		BankKeeper:      bk,
		PriceFeedKeeper: pk,
		cdc:             cdc,
	}
}

// DepositToken subtract the given token amount from user's wallet, and send him the corresponding commercio credits
// If user's wallet hasn't got enough funds or if there will be some problems while adding credits to the wallet, an error will occur
func (keeper Keeper) DepositToken(ctx sdk.Context, user sdk.AccAddress, token sdk.Coins) (sdk.Coins, error) {

	//Subtract the given token amount from user's wallet
	_, err := keeper.BankKeeper.SubtractCoins(ctx, user, token)
	if err != nil {
		return nil, err
	}
	//get token's current price
	tokenPrice := keeper.PriceFeedKeeper.GetPrice(ctx, types.DefaultBondDenom)

	//get the token value = tokens amount * token price
	tokenValue := tokenPrice.Mul(token.AmountOf(types.DefaultBondDenom))

	//get credits' current price
	creditsPrice := keeper.PriceFeedKeeper.GetPrice(ctx, types.DefaultCreditsDenom)

	//get credits' amount = token value / credits price
	creditsAmount := tokenValue.Quo(creditsPrice)

	//add credits to users wallet
	credits := sdk.NewCoins(sdk.NewCoin(types.DefaultCreditsDenom, creditsAmount))
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
	creditsPrice := keeper.PriceFeedKeeper.GetPrice(ctx, types.DefaultCreditsDenom)

	//get the credits' value = credits amount * credits price
	creditsValue := creditsPrice.Mul(credits.AmountOf(types.DefaultCreditsDenom))

	//get token's current price
	tokenPrice := keeper.PriceFeedKeeper.GetPrice(ctx, types.DefaultBondDenom)

	//get tokens' amount = credits value / token price
	tokensAmount := creditsValue.Quo(tokenPrice)
	tokens := sdk.NewCoins(sdk.NewCoin(types.DefaultBondDenom, tokensAmount))
	credits, err = keeper.BankKeeper.AddCoins(ctx, user, tokens)
	if err != nil {
		return nil, err
	}

	return tokens, err
}
