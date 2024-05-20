package apptesting

// import (
// 	//"encoding/json"
// 	"fmt"
// 	//"testing"
// 	"time"

// 	"github.com/cosmos/cosmos-sdk/baseapp"
// 	"github.com/cosmos/cosmos-sdk/client"

// 	//cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
// 	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
// 	//"cosmossdk.io/simapp"
// 	"cosmossdk.io/math"
// 	"cosmossdk.io/store/rootmulti"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/types/tx/signing"
// 	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
// 	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

// 	//"github.com/cosmos/cosmos-sdk/x/staking"
// 	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

// 	//"github.com/stretchr/testify/require"
// 	"cosmossdk.io/log"
// 	abci "github.com/cometbft/cometbft/abci/types"
// 	"github.com/cometbft/cometbft/crypto/ed25519"
// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
// 	dbm "github.com/cosmos/cosmos-db"
// 	"github.com/stretchr/testify/suite"

// 	//"github.com/commercionetwork/commercionetwork/x/ibc-address-limiter/types"

// 	"github.com/commercionetwork/commercionetwork/app"
// 	"github.com/commercionetwork/commercionetwork/testutil/simapp"
// )

// //var testUser3, _ = sdk.AccAddressFromBech32("cosmos14lultfckehtszvzw4ehu0apvsr77afvyhgqhwh")

// type KeeperTestHelper struct {
// 	suite.Suite

// 	App         *app.App
// 	Ctx         sdk.Context
// 	QueryHelper *baseapp.QueryServiceTestHelper
// 	//queryClient types.QueryClient
// 	TestAccs []sdk.AccAddress
// }

// var (
// 	SecondaryDenom  = "uccc"
// 	SecondaryAmount = math.NewInt(100000000)
// )

// // Setup sets up basic environment for suite (App, Ctx, and test accounts)
// func (s *KeeperTestHelper) Setup() {
// 	//s.App = app.Setup(false)
// 	s.App = simapp.New("")
// 	s.Ctx = s.App.BaseApp.NewContext(false, tmproto.Header{})
// 	s.QueryHelper = &baseapp.QueryServiceTestHelper{
// 		GRPCQueryRouter: s.App.GRPCQueryRouter(),
// 		Ctx:             s.Ctx,
// 	}

// 	s.SetEpochStartTime()
// 	//s.App.GovernmentKeeper.SetGovernmentAddress(s.Ctx, testUser3)
// 	s.TestAccs = CreateRandomAccounts(3)
// }

// func (s *KeeperTestHelper) SetupTestForInitGenesis() {
// 	// Setting to True, leads to init genesis not running
// 	s.App = app.Setup(true)
// 	s.Ctx = s.App.BaseApp.NewContext(true, tmproto.Header{})
// }

// func (s *KeeperTestHelper) SetEpochStartTime() {
// 	epochsKeeper := s.App.EpochsKeeper

// 	for _, epoch := range epochsKeeper.AllEpochInfos(s.Ctx) {
// 		epoch.StartTime = s.Ctx.BlockTime()
// 		epochsKeeper.DeleteEpochInfo(s.Ctx, epoch.Identifier)
// 		epochsKeeper.SetEpochInfo(s.Ctx, epoch)
// 		/*if err != nil {
// 			panic(err)
// 		}*/
// 	}
// }

// // CreateTestContext creates a test context.
// func (s *KeeperTestHelper) CreateTestContext() sdk.Context {
// 	ctx, _ := s.CreateTestContextWithMultiStore()
// 	return ctx
// }

// // CreateTestContextWithMultiStore creates a test context and returns it together with multi store.
// func (s *KeeperTestHelper) CreateTestContextWithMultiStore() (sdk.Context, sdk.CommitMultiStore) {
// 	db := dbm.NewMemDB()
// 	logger := log.NewNopLogger()

// 	ms := rootmulti.NewStore(db, logger)

// 	return sdk.NewContext(ms, tmproto.Header{}, false, logger), ms
// }

