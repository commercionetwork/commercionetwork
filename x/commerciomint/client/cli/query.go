package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group commerciomint queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1

	cmd.AddCommand(
		CmdGetEtps(),
		CmdGetAllEtps(),
		CmdGetEtp(),

		//CmdGetConversionRate(), // TODO DELTE COMMAND
		//CmdGetFreezePeriod(), // TODO DELTE COMMAND
		CmdGetParams(),
	)

	return cmd
}
