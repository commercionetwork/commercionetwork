package client

import (
	commercioauthcmd "commercio-network/x/commercioauth/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group CommercioAUTH queries under a subcommand
	authCmd := &cobra.Command{
		Use:   "commercioauth",
		Short: "CommercioAUTH querying commands",
	}

	authCmd.AddCommand(client.GetCommands(
		commercioauthcmd.GetCmdReadAccount(mc.storeKey, mc.cdc),
	)...)

	return authCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:   "commercioauth",
		Short: "CommercioAUTH transactions subcommands",
	}

	govTxCmd.AddCommand(client.PostCommands(
		commercioauthcmd.GetCmdRegisterAccount(mc.cdc),
	)...)

	return govTxCmd
}
