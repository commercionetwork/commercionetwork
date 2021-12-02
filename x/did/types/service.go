package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (s ServiceNew) Validate() error {
	if s.ID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"id\" is required")
	}

	if s.Type == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"type\" is required")
	}

	if s.ServiceEndpoint == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"serviceEndpoint\" is required")
	}

	return nil
}
