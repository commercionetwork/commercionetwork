package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Funder struct {
	Address sdk.AccAddress `json:"address"`
}

type Funders []Funder

func (funders Funders) AppendFunderIfMissing(funder Funder) Funders {
	for _, ele := range funders {
		if ele.Equals(funder) {
			return funders
		}
	}
	return append(funders, funder)
}

func (funders Funders) FindFunder(funder Funder) Funder {
	var foundFunder Funder
	for _, ele := range funders {
		if ele.Equals(funder) {
			foundFunder = ele
		}
	}
	return foundFunder
}

func (funder Funder) Equals(obj Funder) bool {
	return funder.Address.Equals(obj.Address)
}
