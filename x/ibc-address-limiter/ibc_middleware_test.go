package ibc_address_limit_test

/*
import (
	"fmt"
	"strings"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/stretchr/testify/suite"

	"github.com/commercionetwork/commercionetwork/app/apptesting"
	commercionetworkibctesting "github.com/commercionetwork/commercionetwork/x/ibc-address-limiter/testutil"
	"github.com/commercionetwork/commercionetwork/x/ibc-address-limiter/types"

	//"github.com/commercionetwork/commercionetwork/app"
)

type MiddlewareTestSuite struct {
	apptesting.KeeperTestHelper

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *commercionetworkibctesting.TestChain
	chainB *commercionetworkibctesting.TestChain
	path   *ibctesting.Path
}

// Setup
func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func NewTransferPath(chainA, chainB *commercionetworkibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA.TestChain, chainB.TestChain)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointA.ChannelConfig.Version = transfertypes.Version
	path.EndpointB.ChannelConfig.Version = transfertypes.Version
	return path
}

func (suite *MiddlewareTestSuite) SetupTest() {
	suite.Setup()
	//ibctesting.DefaultTestingAppInit = commercionetworkibctesting.SetupTestingApp
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = &commercionetworkibctesting.TestChain{
		TestChain: suite.coordinator.GetChain(ibctesting.GetChainID(1)),
	}

	suite.chainB = &commercionetworkibctesting.TestChain{
		TestChain: suite.coordinator.GetChain(ibctesting.GetChainID(2)),
	}
	suite.path = NewTransferPath(suite.chainA, suite.chainB)

	suite.coordinator.Setup(suite.path)
}

// Helpers
func (suite *MiddlewareTestSuite) MessageFromAToB(denom string, amount sdk.Int) sdk.Msg {
	coin := sdk.NewCoin(denom, amount)
	port := suite.path.EndpointA.ChannelConfig.PortID
	channel := suite.path.EndpointA.ChannelID
	accountFrom := suite.chainA.SenderAccount.GetAddress().String()
	accountTo := suite.chainB.SenderAccount.GetAddress().String()
	timeoutHeight := clienttypes.NewHeight(0, 100)
	return transfertypes.NewMsgTransfer(
		port,
		channel,
		coin,
		accountFrom,
		accountTo,
		timeoutHeight,
		0,
	)
}

func (suite *MiddlewareTestSuite) MessageFromBToA(denom string, amount sdk.Int) sdk.Msg {
	coin := sdk.NewCoin(denom, amount)
	port := suite.path.EndpointB.ChannelConfig.PortID
	channel := suite.path.EndpointB.ChannelID
	accountFrom := suite.chainB.SenderAccount.GetAddress().String()
	accountTo := suite.chainA.SenderAccount.GetAddress().String()
	timeoutHeight := clienttypes.NewHeight(0, 100)
	return transfertypes.NewMsgTransfer(
		port,
		channel,
		coin,
		accountFrom,
		accountTo,
		timeoutHeight,
		0,
	)
}

// Tests that a receiver address longer than 4096 is not accepted
func (suite *MiddlewareTestSuite) TestInvalidReceiver() {
	suite.Setup()
	msg := transfertypes.NewMsgTransfer(
		suite.path.EndpointB.ChannelConfig.PortID,
		suite.path.EndpointB.ChannelID,
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1)),
		suite.chainB.SenderAccount.GetAddress().String(),
		strings.Repeat("x", 4097),
		clienttypes.NewHeight(0, 100),
		0,
	)
	_, ack, _ := suite.FullSendBToA(msg)
	suite.Require().Contains(string(ack), "error",
		"acknowledgment is not an error")
	suite.Require().Contains(string(ack), sdkerrors.ErrInvalidAddress.Error(),
		"acknowledgment error is not of the right type")
}

func (suite *MiddlewareTestSuite) FullSendBToA(msg sdk.Msg) (*sdk.Result, string, error) {
	sendResult, err := suite.chainB.SendMsgsNoCheck(msg)
	suite.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(sendResult.GetEvents())
	suite.Require().NoError(err)

	err = suite.path.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	res, err := suite.path.EndpointA.RecvPacketWithResult(packet)
	suite.Require().NoError(err)

	ack, _ := ibctesting.ParseAckFromEvents(res.GetEvents())

	err = suite.path.EndpointA.UpdateClient()
	suite.Require().NoError(err)
	err = suite.path.EndpointB.UpdateClient()
	suite.Require().NoError(err)

	return sendResult, string(ack), err
}

func (suite *MiddlewareTestSuite) FullSendAToB(msg sdk.Msg) (*sdk.Result, string, error) {
	sendResult, err := suite.chainA.SendMsgsNoCheck(msg)
	if err != nil {
		return nil, "", err
	}

	packet, err := ibctesting.ParsePacketFromEvents(sendResult.GetEvents())
	if err != nil {
		return nil, "", err
	}

	err = suite.path.EndpointB.UpdateClient()
	if err != nil {
		return nil, "", err
	}

	res, err := suite.path.EndpointB.RecvPacketWithResult(packet)
	if err != nil {
		return nil, "", err
	}

	ack, err := ibctesting.ParseAckFromEvents(res.GetEvents())
	if err != nil {
		return nil, "", err
	}

	err = suite.path.EndpointA.UpdateClient()
	if err != nil {
		return nil, "", err
	}
	err = suite.path.EndpointB.UpdateClient()
	if err != nil {
		return nil, "", err
	}

	return sendResult, string(ack), nil
}

func (suite *MiddlewareTestSuite) AssertReceive(success bool, msg sdk.Msg) (string, error) {
	_, ack, err := suite.FullSendBToA(msg)
	if success {
		suite.Require().NoError(err)
		suite.Require().NotContains(string(ack), "error",
			"acknowledgment is an error")
	} else {
		suite.Require().Contains(string(ack), "error",
			"acknowledgment is not an error")
		suite.Require().Contains(string(ack), types.ErrUnauthorized.Error(),
			"acknowledgment error is not of the right type")
	}
	return ack, err
}

func (suite *MiddlewareTestSuite) AssertSend(success bool, msg sdk.Msg) (*sdk.Result, error) {
	r, _, err := suite.FullSendAToB(msg)
	if success {
		suite.Require().NoError(err, "IBC send failed. Expected success. %s", err)
	} else {
		suite.Require().Error(err, "IBC send succeeded. Expected failure")

	}
	return r, err
}

// Tests
// Test that Sending IBC messages works when the middleware isn't configured
func (suite *MiddlewareTestSuite) TestSendTransferNoContract() {
	one := sdk.NewInt(1)
	suite.AssertSend(true, suite.MessageFromAToB(sdk.DefaultBondDenom, one))
}

// Test that Receiving IBC messages works when the middleware isn't configured
func (suite *MiddlewareTestSuite) TestReceiveTransferNoContract() {
	one := sdk.NewInt(1)
	suite.AssertReceive(true, suite.MessageFromBToA(sdk.DefaultBondDenom, one))
}

func (suite *MiddlewareTestSuite) fullSendTest(native bool) {
	// Get the denom and amount to send
	denom := sdk.DefaultBondDenom
	//channel := "channel-0"
	if !native {
		denomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom("transfer", "channel-0", denom))
		fmt.Println(denomTrace)
		denom = denomTrace.IBCDenom()
	}

	// Setup contract
	suite.chainA.StoreContractCode(&suite.Suite, "./bytecode/address_limiter.wasm")
}

// Test address limiting on sends
func (suite *MiddlewareTestSuite) TestSendTransferWithAddressLimitingNative() {
	// Sends denom=stake from A->B. Address limit receives "stake" in the packet. Nothing to do in the contract
	suite.fullSendTest(true)
}

func (suite *MiddlewareTestSuite) TestUnsetAddressLimitingContract() {
	// Setup contract
	suite.chainA.StoreContractCode(&suite.Suite, "./bytecode/address_limiter.wasm")
	addr := suite.chainA.InstantiateALContract(&suite.Suite, []sdk.Address{})
	suite.chainA.RegisterAddressLimitingContract(addr)

	// Unset the contract param
	params, err := types.NewParams("")
	suite.Require().NoError(err)
	app := suite.chainA.GetApp()
	paramSpace, ok := app.ParamsKeeper.GetSubspace(types.ModuleName)
	suite.Require().True(ok)
	// N.B.: this panics if validation fails.
	paramSpace.SetParamSet(suite.chainA.GetContext(), &params)
}
*/
