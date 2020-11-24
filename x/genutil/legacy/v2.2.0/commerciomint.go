package v2_2_0

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/commercionetwork/commercionetwork/x/commerciomint"
	commerciomintTypes "github.com/commercionetwork/commercionetwork/x/commerciomint/types"

	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func commercioMintMigrate(appState genutil.AppMap, govAddress sdk.AccAddress) genutil.AppMap {
	// Instead of initializing our own codec, use the cosmos x/auth one
	// Because we must trigger x/supply/internal/types init() function,
	// call NewEmptyModuleAccount.
	// Calling a function method triggers its init() function if it haven't been triggered already.
	// All because otherwise, we could not unmarshal cosmos-sdk/ModuleAccount accounts.
	// Sad, isn't it?
	_ = supply.NewEmptyModuleAccount("fu")
	accountsCdc := authTypes.ModuleCdc

	// 1. remove uccc from every account, and move commerciomint ucommercio's to government address
	if appState[auth.ModuleName] != nil {
		var old auth.GenesisState
		accountsCdc.MustUnmarshalJSON(appState[auth.ModuleName], &old)

		newAccounts := make(authexported.GenesisAccounts, 0, len(old.Accounts))

		var mintUcomm sdk.Coins
		for _, account := range old.Accounts {
			coins := account.GetCoins()
			comAmount := coins.AmountOf("ucommercio")

			account.SetCoins(
				sdk.NewCoins(
					sdk.NewCoin("ucommercio", comAmount),
				),
			)

			accIface := interface{}(account)

			if macc, ok := accIface.(*supply.ModuleAccount); ok {

				if macc.GetName() == "accreditations" {
					continue
				}
				if mintUcomm == nil {
					if macc.GetName() == "commerciomint" {
						mintUcomm = macc.GetCoins()

						// reset commerciomint coins
						err := account.SetCoins(sdk.NewCoins())
						if err != nil {
							panic(err)
						}
					}

				}
			}

			newAccounts = append(newAccounts, account)
		}

		if mintUcomm != nil {
			for i := 0; i < len(newAccounts); i++ {
				if newAccounts[i].GetAddress().Equals(govAddress) {
					newAccounts[i].SetCoins(
						newAccounts[i].GetCoins().Add(mintUcomm...),
					)
					break
				}
			}
		}

		// delete old state
		delete(appState, auth.ModuleName)

		// append new state with old params and new accounts
		appState[auth.ModuleName] = accountsCdc.MustMarshalJSON(
			auth.NewGenesisState(old.Params, newAccounts),
		)
	}

	// 2. commerciomint state must change completely according to its new state schema
	mintCdc := codec.New()
	codec.RegisterCrypto(mintCdc)
	commerciomintTypes.RegisterCodec(mintCdc)

	if appState[commerciomintTypes.ModuleName] != nil {
		delete(appState, commerciomintTypes.ModuleName)

		appState[commerciomintTypes.ModuleName] = mintCdc.MustMarshalJSON(
			commerciomint.DefaultGenesisState(),
		)
	}

	// 3. staking module must only contain ucommercio
	if appState[supply.ModuleName] != nil {
		var old supply.GenesisState
		supply.ModuleCdc.MustUnmarshalJSON(appState[supply.ModuleName], &old)

		commAmount := old.Supply.AmountOf("ucommercio")

		old.Supply = sdk.NewCoins(sdk.NewCoin("ucommercio", commAmount))

		delete(appState, supply.ModuleName)
		appState[supply.ModuleName] = supply.ModuleCdc.MustMarshalJSON(
			old,
		)
	}

	return appState
}
