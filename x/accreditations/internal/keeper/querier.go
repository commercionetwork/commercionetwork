package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryGetInvites:
			return queryGetInvites(ctx, path[1:], keeper)
		case types.QueryGetTrustedServiceProviders:
			return queryGetSigners(ctx, path[1:], keeper)
		case types.QueryGetPoolFunds:
			return queryGetPoolFunds(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

type InviteResponse struct {
	Recipient sdk.AccAddress `json:"recipient"`
	Sender    sdk.AccAddress `json:"sender"`
}

func queryGetInvites(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {

	// Get the list of invites
	var invites []types.Invite
	if len(path) == 0 {

		// Get all the invites
		invites = keeper.GetInvites(ctx)

	} else {

		// A user has been specified, get only his invites
		address, _ := sdk.AccAddressFromBech32(path[0])
		if invite, found := keeper.GetInvite(ctx, address); found {
			invites = []types.Invite{invite}
		}
	}

	if invites == nil {
		invites = make([]types.Invite, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, invites)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSigners(ctx sdk.Context, _ []string, keeper Keeper) (res []byte, err sdk.Error) {
	signers := keeper.GetTrustedServiceProviders(ctx)
	if signers == nil {
		signers = make([]sdk.AccAddress, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, signers)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetPoolFunds(ctx sdk.Context, _ []string, keeper Keeper) (res []byte, err sdk.Error) {
	value := keeper.GetPoolFunds(ctx)
	if value == nil {
		value = make([]sdk.Coin, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, value)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}
