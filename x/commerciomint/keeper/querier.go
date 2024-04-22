package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"

	abci "github.com/cometbft/cometbft/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetEtpRest:
			return queryGetEtp(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryGetEtpsByOwnerRest:
			return queryGetEtpsByOwner(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryGetallEtpsRest:
			return queryGetAllEtps(ctx, k, legacyQuerierCdc)
		case types.QueryConversionRateRest:
			return queryGetConversionRate(ctx, k, legacyQuerierCdc)
		case types.QueryFreezePeriodRest:
			return queryGetFreezePeriod(ctx, k, legacyQuerierCdc)
		case types.QueryGetParamsRest:
			return queryGetParams(ctx, k, legacyQuerierCdc)
		default:
			return nil, errorsmod.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetEtp(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	id := path[0]
	etp, ok := k.GetPositionById(ctx, id)
	if !ok {
		return nil, errorsmod.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("position with id: %s not found!", id))
	}

	etpbz, err := codec.MarshalJSONIndent(legacyQuerierCdc, etp)
	if err != nil {
		return nil, errorsmod.Wrap(sdkErr.ErrUnknownRequest, "could not marshal result to JSON")
	}

	return etpbz, nil
}

func queryGetEtpsByOwner(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	ownerAddr, e := sdk.AccAddressFromBech32(path[0])
	if e != nil {
		return nil, errorsmod.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("invalid address %s", path[0]))
	}
	etps := k.GetAllPositionsOwnedBy(ctx, ownerAddr)
	etpsBz, err := codec.MarshalJSONIndent(legacyQuerierCdc, etps)
	if err != nil {
		return nil, errorsmod.Wrap(sdkErr.ErrUnknownRequest, "could not marshal result to JSON")
	}

	return etpsBz, nil
}

func queryGetAllEtps(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	etps := k.GetAllPositions(ctx)
	etpsBz, err := codec.MarshalJSONIndent(legacyQuerierCdc, etps)
	if err != nil {
		return nil, errorsmod.Wrap(sdkErr.ErrUnknownRequest, "could not marshal result to JSON")
	}

	return etpsBz, nil
}

func queryGetConversionRate(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	return codec.MarshalJSONIndent(legacyQuerierCdc, k.GetConversionRate(ctx))
}

func queryGetFreezePeriod(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	return codec.MarshalJSONIndent(legacyQuerierCdc, k.GetFreezePeriod(ctx))
}

func queryGetParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	paramsBz, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, errorsmod.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return paramsBz, nil
}
