package keeper

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/pricefeed"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types"
)

func TestCdpsForExistingPrice(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(Keeper, bank.Keeper, pricefeed.Keeper, types.Context) error
		wantFail  bool
	}{
		{
			"Each Position opened refers to an existing price",
			func(k Keeper, bk bank.Keeper, pfk pricefeed.Keeper, ctx types.Context) error {
				err := bk.SetCoins(ctx, testCdpOwner, sdk.NewCoins(testCdp.Deposit))
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeed.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

				return k.NewPosition(ctx, testCdpOwner, testCdp.Deposit)
			},
			false,
		},
		{
			"Position opened with corresponding price set to zero values (no value, no expiry)",
			func(k Keeper, bk bank.Keeper, pfk pricefeed.Keeper, ctx types.Context) error {
				err := bk.SetCoins(ctx, testCdpOwner, sdk.NewCoins(testCdp.Deposit))
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeed.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

				err = k.NewPosition(ctx, testCdpOwner, testCdp.Deposit)
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeed.NewPrice(testLiquidityDenom, sdk.NewDec(0), sdk.NewInt(0)))

				return nil
			},
			true,
		},
		{
			"Position opened with corresponding price nonexistant",
			func(k Keeper, bk bank.Keeper, pfk pricefeed.Keeper, ctx types.Context) error {
				err := bk.SetCoins(ctx, testCdpOwner, sdk.NewCoins(testCdp.Deposit))
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeed.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

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
			func(k Keeper, bk bank.Keeper, pfk pricefeed.Keeper, ctx types.Context) error {
				return nil
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, bk, pfk, _, _, k := SetupTestInput()

			require.NoError(t, tt.setupFunc(k, bk, pfk, ctx))

			_, failed := CdpsForExistingPrice(k)(ctx)

			require.Equal(t, tt.wantFail, failed)
		})
	}
}

func TestLiquidityPoolAmountEqualsCdps(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(Keeper, bank.Keeper, pricefeed.Keeper, types.Context) error
		wantFail  bool
	}{
		{
			"One cdp opened equals the value of the liquidity pool",
			func(k Keeper, bk bank.Keeper, pfk pricefeed.Keeper, ctx types.Context) error {
				err := bk.SetCoins(ctx, testCdpOwner, sdk.NewCoins(testCdp.Deposit))
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeed.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

				return k.NewPosition(ctx, testCdpOwner, testCdp.Deposit)
			},
			false,
		},
		{
			"One cdp opened and the liquidity pool is zero",
			func(k Keeper, bk bank.Keeper, pfk pricefeed.Keeper, ctx types.Context) error {
				err := bk.SetCoins(ctx, testCdpOwner, sdk.NewCoins(testCdp.Deposit))
				if err != nil {
					return err
				}

				pfk.SetCurrentPrice(ctx, pricefeed.NewPrice(testLiquidityDenom, sdk.NewDec(10), sdk.NewInt(1000)))

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
			require.NoError(t, tt.setupFunc(k, bk, pfk, ctx))
			_, failed := LiquidityPoolAmountEqualsCdps(k)(ctx)

			require.Equal(t, tt.wantFail, failed)
		})
	}
}
