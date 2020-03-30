package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper, govKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgSetPrice:
			return handleMsgSetPrice(ctx, keeper, msg)
		case types.MsgAddOracle:
			return handleMsgAddOracle(ctx, keeper, govKeeper, msg)
		case types.MsgBlacklistDenom:
			return handleMsgBlacklistDenom(ctx, keeper, govKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, (errMsg))
		}
	}
}

func handleMsgSetPrice(ctx sdk.Context, keeper Keeper, msg types.MsgSetPrice) (*sdk.Result, error) {
	// Check the signer
	if !keeper.IsOracle(ctx, msg.Oracle) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("%s is not an oracle", msg.Oracle.String()))
	}

	for _, denom := range keeper.DenomBlacklist(ctx) {
		if denom == msg.Price.AssetName {
			return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, msg.Price.AssetName+" has been blacklisted")
		}
	}

	// Set the raw price
	if err := keeper.AddRawPrice(ctx, msg.Oracle, msg.Price); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, (err.Error()))
	}
	return &sdk.Result{}, nil
}

func handleMsgAddOracle(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper, msg types.MsgAddOracle) (*sdk.Result, error) {
	gov := govKeeper.GetGovernmentAddress(ctx)

	// Someone who's not the government is trying to add an oracle
	if !(gov.Equals(msg.Signer)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("%s hasn't the rights to add an oracle", msg.Signer))
	}

	keeper.AddOracle(ctx, msg.Oracle)
	return &sdk.Result{}, nil
}

func handleMsgBlacklistDenom(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper, msg types.MsgBlacklistDenom) (*sdk.Result, error) {
	if !msg.Signer.Equals(govKeeper.GetGovernmentAddress(ctx)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("%s hasn't the rights to blacklist a denom", msg.Signer))
	}

	keeper.BlacklistDenom(ctx, msg.Denom)

	return &sdk.Result{}, nil
}
