package v2_0_0

import (
	"github.com/cosmos/cosmos-sdk/codec"
	v038auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v0_38"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

func Migrate(appState genutil.AppMap) genutil.AppMap {
	oldAccountsCdc := codec.New()
	codec.RegisterCrypto(oldAccountsCdc)
	RegisterCodec(oldAccountsCdc)

	v038Codec := codec.New()
	codec.RegisterCrypto(v038Codec)
	v038auth.RegisterCodec(v038Codec)

	if appState[v038auth.ModuleName] != nil {
		var old GenesisState

		// deserialize current state in our old state serialization representation
		oldAccountsCdc.MustUnmarshalJSON(appState[v038auth.ModuleName], &old)

		// create a slice with len(old.Accounts) slots to accommodate 0.38-formatted genesis accounts
		newAccounts := make(v038auth.GenesisAccounts, len(old.Accounts))

		for i, account := range old.Accounts {
			// since we can encounter either a BaseAccount or ModuleAccount,
			// convert account to interface and do type assertion magic to determine what type we're dealing with
			accIface := interface{}(account)

			var newGenAcc v038auth.GenesisAccount

			switch v := accIface.(type) {
			case *BaseAccount: // this is a BaseAccount
				seq, accn := secAccnUint64(v)
				nl := v038auth.NewBaseAccount(v.GetAddress(), v.GetCoins(), getPk(v.PubKey), accn, seq)
				newGenAcc = v038auth.GenesisAccount(nl)
			case *ModuleAccount: // this is a ModuleAccount
				seq, accn := secAccnUint64(v.BaseAccount)
				nl := v038auth.NewBaseAccount(v.GetAddress(), v.GetCoins(), getPk(v.PubKey), accn, seq)
				newGenAcc = v038auth.NewModuleAccount(nl, v.Name, v.Permissions...)
			}

			newAccounts[i] = newGenAcc
		}

		newAccounts = sanitizeGenesisAccounts(newAccounts)

		if err := validateGenAccounts(newAccounts); err != nil {
			panic(err)
		}

		// delete old state
		delete(appState, v038auth.ModuleName)

		// append new state with old params and new accounts
		appState[v038auth.ModuleName] = v038Codec.MustMarshalJSON(
			v038auth.NewGenesisState(old.Params, newAccounts),
		)
	}

	return appState
}
