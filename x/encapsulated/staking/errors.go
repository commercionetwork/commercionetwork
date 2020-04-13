package customstaking

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var ErrMinimumDeposit = sdkerrors.Register(staking.ModuleName, 46, "deposit is not enough")
