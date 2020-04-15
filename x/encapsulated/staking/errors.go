package customstaking

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var ErrMinimumStake = sdkerrors.Register(staking.ModuleName, 46, "delegation amount is not enough")
