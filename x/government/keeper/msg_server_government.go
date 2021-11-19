package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/government/types"
)

func (k msgServer) SetGovAddress(goCtx context.Context, msg *types.MsgSetGovAddress) (*types.MsgSetGovAddressResponse, error) {

	return &types.MsgSetGovAddressResponse{}, nil
}
