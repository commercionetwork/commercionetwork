package commercionetworkibctesting


import (
	"fmt"
	"io/ioutil"

	"github.com/stretchr/testify/require"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	"github.com/commercionetwork/commercionetwork/x/ibc-address-limiter/types"
	"github.com/stretchr/testify/suite"
)

func (chain *TestChain) StoreContractCode(suite *suite.Suite, path string) {
	commercionetworkApp := chain.TestChain.GetSimApp()

	govKeeper := commercionetworkApp.GovKeeper
	wasmCode, err := ioutil.ReadFile(path)
	suite.Require().NoError(err)

	addr := commercionetworkApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)
	src := wasmtypes.StoreCodeProposalFixture(func(p *wasmtypes.StoreCodeProposal) {
		p.RunAs = addr.String()
		p.WASMByteCode = wasmCode
	})

	// when stored
	storedProposal, err := govKeeper.SubmitProposal(chain.GetContext(), src)
	suite.Require().NoError(err)

	// and proposal execute
	handler := govKeeper.Router().GetRoute(storedProposal.ProposalRoute())
	err = handler(chain.GetContext(), storedProposal.GetContent())
	suite.Require().NoError(err)
}

func (chain *TestChain) InstantiateALContract(suite *suite.Suite, addrs_whitelist []sdk.Address) sdk.AccAddress {
	commercionetworkApp := chain.GetApp()
	transferModule := commercionetworkApp.AccountKeeper.GetModuleAddress(transfertypes.ModuleName)
	govModule := commercionetworkApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)

	initMsgBz := []byte(fmt.Sprintf(`{
           "gov_module":  "%s",
           "ibc_module":"%s",
           "addrs_whitelist": %s
        }`,
		govModule, transferModule, addrs_whitelist))

	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(commercionetworkApp.WasmKeeper)
	codeID := uint64(1)
	creator := commercionetworkApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)
	addr, _, err := contractKeeper.Instantiate(chain.GetContext(), codeID, creator, creator, initMsgBz, "address limiting contract", nil)
	suite.Require().NoError(err)
	return addr
}

func (chain *TestChain) InstantiateContract(suite *suite.Suite, msg string) sdk.AccAddress {
	commercionetworkApp := chain.GetApp()
	contractKeeper := wasmkeeper.NewDefaultPermissionKeeper(commercionetworkApp.WasmKeeper)
	codeID := uint64(1)
	creator := commercionetworkApp.AccountKeeper.GetModuleAddress(govtypes.ModuleName)
	addr, _, err := contractKeeper.Instantiate(chain.GetContext(), codeID, creator, creator, []byte(msg), "contract", nil)
	suite.Require().NoError(err)
	return addr
}

func (chain *TestChain) QueryContract(suite *suite.Suite, contract sdk.AccAddress, key []byte) string {
	commercionetworkApp := chain.GetApp()
	state, err := commercionetworkApp.WasmKeeper.QuerySmart(chain.GetContext(), contract, key)
	suite.Require().NoError(err)
	return string(state)
}

func (chain *TestChain) RegisterAddressLimitingContract(addr []byte) {
	addrStr, err := sdk.Bech32ifyAddressBytes("did:com:", addr)
	require.NoError(chain.T, err)
	params, err := types.NewParams(addrStr)
	require.NoError(chain.T, err)
	commercionetworkApp := chain.GetApp()
	paramSpace, ok := commercionetworkApp.ParamsKeeper.GetSubspace(types.ModuleName)
	require.True(chain.T, ok)
	paramSpace.SetParamSet(chain.GetContext(), &params)
}

