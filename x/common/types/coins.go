package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Coins []sdk.Coin

//convert []sdk.Coin to []*sdk.Coin
func (cs Coins) ConvertToCoinsPointers() []*sdk.Coin {
	res := []*sdk.Coin{}

	for _, c := range cs {
		res = append(res, &c)
	}

	return res
}
/*
type CoinsPointers []*sdk.Coin

//returns sdk.Coins
func (cp CoinsPointers) ConvertFromCoinsPointers() sdk.Coins {
	res := sdk.Coins{}

	for _, c := range cp {
		res = append(res, *c)
	}

	return res
}
*/