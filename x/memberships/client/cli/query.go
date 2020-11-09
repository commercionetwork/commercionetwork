package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/memberships/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		getCmdGetInvites(cdc),
		getCmdGetTrustedServiceProviders(cdc),
		getCmdGetPoolFunds(cdc),
		getCmdMembershipForUser(cdc),
		getCmdMemberships(cdc),
		getCmdTspMemberships(cdc),
	)

	return cmd
}

func getCmdGetInvites(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "invites",
		Short: "Get all membership invitations",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getInvitesFunc(cmd, cdc)
		},
	}
}

func getInvitesFunc(cmd *cobra.Command, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetInvites)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return err
	}

	cmd.Println(string(res))

	return nil
}

func getCmdGetTrustedServiceProviders(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "trusted-service-providers",
		Short: "Get all membership invitations for a user",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTrustedServiceProvidersFunc(cmd, cdc)
		},
	}
}

func getTrustedServiceProvidersFunc(cmd *cobra.Command, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetTrustedServiceProviders)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return err
	}

	cmd.Println(string(res))

	return nil
}

func getCmdGetPoolFunds(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool-funds",
		Short: "Get the pool funds amounts",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdGetPoolFundsFunc(cmd, cdc)
		},
	}
}

func getCmdGetPoolFundsFunc(cmd *cobra.Command, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetPoolFunds)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return fmt.Errorf("could not get pool funds schemes: %w", err)
	}

	cmd.Println(string(res))

	return nil
}

func getCmdMembershipForUser(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "membership [user]",
		Short: "Get user membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdMembershipForUserFunc(cmd, args, cdc)
		},
	}
}

func getCmdMembershipForUserFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetMembership, args[0])
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return fmt.Errorf("could not get membership for user: %w", err)
	}

	cmd.Println(string(res))

	return nil
}

func getCmdMemberships(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "memberships",
		Short: "Get all memberships",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdMembershipsFunc(cmd, cdc)
		},
	}
}

func getCmdMembershipsFunc(cmd *cobra.Command, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGetMemberships)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return fmt.Errorf("could not get memberships: %w", err)
	}

	cmd.Println(string(res))

	return nil
}

func getCmdTspMemberships(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sold [tsp-address]",
		Short: "Get tsp-address memberships",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdTspMembershipsFunc(cmd, args, cdc)
		},
	}
}

func getCmdTspMembershipsFunc(cmd *cobra.Command, args []string, cdc *codec.Codec) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetTspMemberships, args[0])
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return fmt.Errorf("could not get memberships for tsp: %w", err)
	}

	cmd.Println(string(res))

	return nil
}
