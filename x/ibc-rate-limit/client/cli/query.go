package cli

import (
	"github.com/spf13/cobra"

	//"github.com/osmosis-labs/osmosis/v13/osmoutils/osmocli"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/client/cli"
	"github.com/commercionetwork/commercionetwork/x/ibc-rate-limit/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := cli.GetQueryCmd(types.ModuleName)

	/*cmd.AddCommand(
		cli.GetParams[*types.QueryParamsRequest](
			types.ModuleName, types.NewQueryClient),
	)*/

	return cmd
}