// // CreateTestContext creates a test context.
// func (s *KeeperTestHelper) Commit() {
// 	oldHeight := s.Ctx.BlockHeight()
// 	oldHeader := s.Ctx.BlockHeader()
// 	s.App.Commit()
// 	newHeader := tmproto.Header{Height: oldHeight + 1, ChainID: oldHeader.ChainID, Time: oldHeader.Time.Add(time.Second)}
// 	s.App.BeginBlock(abci.RequestBeginBlock{Header: newHeader})
// 	s.Ctx = s.App.NewContext(false, newHeader)

// }

// /*
// // FundAcc funds target address with specified amount.
// func (s *KeeperTestHelper) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
// 	err := app.FundAccount(s.App.BankKeeper, s.Ctx, acc, amounts)
// 	s.Require().NoError(err)
// }*/
// /*
// // FundModuleAcc funds target modules with specified amount.
// func (s *KeeperTestHelper) FundModuleAcc(moduleName string, amounts sdk.Coins) {
// 	err := app.FundModuleAccount(s.App.BankKeeper, s.Ctx, moduleName, amounts)
// 	s.Require().NoError(err)
// }*/
// /*
// func (s *KeeperTestHelper) MintCoins(coins sdk.Coins) {
// 	err := s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, coins)
// 	s.Require().NoError(err)
// }*/

// // SetupValidator sets up a validator and returns the ValAddress.
// func (s *KeeperTestHelper) SetupValidator(bondStatus stakingtypes.BondStatus) sdk.ValAddress {
// 	valPub := secp256k1.GenPrivKey().PubKey()
// 	valAddr := sdk.ValAddress(valPub.Address())
// 	//bondDenom := s.App.StakingKeeper.GetParams(s.Ctx).BondDenom
// 	//selfBond := sdk.NewCoins(sdk.Coin{Amount: math.NewInt(100), Denom: bondDenom})

// 	//s.FundAcc(sdk.AccAddress(valAddr), selfBond)

// 	//stakingHandler := staking.NewHandler(s.App.StakingKeeper)
// 	//stakingCoin := sdk.NewCoin(sdk.DefaultBondDenom, selfBond[0].Amount)
// 	//ZeroCommission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
// 	//msg, err := stakingtypes.NewMsgCreateValidator(valAddr, valPub, stakingCoin, stakingtypes.Description{}, ZeroCommission, math.OneInt())
// 	//s.Require().NoError(err)
// 	//res, err := stakingHandler(s.Ctx, msg)
// 	//s.Require().NoError(err)
// 	//s.Require().NotNil(res)

// 	val, found := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
// 	s.Require().True(found)

// 	val = val.UpdateStatus(bondStatus)
// 	s.App.StakingKeeper.SetValidator(s.Ctx, val)

// 	consAddr, err := val.GetConsAddr()
// 	s.Suite.Require().NoError(err)

// 	signingInfo := slashingtypes.NewValidatorSigningInfo(
// 		consAddr,
// 		s.Ctx.BlockHeight(),
// 		0,
// 		time.Unix(0, 0),
// 		false,
// 		0,
// 	)
// 	s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, consAddr, signingInfo)

// 	return valAddr
// }

// // SetupMultipleValidators setups "numValidator" validators and returns their address in string
// func (s *KeeperTestHelper) SetupMultipleValidators(numValidator int) []string {
// 	valAddrs := []string{}
// 	for i := 0; i < numValidator; i++ {
// 		valAddr := s.SetupValidator(stakingtypes.Bonded)
// 		valAddrs = append(valAddrs, valAddr.String())
// 	}
// 	return valAddrs
// }

// /*
// // BeginNewBlock starts a new block.
// func (s *KeeperTestHelper) BeginNewBlock(executeNextEpoch bool) {
// 	var valAddr []byte

// 	validators := s.App.StakingKeeper.GetAllValidators(s.Ctx)
// 	if len(validators) >= 1 {
// 		valAddrFancy, err := validators[0].GetConsAddr()
// 		s.Require().NoError(err)
// 		valAddr = valAddrFancy.Bytes()
// 	} else {
// 		valAddrFancy := s.SetupValidator(stakingtypes.Bonded)
// 		validator, _ := s.App.StakingKeeper.GetValidator(s.Ctx, valAddrFancy)
// 		valAddr2, _ := validator.GetConsAddr()
// 		valAddr = valAddr2.Bytes()
// 	}

