package keeper

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/government/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
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

// TODO: SetGovAddress method does nothing. In the future should reconfigure government address
func (k msgServer) SetGovAddress(goCtx context.Context, msg *types.MsgSetGovAddress) (*types.MsgSetGovAddressResponse, error) {

	return &types.MsgSetGovAddressResponse{}, nil
}

func (k msgServer) FixSupply(goCtx context.Context, msg *types.MsgFixSupply) (*types.MsgFixSupplyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	gov := k.GetGovernmentAddress(ctx)

	signer, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	if !(gov.Equals(signer)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot fix supply. Gov must be %s", msg.Sender, gov.String()))
	}

	err = k.FixSupplyKeeper(
		ctx,
		signer,
		*msg.Amount,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgFixSupplyResponse{}, nil
}
