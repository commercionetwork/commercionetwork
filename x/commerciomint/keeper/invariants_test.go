package keeper

import (
	"fmt"
	"testing"

	pricefeedKeeper "github.com/commercionetwork/commercionetwork/x/pricefeed/keeper"
	pricefeedTypes "github.com/commercionetwork/commercionetwork/x/pricefeed/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

func TestValidatePositions(t *testing.T) {
	tests := []struct {
		name      string
		positions []types.Position
		wantErr   bool
	}{
		{"nil", nil, false},
		{"no owner", []types.Position{{}}, true},
		{"no deposit", []types.Position{{Owner: []byte("ciao")}}, true},
		{"invalid liquidity", []types.Position{{Owner: []byte("ciao"), Deposit: sdk.NewCoins(sdk.NewInt64Coin("test", 10))}}, true},
		{"invalid block height", []types.Position{{Owner: []byte("ciao"),
			Deposit: sdk.NewCoins(sdk.NewInt64Coin("test", 10)), Credits: sdk.NewInt64Coin("credit", 100),
		}}, true},
		{"ok", []types.Position{{Owner: []byte("ciao"), CreatedAt: 10,
			Deposit: sdk.NewCoins(sdk.NewInt64Coin("test", 10)), Credits: sdk.NewInt64Coin("credit", 100),
		}}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, _, _, k := SetupTestInput()
			ctx = ctx.WithBlockHeight(5)
			for _, pos := range tt.positions {
				k.SetPosition(ctx, pos)
			}
			msg, failed := ValidateAllPositions(k)(ctx)
			require.Equal(t, tt.wantErr, failed)
			t.Log(msg)
		})
	}
}

func TestPositionsForExistingPrice(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(Keeper, bank.Keeper, pricefeedKeeper.Keeper, sdk.Context) error
		wantFail  bool
	}{
		{
			"Each Position opened refers to an existing price",
			func(k Keeper, bk bank.Keeper, pfk pricefeedKeeper.Keeper, ctx sdk.Context) error {
				err := bk.SetCoins(ctx, testCdpOwner, testCdp.Deposit)
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeedTypes.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

				return k.NewPosition(ctx, testCdpOwner, testCdp.Deposit)
			},
			false,
		},
		{
			"Position opened with corresponding price set to zero values (no value, no expiry)",
			func(k Keeper, bk bank.Keeper, pfk pricefeedKeeper.Keeper, ctx sdk.Context) error {
				err := bk.SetCoins(ctx, testCdpOwner, testCdp.Deposit)
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeedTypes.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

				err = k.NewPosition(ctx, testCdpOwner, testCdp.Deposit)
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeedTypes.NewPrice(testLiquidityDenom, sdk.NewDec(0), sdk.NewInt(0)))

				return nil
			},
			true,
		},
		{
			"Position opened with corresponding price nonexistant",
			func(k Keeper, bk bank.Keeper, pfk pricefeedKeeper.Keeper, ctx sdk.Context) error {
				ctx = ctx.WithBlockHeight(2)
				err := bk.SetCoins(ctx, testCdpOwner, testCdp.Deposit)
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeedTypes.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

				err = k.NewPosition(ctx, testCdpOwner, testCdp.Deposit)
				if err != nil {
					return err
				}

				// explicitly remove price from pricefeed storage to force the invariant to break
				store := ctx.KVStore(pfk.StoreKey)
				store.Delete([]byte("pricefeed:currentPrices:" + testLiquidityDenom))

				return nil
			},
			true,
		},
		{
			"No cdps and no prices",
			func(k Keeper, bk bank.Keeper, pfk pricefeedKeeper.Keeper, ctx sdk.Context) error {
				return nil
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, bk, pfk, _, _, k := SetupTestInput()
			ctx = ctx.WithBlockHeight(5)

			require.NoError(t, tt.setupFunc(k, bk, pfk, ctx))

			_, failed := PositionsForExistingPrice(k)(ctx)

			require.Equal(t, tt.wantFail, failed)
		})
	}
}

func TestLiquidityPoolAmountEqualsPositions(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(Keeper, bank.Keeper, pricefeedKeeper.Keeper, sdk.Context) error
		wantFail  bool
	}{
		{
			"One cdp opened equals the value of the liquidity pool",
			func(k Keeper, bk bank.Keeper, pfk pricefeedKeeper.Keeper, ctx sdk.Context) error {
				err := bk.SetCoins(ctx, testCdpOwner, testCdp.Deposit)
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeedTypes.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

				return k.NewPosition(ctx, testCdpOwner, testCdp.Deposit)
			},
			false,
		},
		{
			"One cdp opened and the liquidity pool is zero",
			func(k Keeper, bk bank.Keeper, pfk pricefeedKeeper.Keeper, ctx sdk.Context) error {
				err := bk.SetCoins(ctx, testCdpOwner, testCdp.Deposit)
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeedTypes.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

				err = k.NewPosition(ctx, testCdpOwner, testCdp.Deposit)
				if err != nil {
					return err
				}

				macc := k.GetMintModuleAccount(ctx)

				if err := macc.SetCoins(sdk.NewCoins()); err != nil {
					return fmt.Errorf("could not set zero coins to pricefeed account")
				}

				k.supplyKeeper.SetModuleAccount(ctx, macc)

				return nil
			},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, bk, pfk, _, _, k := SetupTestInput()
			ctx = ctx.WithBlockHeight(1)
			require.NoError(t, tt.setupFunc(k, bk, pfk, ctx))
			_, failed := LiquidityPoolAmountEqualsPositions(k)(ctx)

			require.Equal(t, tt.wantFail, failed)
		})
	}
}
