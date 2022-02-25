package keeper

import (
	"context"
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context, bankKeeper.Keeper) {
	keeper, ctx := setupKeeper(t)
	keeper.govKeeper.SetGovernmentAddress(ctx, TestFunder)
	return NewMsgServerImpl(*keeper), sdk.WrapSDKContext(ctx), keeper.bankKeeper
}

func TestSetParams(t *testing.T) {
	srv, ctx, _ := setupMsgServer(t)
	for _, tc := range []struct {
		desc string
		msg  *types.MsgSetParams
		err  error
	}{
		{
			desc: "regular params",
			msg: &types.MsgSetParams{
				Government:           TestFunder.String(),
				DistrEpochIdentifier: types.EpochDay,
				EarnRate:             sdk.NewDecWithPrec(5, 1),
			},
			err: nil,
		},
		{
			desc: "inavlid government address",
			msg: &types.MsgSetParams{
				Government:           valDelAddr.String(),
				DistrEpochIdentifier: types.EpochDay,
				EarnRate:             sdk.NewDecWithPrec(5, 1),
			},
			err: sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set params", valDelAddr.String())),
		},
		{
			desc: "invalid epoch identifier",
			msg: &types.MsgSetParams{
				Government:           TestFunder.String(),
				DistrEpochIdentifier: "",
				EarnRate:             sdk.NewDecWithPrec(5, 1),
			},
			err: sdkErr.Wrap(sdkErr.ErrInvalidType, "invalid epoch identifier: "),
		},
		{
			desc: "invalid earn rate",
			msg: &types.MsgSetParams{
				Government:           TestFunder.String(),
				DistrEpochIdentifier: types.EpochDay,
				EarnRate:             sdk.NewDec(-1),
			},
			err: sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("invalid vbr earn rate: %s", sdk.NewDec(-1))),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			_, err := srv.SetParams(ctx, tc.msg)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
