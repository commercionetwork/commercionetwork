package client

import (
	commerciodocscmd "commercio-network/x/commerciodocs/client/cli"
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
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:   "commerciodocs",
		Short: "CommercioDOCS querying commands",
	}

	govQueryCmd.AddCommand(client.GetCommands(
		commerciodocscmd.GetCmdReadDocumentMetadata(mc.storeKey, mc.cdc),
		commerciodocscmd.GetCmdListAuthorizedReaders(mc.storeKey, mc.cdc),
	)...)

	return govQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:   "commerciodocs",
		Short: "CommercioDOCS transactions subcommands",
	}

	govTxCmd.AddCommand(client.PostCommands(
		commerciodocscmd.GetCmdStoreDocument(mc.cdc),
		commerciodocscmd.GetCmdShareDocument(mc.cdc),
	)...)

	return govTxCmd
}
