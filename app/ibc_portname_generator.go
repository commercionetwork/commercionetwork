package app

/*package app

import (
	"strings"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CommercioIBCPortNameGenerator uses simple translation to generate port-id from contract address.
type CommercioIBCPortNameGenerator struct{}

var AccountAddressPrefixTranslation = "did.com."

// PortIDForContract coverts contract into port-id in the format "wasm.did.com.<contract-address>"
func (CommercioIBCPortNameGenerator) PortIDForContract(ctx sdk.Context, addr sdk.AccAddress) string {
	addrStr := addr.String()
	if strings.HasPrefix(addrStr, AccountAddressPrefix) {
		addrStr = AccountAddressPrefixTranslation + addrStr[len(AccountAddressPrefix):]
	}
	return wasmkeeper.GetPortIDPrefix() + addrStr
}

// ContractFromPortID reads the contract address from "wasm.did.com.<contract-address>" in the port-id.
func (CommercioIBCPortNameGenerator) ContractFromPortID(ctx sdk.Context, portID string) (sdk.AccAddress, error) {
	if !strings.HasPrefix(portID, wasmkeeper.GetPortIDPrefix()) {
		return nil, sdkerrors.Wrapf(wasmtypes.ErrInvalid, "without prefix")
	}
	portIDaddr := portID[len(wasmkeeper.GetPortIDPrefix()):]
	if strings.HasPrefix(portIDaddr, AccountAddressPrefixTranslation) {
		portIDaddr = AccountAddressPrefix + portIDaddr[len(AccountAddressPrefixTranslation):]
	}
	return sdk.AccAddressFromBech32(portIDaddr)
}
*/
