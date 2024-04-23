package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnauthorized		 = errorsmod.Register(ModuleName, 2, "unauthorized")
	ErrBadMessage        = errorsmod.Register(ModuleName, 3, "bad message")
	ErrContractError     = errorsmod.Register(ModuleName, 4, "contract error")
)
