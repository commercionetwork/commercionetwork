package keeper

// func TestIdentityQuerySingle(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	wctx := sdk.WrapSDKContext(ctx)
// 	msgs := createNIdentity(keeper, ctx, 2)
// 	for _, tc := range []struct {
// 		desc     string
// 		request  *types.QueryResolveDidRequest
// 		response *types.QueryResolveDidResponse
// 		err      error
// 	}{
// 		{
// 			desc:     "First",
// 			request:  &types.QueryResolveDidRequest{ID: msgs[0].ID},
// 			response: &types.QueryResolveDidResponse{DidDocument: &msgs[0]},
// 		},
// 		{
// 			desc:     "Second",
// 			request:  &types.QueryResolveDidRequest{ID: msgs[1].ID},
// 			response: &types.QueryResolveDidResponse{DidDocument: &msgs[1]},
// 		},
// 		{
// 			desc:    "KeyNotFound",
// 			request: &types.QueryResolveDidRequest{ID: "x"},
// 			err:     sdkerrors.ErrKeyNotFound,
// 		},
// 		{
// 			desc: "InvalidRequest",
// 			err:  status.Error(codes.InvalidArgument, "invalid request"),
// 		},
// 	} {
// 		tc := tc
// 		t.Run(tc.desc, func(t *testing.T) {
// 			response, err := keeper.Identity(wctx, tc.request)
// 			if tc.err != nil {
// 				require.ErrorIs(t, err, tc.err)
// 			} else {
// 				require.Equal(t, tc.response, response)
// 			}
// 		})
// 	}
// }
