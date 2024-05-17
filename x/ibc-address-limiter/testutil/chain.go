package commercionetworkibctesting

// import (
// 	"encoding/json"

// 	"github.com/commercionetwork/commercionetwork/app/params"
// 	"github.com/commercionetwork/commercionetwork/testutil/simapp"
// 	"github.com/cosmos/cosmos-sdk/baseapp"
// 	"github.com/cosmos/cosmos-sdk/client"
// 	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/ibc-go/testing/simapp/helpers"
// 	ibctesting "github.com/cosmos/ibc-go/v8/testing"

// 	//"github.com/cosmos/ibc-go/v8/testing/simapp/helpers"
// 	abci "github.com/cometbft/cometbft/abci/types"
// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

// 	"github.com/commercionetwork/commercionetwork/app"
// )

// type TestChain struct {
// 	*ibctesting.TestChain
// }

// func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
// 	encodingConfig := params.MakeEncodingConfig()
// 	cdc := encodingConfig.Marshaler
// 	commercionetworkApp := /*app.Setup(false)*/ simapp.New("")

// 	return commercionetworkApp, app.NewDefaultGenesisState(cdc)
// }

// // SendMsgsNoCheck overrides ibctesting.TestChain.SendMsgs so that it doesn't check for errors. That should be handled by the caller
// func (chain *TestChain) SendMsgsNoCheck(msgs ...sdk.Msg) (*sdk.Result, error) {
// 	// ensure the chain has the latest time
// 	chain.Coordinator.UpdateTimeForChain(chain.TestChain)

// 	_, r, err := SignAndDeliver(
// 		chain.TxConfig,
// 		chain.App.GetBaseApp(),
// 		chain.GetContext().BlockHeader(),
// 		msgs,
// 		chain.ChainID,
// 		[]uint64{chain.SenderAccount.GetAccountNumber()},
// 		[]uint64{chain.SenderAccount.GetSequence()},
// 		chain.SenderPrivKey,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// SignAndDeliver calls app.Commit()
// 	chain.NextBlock()

// 	// increment sequence for successful transaction execution
// 	err = chain.SenderAccount.SetSequence(chain.SenderAccount.GetSequence() + 1)
// 	if err != nil {
// 		return nil, err
// 	}

// 	chain.Coordinator.IncrementTime()

// 	return r, nil
// }

// // SignAndDeliver signs and delivers a transaction without asserting the results. This overrides the function
// // from ibctesting
// func SignAndDeliver(
// 	txCfg client.TxConfig, app *baseapp.BaseApp, header tmproto.Header, msgs []sdk.Msg,
// 	chainID string, accNums, accSeqs []uint64, priv ...cryptotypes.PrivKey,
// ) (sdk.GasInfo, *sdk.Result, error) {
// 	tx, _ := helpers.GenTx(
// 		txCfg,
// 		msgs,
// 		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
// 		helpers.DefaultGenTxGas,
// 		chainID,
// 		accNums,
// 		accSeqs,
// 		priv...,
// 	)

// 	// Simulate a sending a transaction and committing a block
// 	app.BeginBlock(abci.RequestBeginBlock{Header: header})
// 	gInfo, res, err := app.Deliver(txCfg.TxEncoder(), tx)

// 	//app.EndBlock(abci.RequestEndBlock{})
// 	//app.Commit()

// 	return gInfo, res, err
// }

// // Move epochs to the future to avoid issues with minting
// /*func (chain *TestChain) MoveEpochsToTheFuture() error {
// 	epochsKeeper := chain.GetApp().EpochsKeeper
// 	ctx := chain.GetContext()
// 	for _, epoch := range epochsKeeper.AllEpochInfos(ctx) {
// 		epoch.StartTime = ctx.BlockTime().Add(time.Hour * 24 * 30)
// 		epochsKeeper.DeleteEpochInfo(chain.GetContext(), epoch.Identifier)
// 		epochsKeeper.SetEpochInfo(ctx, epoch)
// 		/*if err != nil {
// 			return err
// 		}*
// 	}
// 	return nil
// }*/

// // GetApp returns the current chain's app
// func (chain *TestChain) GetApp() *app.App {
// 	v, _ := chain.App.(*app.App)
// 	return v
// }
