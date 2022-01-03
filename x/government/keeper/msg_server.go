package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/government/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) SetGovAddress(goCtx context.Context, msg *types.MsgSetGovAddress) (*types.MsgSetGovAddressResponse, error) {

	return &types.MsgSetGovAddressResponse{}, nil
}
