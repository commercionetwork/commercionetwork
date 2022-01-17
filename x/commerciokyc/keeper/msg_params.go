package keeper

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SetParams(goCtx context.Context, msg *types.MsgSetParams) (*types.MsgSetParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	gov := k.GovKeeper.GetGovernmentAddress(ctx)
	msgSigner, e := sdk.AccAddressFromBech32(msg.Signer)
	if e != nil {
		return nil, e
	}
	if !gov.Equals(msgSigner) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set commerciokyc params", msg.Signer))
	}

	k.UpdateParams(ctx, *msg.Params)

	return &types.MsgSetParamsResponse{}, nil
}
