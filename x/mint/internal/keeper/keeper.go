package keeper

import (
	pricefeed "github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	StoreKey sdk.StoreKey

	BankKeeper      bank.BaseKeeper
	PriceFeedKeeper pricefeed.Keeper

	cdc *codec.Codec
}

func NewKeeper(storekey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey: storekey,
		cdc:      cdc,
	}
}

func (keeper Keeper) GetCredits(user sdk.AccAddress, tokenAmt sdk.Coins) sdk.Coins {

}
