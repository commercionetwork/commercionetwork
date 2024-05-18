package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrUnauthorized		 = sdkerrors.Register(ModuleName, 2, "unauthorized")
	ErrBadMessage        = sdkerrors.Register(ModuleName, 3, "bad message")
	ErrContractError     = sdkerrors.Register(ModuleName, 4, "contract error")
)
