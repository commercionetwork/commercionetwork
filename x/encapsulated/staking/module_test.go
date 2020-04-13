package customstaking_test

import (
	"testing"

	app2 "github.com/commercionetwork/commercionetwork/app"

	"github.com/commercionetwork/commercionetwork/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	customstaking "github.com/commercionetwork/commercionetwork/x/encapsulated/staking"
)

func TestHandler_MinimumStakeToJoin_MinimumDepositNotEnough(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{})

	tokens := sdk.TokensFromConsensusPower(1000)
	testAddrs := simapp.AddTestAddrs(app, ctx, 1, tokens)
	valAddrs := simapp.ConvertAddrsToValAddrs(testAddrs)
	pks := simapp.CreateTestPubKeys(1)

	validators := app.StakingKeeper.GetAllValidators(ctx)
	require.Len(t, validators, 0)

	commission := staking.NewCommissionRates(types.NewDecWithPrec(5, 1), types.NewDecWithPrec(5, 1), types.NewDec(0))
	msg := staking.NewMsgCreateValidator(
		valAddrs[0], pks[0], types.NewCoin(app2.DefaultBondDenom, types.NewInt(1000)), staking.Description{}, commission, types.OneInt(),
	)

	sh := customstaking.NewHandler(app.StakingKeeper)

	res, err := sh(ctx, msg)
	require.EqualError(t, err, customstaking.ErrMinimumDeposit.Error())
	require.Nil(t, res)

	// end block to bond validator
	staking.EndBlocker(ctx, app.StakingKeeper)

	validators = app.StakingKeeper.GetAllValidators(ctx)
	require.Len(t, validators, 0)
}

func TestHandler_MinimumStakeToJoin_MinimumDepositEnough(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{})

	tokens := sdk.TokensFromConsensusPower(60000)
	testAddrs := simapp.AddTestAddrs(app, ctx, 1, tokens)
	valAddrs := simapp.ConvertAddrsToValAddrs(testAddrs)
	pks := simapp.CreateTestPubKeys(1)

	validators := app.StakingKeeper.GetAllValidators(ctx)
	require.Len(t, validators, 0)

	commission := staking.NewCommissionRates(types.NewDecWithPrec(5, 1), types.NewDecWithPrec(5, 1), types.NewDec(0))
	msg := staking.NewMsgCreateValidator(
		valAddrs[0], pks[0], types.NewCoin(app2.DefaultBondDenom, customstaking.MinimumDeposit.Amount), staking.Description{}, commission, types.OneInt(),
	)

	sh := customstaking.NewHandler(app.StakingKeeper)

	res, err := sh(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond validator
	staking.EndBlocker(ctx, app.StakingKeeper)

	validators = app.StakingKeeper.GetAllValidators(ctx)
	require.Len(t, validators, 1)
}
