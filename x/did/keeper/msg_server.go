package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SetDid(goCtx context.Context, msg *types.MsgSetDid) (*types.MsgSetDidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var didDocument = types.DidDocument{
		Context: msg.Context,
		ID:      msg.ID,
		// PubKeys: msg.PubKeys,
		//Proof:   msg.Proof,
		// Service: msg.Service,
	}

	id := k.AppendId(
		ctx,
		didDocument,
	)

	return &types.MsgSetDidResponse{
		ID: id,
	}, nil
}
