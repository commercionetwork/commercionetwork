package main

import (
	"os"

	"github.com/commercionetwork/commercionetwork/app"
	"github.com/commercionetwork/commercionetwork/cmd/commercionetworkd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
