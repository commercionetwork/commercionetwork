package v1_5_0

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

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

// getPk converts our Pubkey old serialization format into a proper Tendermint crypto.Pubkey.
func getPk(p PubKey) crypto.PubKey {
	var pk crypto.PubKey

	if p.Value != "" {
		// decode public key base64
		pkHex, err := base64.StdEncoding.DecodeString(p.Value)
		if err != nil {
			panic(err)
		}

		// copy pkHex into castPk for type coherence
		castPk := secp256k1.PubKeySecp256k1{}
		copy(castPk[:], pkHex)

		// transform castPk into a bech32 pubkey...
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, castPk)
		if err != nil {
			panic(err)
		}

		// ...then finally transform the bech32 into a crypto.PublicKey instance
		pk, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, pks)
		if err != nil {
			panic(err)
		}
	}

	return pk
}

// secAccnUint64 converts sequence number and account number strings to their uint64 values.
func secAccnUint64(b *BaseAccount) (accNum uint64, seq uint64) {
	var err error

	accNum, err = strconv.ParseUint(b.AccountNumber, 10, 64)
	if err != nil {
		panic(fmt.Errorf("could not convert account_number to uint64: %w", err))
	}

	seq, err = strconv.ParseUint(b.Sequence, 10, 64)
	if err != nil {
		panic(fmt.Errorf("could not convert sequence to uint64: %w", err))
	}

	return
}
