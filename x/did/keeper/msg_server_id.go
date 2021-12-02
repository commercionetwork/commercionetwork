package keeper

// func (k msgServer) SetIdentity(goCtx context.Context, msg *types.MsgSetIdentity) (*types.MsgSetIdentityResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(goCtx)

// 	var didDocument = types.DidDocument{
// 		Context: msg.Context,
// 		ID:      msg.ID,
// 		PubKeys: msg.PubKeys,
// 		//Proof:   msg.Proof,
// 		Service: msg.Service,
// 	}

// 	id := k.AppendId(
// 		ctx,
// 		didDocument,
// 	)

// 	return &types.MsgSetIdentityResponse{
// 		ID: id,
// 	}, nil
// }
