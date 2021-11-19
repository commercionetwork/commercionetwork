package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) MintCCC(goCtx context.Context, msg *types.MsgMintCCC) (*types.MsgMintCCCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var depositAmount int64
	for _, denom := range msg.DepositAmount {
		if denom.Denom == types.CreditsDenom {
			depositAmount = denom.Amount.Int64()
			break
		}
	}
	var postion = types.Position{
		Owner:      msg.Depositor,
		Collateral: depositAmount,
		ID:         msg.ID,
	}

	err := k.NewPosition(
		ctx,
		postion,
	)
	if err != nil {
		return &types.MsgMintCCCResponse{},  errors.Wrap(errors.ErrInvalidRequest, err.Error())
	}
	return &types.MsgMintCCCResponse{
		ID: msg.ID,
	}, nil
}

// TODO IMPLEMENTATION
func (k msgServer) BurnCCC(goCtx context.Context, msg *types.MsgBurnCCC) (*types.MsgBurnCCCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	residualAmount, err := k.RemoveCCC(
		ctx,
		signer,
		msg.ID,
		*msg.Amount,
	)
	if err != nil {
		return &types.MsgBurnCCCResponse{}, err
	}
	residualCredits := sdk.NewCoin(types.CreditsDenom, residualAmount)

	return &types.MsgBurnCCCResponse{
		ID:       msg.ID,
		Residual: &residualCredits,
	}, nil
}
