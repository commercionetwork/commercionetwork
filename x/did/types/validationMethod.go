package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

func (v VerificationMethod) Validate() error {
	if v.ID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"id\" is required")
	}

	if v.Type == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"type\" is required")
	}

	if v.Controller == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"serviceEndpoint\" is required")
	}

	// TODO

	return nil
}
