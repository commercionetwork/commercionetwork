package cli

import (
	"fmt"
    "strconv"
	"github.com/spf13/cobra"

    "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
)

var _ = strconv.Itoa(0)

func CmdCurrentUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Get the currently active upgrade",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCurrentUpgradeFunc(cmd, args)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

    return cmd
}

func getCurrentUpgradeFunc(cmd *cobra.Command, _ []string) error {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrent)
	res, _, err := clientCtx.QueryWithData(route, nil)

	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}