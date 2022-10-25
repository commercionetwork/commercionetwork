package keeper

import (
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper_Etp(t *testing.T) {

	type args struct {
		req *types.QueryEtpRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryEtpResponse
		wantErr bool
	}{
		{
			name:    "invalid request",
			wantErr: true,
		},
		{
			name: "valid ID but empty store",
			args: args{
				req: &types.QueryEtpRequest{
					ID: "$" + testEtp.ID,
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{req: &types.QueryEtpRequest{ID: testEtp.ID}},
			want: &types.QueryEtpResponse{
				Position: &testEtp,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			wctx := sdk.WrapSDKContext(ctx)

			require.NoError(t, k.SetPosition(ctx, testEtp))

			got, err := k.Etp(wctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.Params() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.Params() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Etps(t *testing.T) {
	type args struct {
		req *types.QueryEtpsRequest
	}
	tests := []struct {
		name    string
		args    args
		etps    []types.Position
		want    *types.QueryEtpsResponse
		wantErr bool
	}{
		{
			name:    "invalid request",
			wantErr: true,
		},
		{
			name: "empty store",
			args: args{
				req: &types.QueryEtpsRequest{},
			},
			etps: []types.Position{},
			want: &types.QueryEtpsResponse{
				Positions: []*types.Position{},
				Pagination: &query.PageResponse{
					Total: 0,
				},
			},
		},
		{
			name: "find one",
			args: args{
				req: &types.QueryEtpsRequest{},
			},
			etps: []types.Position{testEtp},
			want: &types.QueryEtpsResponse{
				Positions: []*types.Position{&testEtp},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			name: "find more than one",
			args: args{
				req: &types.QueryEtpsRequest{},
			},
			etps: []types.Position{testEtp, testEtp1, testEtpAnotherOwner},
			want: &types.QueryEtpsResponse{
				Positions: []*types.Position{&testEtp, &testEtp1, &testEtpAnotherOwner},
				Pagination: &query.PageResponse{
					Total: 3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			wctx := sdk.WrapSDKContext(ctx)

			for _, etp := range tt.etps {
				require.NoError(t, k.SetPosition(ctx, etp))
			}

			got, err := k.Etps(wctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.Etps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				require.Len(t, got.Positions, len(tt.want.Positions))
				for _, want := range tt.want.Positions {
					require.Contains(t, got.Positions, want)
				}

				require.Equal(t, tt.want.Pagination, got.Pagination)
			}
		})
	}
}

func TestKeeper_Pagination(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	wctx := sdk.WrapSDKContext(ctx)

	etps := []types.Position{testEtp, testEtp1, testEtpAnotherOwner}
	for _, etp := range etps {
		require.NoError(t, k.SetPosition(ctx, etp))
	}

	t.Run("Etps ByOffset with limit 1", func(t *testing.T) {
		req := &types.QueryEtpsRequest{
			Pagination: &query.PageRequest{
				Limit: 1,
			},
		}

		var result []types.Position
		for i := 1; i <= len(etps); i++ {
			got, err := k.Etps(wctx, req)
			require.NoError(t, err)
			for _, p := range got.Positions {
				result = append(result, *p)
			}
			req.Pagination.Offset = uint64(i)
		}

		expected := k.GetAllPositions(ctx)

		assert.Len(t, result, len(expected))
		for _, want := range expected {
			assert.Contains(t, result, *want)
		}
	})

	t.Run("Etps ByKey with limit 1", func(t *testing.T) {

		var next []byte

		req := &types.QueryEtpsRequest{
			Pagination: &query.PageRequest{
				Key:   next,
				Limit: 1,
			},
		}

		var result []types.Position
		for i := 1; i <= len(etps); i++ {
			got, err := k.Etps(wctx, req)
			require.NoError(t, err)
			for _, p := range got.Positions {
				result = append(result, *p)
			}

			req.Pagination.Key = got.Pagination.NextKey
		}

		expected := k.GetAllPositions(ctx)

		assert.Len(t, result, len(expected))
		for _, want := range expected {
			assert.Contains(t, result, *want)
		}
	})

	t.Run("EtpsByOwner ByOffset with limit 1", func(t *testing.T) {
		req := &types.QueryEtpsByOwnerRequest{
			Owner:      testEtpOwner.String(),
			Pagination: &query.PageRequest{Limit: 1},
		}

		expected := k.GetAllPositionsOwnedBy(ctx, testEtpOwner)

		var result []types.Position
		for i := 1; i <= len(expected); i++ {
			got, err := k.EtpsByOwner(wctx, req)
			require.NoError(t, err)
			for _, p := range got.Positions {
				result = append(result, *p)
			}
			req.Pagination.Offset = uint64(i)
		}

		assert.Len(t, result, len(expected))
		for _, want := range expected {
			assert.Contains(t, result, *want)
		}
	})

	t.Run("Etps ByKey with limit 1", func(t *testing.T) {

		var next []byte

		req := &types.QueryEtpsByOwnerRequest{
			Owner:      testEtpOwner.String(),
			Pagination: &query.PageRequest{Key: next, Limit: 1},
		}

		expected := k.GetAllPositionsOwnedBy(ctx, testEtpOwner)

		var result []types.Position
		for i := 1; i <= len(expected); i++ {
			got, err := k.EtpsByOwner(wctx, req)
			require.NoError(t, err)
			for _, p := range got.Positions {
				result = append(result, *p)
			}

			req.Pagination.Key = got.Pagination.NextKey
		}

		assert.Len(t, result, len(expected))
		for _, want := range expected {
			assert.Contains(t, result, *want)
		}
	})

}

func TestKeeper_EtpsByOwner(t *testing.T) {
	type args struct {
		req *types.QueryEtpsByOwnerRequest
	}
	tests := []struct {
		name    string
		args    args
		etps    []types.Position
		want    *types.QueryEtpsResponse
		wantErr bool
	}{
		{
			name:    "invalid request",
			wantErr: true,
		},
		{
			name: "invalid address",
			args: args{
				req: &types.QueryEtpsByOwnerRequest{},
			},
			wantErr: true,
		},
		{
			name: "empty store",
			args: args{
				req: &types.QueryEtpsByOwnerRequest{
					Owner:      testEtpOwner.String(),
					Pagination: &query.PageRequest{},
				},
			},
			etps: []types.Position{},
			want: &types.QueryEtpsResponse{
				Positions: []*types.Position{},
				Pagination: &query.PageResponse{
					Total: 0,
				},
			},
		},
		{
			name: "find one",
			args: args{
				req: &types.QueryEtpsByOwnerRequest{
					Owner:      testEtpOwner.String(),
					Pagination: &query.PageRequest{},
				},
			},
			etps: []types.Position{testEtp},
			want: &types.QueryEtpsResponse{
				Positions: []*types.Position{&testEtp},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			name: "find more than one among same owner",
			args: args{
				req: &types.QueryEtpsByOwnerRequest{
					Owner:      testEtpOwner.String(),
					Pagination: &query.PageRequest{},
				},
			},
			etps: []types.Position{testEtp, testEtp1, testEtp2},
			want: &types.QueryEtpsResponse{
				Positions: []*types.Position{&testEtp, &testEtp1, &testEtp2},
				Pagination: &query.PageResponse{
					Total: 3,
				},
			},
		},

		{
			name: "find more than one among another owner",
			args: args{
				req: &types.QueryEtpsByOwnerRequest{
					Owner:      testEtpOwner.String(),
					Pagination: &query.PageRequest{},
				},
			},
			etps: []types.Position{testEtp, testEtp1, testEtpAnotherOwner},
			want: &types.QueryEtpsResponse{
				Positions: []*types.Position{&testEtp, &testEtp1},
				Pagination: &query.PageResponse{
					Total: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			wctx := sdk.WrapSDKContext(ctx)

			for _, etp := range tt.etps {
				require.NoError(t, k.SetPosition(ctx, etp))
			}

			got, err := k.EtpsByOwner(wctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.EtpsByOwner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Len(t, got.Positions, len(tt.want.Positions))
				for _, want := range tt.want.Positions {
					require.Contains(t, got.Positions, want)
				}

				require.Equal(t, tt.want.Pagination, got.Pagination)
			}
		})
	}
}
