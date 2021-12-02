package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (msg *MsgSetDid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.ID)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetDid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetDid) Route() string {
	return RouterKey
}

func (msg *MsgSetDid) Type() string {
	return MsgTypeSetDid
}

func (msg *MsgSetDid) ValidateBasic() error {
	// _, err := sdk.AccAddressFromBech32(msg.ID)
	// if err != nil {
	// 	return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	// }

	// if msg.Context != ContextDidV1 {
	// 	return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Invalid context, must be https://www.w3.org/ns/did/v1")
	// }

	// controller, _ := sdk.AccAddressFromBech32(msg.ID)

	// for _, key := range msg.PubKeys {
	// 	if err := key.Validate(); err != nil {
	// 		return err
	// 	}
	// 	keycontroller, _ := sdk.AccAddressFromBech32(key.Controller)
	// 	if !controller.Equals(keycontroller) {
	// 		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Public key controller must match did document id")
	// 	}
	// }

	// var pubKeys PubKeys
	// pubKeys = msg.PubKeys
	// if err := pubKeys.noDuplicates(); err != nil {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	// }

	// if !pubKeys.HasVerificationAndSignatureKey() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "specified public keys are not in the correct format")
	// }
	// /*
	// 	if msg.Proof == nil {
	// 		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "proof not provided")
	// 	}

	// 	if err := msg.Proof.Validate(); err != nil {
	// 		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("proof validation error: %s", err.Error()))
	// 	}
	// */
	// // we have some service, we should validate 'em
	// if msg.Service != nil {
	// 	for i, service := range msg.Service {
	// 		err := service.Validate()
	// 		if err != nil {
	// 			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("service %d validation failed: %s", i, err.Error()))
	// 		}
	// 	}
	// }
	// /*
	// 	if err := msg.VerifyProof(); err != nil {
	// 		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
	// 	}*/

	// if err := msg.lengthLimits(); err != nil {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	// }

	return nil
}