// 	s.BeginNewBlockWithProposer(executeNextEpoch, valAddr)
// }

// // BeginNewBlockWithProposer begins a new block with a proposer.
// func (s *KeeperTestHelper) BeginNewBlockWithProposer(executeNextEpoch bool, proposer sdk.ValAddress) {
// 	validator, found := s.App.StakingKeeper.GetValidator(s.Ctx, proposer)
// 	s.Assert().True(found)

// 	valConsAddr, err := validator.GetConsAddr()
// 	s.Require().NoError(err)

// 	valAddr := valConsAddr.Bytes()

// 	epochIdentifier := s.App.SuperfluidKeeper.GetEpochIdentifier(s.Ctx)
// 	epoch := s.App.EpochsKeeper.GetEpochInfo(s.Ctx, epochIdentifier)
// 	newBlockTime := s.Ctx.BlockTime().Add(5 * time.Second)
// 	if executeNextEpoch {
// 		newBlockTime = s.Ctx.BlockTime().Add(epoch.Duration).Add(time.Second)
// 	}

// 	header := tmtypes.Header{Height: s.Ctx.BlockHeight() + 1, Time: newBlockTime}
// 	newCtx := s.Ctx.WithBlockTime(newBlockTime).WithBlockHeight(s.Ctx.BlockHeight() + 1)
// 	s.Ctx = newCtx
// 	lastCommitInfo := abci.LastCommitInfo{
// 		Votes: []abci.VoteInfo{{
// 			Validator:       abci.Validator{Address: valAddr, Power: 1000},
// 			SignedLastBlock: true,
// 		}},
// 	}
// 	reqBeginBlock := abci.RequestBeginBlock{Header: header, LastCommitInfo: lastCommitInfo}

// 	fmt.Println("beginning block ", s.Ctx.BlockHeight())
// 	s.App.BeginBlocker(s.Ctx, reqBeginBlock)
// 	s.Ctx = s.App.NewContext(false, reqBeginBlock.Header)
// }*/

// // EndBlock ends the block, and runs commit
// func (s *KeeperTestHelper) EndBlock() {
// 	reqEndBlock := abci.RequestEndBlock{Height: s.Ctx.BlockHeight()}
// 	s.App.EndBlocker(s.Ctx, reqEndBlock)
// }

// func (s *KeeperTestHelper) RunMsg(msg sdk.Msg) (*sdk.Result, error) {
// 	// cursed that we have to copy this internal logic from SDK
// 	router := s.App.MsgServiceRouter()
// 	if handler := router.Handler(msg); handler != nil {
// 		// ADR 031 request type routing
// 		return handler(s.Ctx, msg)
// 	}
// 	s.FailNow("msg %v could not be ran", msg)
// 	return nil, fmt.Errorf("msg %v could not be ran", msg)
// }

// // AllocateRewardsToValidator allocates reward tokens to a distribution module then allocates rewards to the validator address.
// func (s *KeeperTestHelper) AllocateRewardsToValidator(valAddr sdk.ValAddress, rewardAmt math.Int) {
// 	validator, found := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
// 	s.Require().True(found)

// 	// allocate reward tokens to distribution module
// 	//coins := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, rewardAmt)}
// 	//err := simapp.FundModuleAccount(s.App.BankKeeper, s.Ctx, distrtypes.ModuleName, coins)
// 	//s.Require().NoError(err)

// 	// allocate rewards to validator
// 	s.Ctx = s.Ctx.WithBlockHeight(s.Ctx.BlockHeight() + 1)
// 	decTokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDec(20000)}}
// 	s.App.DistrKeeper.AllocateTokensToValidator(s.Ctx, validator, decTokens)
// }

