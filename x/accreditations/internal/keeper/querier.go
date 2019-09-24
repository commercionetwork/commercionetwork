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
		case types.QueryGetAccrediter:
			return queryGetAccrediter(ctx, path[1:], keeper)
		case types.QueryGetSigners:
			return queryGetSigners(ctx, path[1:], keeper)
		case types.QueryGetPoolFunds:
			return queryGetPoolFunds(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

type AccrediterResponse struct {
	User       sdk.AccAddress `json:"user"`
	Accrediter sdk.AccAddress `json:"accrediter"`
}

func queryGetAccrediter(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	accreditation := keeper.GetAccreditation(ctx, address)
	response := AccrediterResponse{
		Accrediter: accreditation.Accrediter,
		User:       address,
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, response)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSigners(ctx sdk.Context, _ []string, keeper Keeper) (res []byte, err sdk.Error) {
	signers := keeper.GetTrustedSigners(ctx)
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
