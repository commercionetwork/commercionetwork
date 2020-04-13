package customstaking_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/simapp"
)

func TestStaking(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, abci.Header{})

	tokens := types.TokensFromConsensusPower(1000)
	testAddrs := simapp.AddTestAddrs(app, ctx, 1, tokens)
	valAddrs := simapp.ConvertAddrsToValAddrs(testAddrs)
	pks := simapp.CreateTestPubKeys(1)

	validators := app.StakingKeeper.GetAllValidators(ctx)
	require.Len(t, validators, 0)

	commission := staking.NewCommissionRates(types.NewDecWithPrec(5, 1), types.NewDecWithPrec(5, 1), types.NewDec(0))
	msg := staking.NewMsgCreateValidator(
		valAddrs[0], pks[0], types.NewCoin(types.DefaultBondDenom, types.NewInt(1000)), staking.Description{}, commission, types.OneInt(),
	)

	sh := staking.NewHandler(app.StakingKeeper)

	res, err := sh(ctx, msg)
	require.NoError(t, err)
	require.NotNil(t, res)

	// end block to bond validator
	staking.EndBlocker(ctx, app.StakingKeeper)

	validators = app.StakingKeeper.GetAllValidators(ctx)
	require.Len(t, validators, 1)
}