// // BuildTx builds a transaction.
// func (s *KeeperTestHelper) BuildTx(
// 	txBuilder client.TxBuilder,
// 	msgs []sdk.Msg,
// 	sigV2 signing.SignatureV2,
// 	memo string, txFee sdk.Coins,
// 	gasLimit uint64,
// ) authsigning.Tx {
// 	err := txBuilder.SetMsgs(msgs[0])
// 	s.Require().NoError(err)

// 	err = txBuilder.SetSignatures(sigV2)
// 	s.Require().NoError(err)

// 	txBuilder.SetMemo(memo)
// 	txBuilder.SetFeeAmount(txFee)
// 	txBuilder.SetGasLimit(gasLimit)

// 	return txBuilder.GetTx()
// }

// /*
// // StateNotAltered validates that app state is not altered. Fails if it is.
// func (s *KeeperTestHelper) StateNotAltered() {
// 	oldState := s.App.ExportState(s.Ctx)
// 	s.App.Commit()
// 	newState := s.App.ExportState(s.Ctx)
// 	s.Require().Equal(oldState, newState)
// }*/

// // CreateRandomAccounts is a function return a list of randomly generated AccAddresses
// func CreateRandomAccounts(numAccts int) []sdk.AccAddress {
// 	testAddrs := make([]sdk.AccAddress, numAccts)
// 	for i := 0; i < numAccts; i++ {
// 		pk := ed25519.GenPrivKey().PubKey()
// 		testAddrs[i] = sdk.AccAddress(pk.Address())
// 	}

// 	return testAddrs
// }

// /*
// func TestMessageAuthzSerialization(t *testing.T, msg sdk.Msg) {
// 	someDate := time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)
// 	const (
// 		mockGranter string = "cosmos1abc"
// 		mockGrantee string = "cosmos1xyz"
// 	)

// 	var (
// 		mockMsgGrant  authz.MsgGrant
// 		mockMsgRevoke authz.MsgRevoke
// 		mockMsgExec   authz.MsgExec
// 	)

// 	// Authz: Grant Msg
// 	typeURL := sdk.MsgTypeURL(msg)
// 	grant, err := authz.NewGrant(someDate, authz.NewGenericAuthorization(typeURL), someDate.Add(time.Hour))
// 	require.NoError(t, err)

// 	msgGrant := authz.MsgGrant{Granter: mockGranter, Grantee: mockGrantee, Grant: grant}
// 	msgGrantBytes := json.RawMessage(sdk.MustSortJSON(authzcodec.ModuleCdc.MustMarshalJSON(&msgGrant)))
// 	err = authzcodec.ModuleCdc.UnmarshalJSON(msgGrantBytes, &mockMsgGrant)
// 	require.NoError(t, err)

// 	// Authz: Revoke Msg
// 	msgRevoke := authz.MsgRevoke{Granter: mockGranter, Grantee: mockGrantee, MsgTypeUrl: typeURL}
// 	msgRevokeByte := json.RawMessage(sdk.MustSortJSON(authzcodec.ModuleCdc.MustMarshalJSON(&msgRevoke)))
// 	err = authzcodec.ModuleCdc.UnmarshalJSON(msgRevokeByte, &mockMsgRevoke)
// 	require.NoError(t, err)

// 	// Authz: Exec Msg
// 	msgAny, err := cdctypes.NewAnyWithValue(msg)
// 	require.NoError(t, err)
// 	msgExec := authz.MsgExec{Grantee: mockGrantee, Msgs: []*cdctypes.Any{msgAny}}
// 	execMsgByte := json.RawMessage(sdk.MustSortJSON(authzcodec.ModuleCdc.MustMarshalJSON(&msgExec)))
// 	err = authzcodec.ModuleCdc.UnmarshalJSON(execMsgByte, &mockMsgExec)
// 	require.NoError(t, err)
// 	require.Equal(t, msgExec.Msgs[0].Value, mockMsgExec.Msgs[0].Value)
// }*/

// func GenerateTestAddrs() (string, string) {
// 	pk1 := ed25519.GenPrivKey().PubKey()
// 	validAddr := sdk.AccAddress(pk1.Address()).String()
// 	invalidAddr := sdk.AccAddress("invalid").String()
// 	return validAddr, invalidAddr
// }
