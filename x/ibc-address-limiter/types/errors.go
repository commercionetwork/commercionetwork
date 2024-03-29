package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUnauthorized		 = sdkerrors.Register(ModuleName, 2, "unauthorized")
	ErrBadMessage        = sdkerrors.Register(ModuleName, 3, "bad message")
	ErrContractError     = sdkerrors.Register(ModuleName, 4, "contract error")
)
