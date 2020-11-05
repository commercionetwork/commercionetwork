package keeper

import (
	"fmt"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

var testMsgMintCCC = types.NewMsgMintCCC(testEtp.Owner, sdk.NewCoins(sdk.NewCoin("uccc", testEtp.Collateral)))
var testMsgBurnCCC = types.NewMsgBurnCCC(testEtp.Owner, testEtp.ID, testEtp.Credits)

func TestHandler_handleMsgMintCCC(t *testing.T) {
	ctx, bk, _, _, _, k := SetupTestInput()
	handler := NewHandler(k)
	ctx = ctx.WithBlockHeight(5)

	// Test setup
	_, _ = bk.AddCoins(ctx, testEtp.Owner, sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(200))))
	balance := bk.GetCoins(ctx, testEtpOwner)

	// Check balance
	require.Equal(t, "200ucommercio", balance.String())

	actual, err := handler(ctx, testMsgMintCCC)
	require.NoError(t, err)
	require.Equal(t, &sdk.Result{Log: "mint successful"}, actual)

	// Check final balance
	balance = bk.GetCoins(ctx, testEtpOwner)
	require.Equal(t, "100uccc", balance.String())
}

func TestHandler_handleMsgBurnCCC(t *testing.T) {
	ctx, bk, _, _, _, k := SetupTestInput()
	handler := NewHandler(k)
	ctx = ctx.WithBlockHeight(5)

	_, _ = bk.AddCoins(ctx, k.supplyKeeper.GetModuleAddress(types.ModuleName), sdk.NewCoins(sdk.NewCoin("ucommercio", testEtp.Collateral)))
	_ = bk.SetCoins(ctx, testEtp.Owner, sdk.NewCoins(testEtp.Credits))
	require.Equal(t, 0, len(k.GetAllPositions(ctx)))
	k.SetPosition(ctx, testEtp)
	require.Equal(t, 1, len(k.GetAllPositions(ctx)))

	expected := &sdk.Result{Log: "burn successful"}
	actual, err := handler(ctx, testMsgBurnCCC)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestHandler_handleMsgSetCCCConversionRate(t *testing.T) {
	ctx, _, _, gk, _, k := SetupTestInput()
	govAddr := []byte("governance")
	gk.SetGovernmentAddress(ctx, govAddr)
	handler := NewHandler(k)
	ctx = ctx.WithBlockHeight(5)

	msg := types.NewMsgSetCCCConversionRate(govAddr, sdk.NewDec(3))

	expected := &sdk.Result{Log: "conversion rate changed successfully to 3.000000000000000000"}
	actual, err := handler(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	msg = types.NewMsgSetCCCConversionRate(govAddr, sdk.NewDec(0))
	_, err = handler(ctx, msg)
	require.Error(t, err)
	require.Equal(t, "invalid request: conversion rate cannot be zero", err.Error())

	msg = types.NewMsgSetCCCConversionRate([]byte("invalidAddr"), sdk.NewDec(3))
	_, err = handler(ctx, msg)
	require.Error(t, err)
	require.Equal(t, "unauthorized: cosmos1d9h8vctvd9jyzerywgt84wdv cannot set conversion rate", err.Error())
}

func TestHandler_InvalidMsg(t *testing.T) {
	ctx, _, _, _, _, k := SetupTestInput()
	handler := NewHandler(k)
	ctx = ctx.WithBlockHeight(5)

	invalidMsg := sdk.NewTestMsg()
	errMsg := fmt.Sprintf("unrecognized %s message type: %v", types.ModuleName, invalidMsg.Type())
	expected := sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)

	_, err := handler(ctx, invalidMsg)
	require.Error(t, err)
	require.Equal(t, expected.Error(), err.Error())
}
