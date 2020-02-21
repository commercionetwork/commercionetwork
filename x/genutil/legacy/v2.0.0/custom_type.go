package v2_0_0

// DONTCOVER
// nolint

// PSA: this file only serves as a bridge between Cosmos SDK we used before v0.38 (we used master, commit 92ea174ea6e6),
// and should NEVER be used to instantiate new accounts or stuff like that.

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v038auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v0_38"
	"github.com/tendermint/tendermint/crypto"

	v034auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v0_34"
)

// PubKey is a type used as a bridge between an old serialization type used by Cosmos/Tendermint.
// It doesn't implement the crypto.PubKey interface, it's only needed to make JSON unmarshaling work.
type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type (
	// partial interface needed only for amino encoding and sanitization
	Account interface {
		GetAddress() sdk.AccAddress
		GetAccountNumber() string
		GetCoins() sdk.Coins
		SetCoins(sdk.Coins) error
	}

	GenesisAccount interface {
		Account

		Validate() error
	}

	GenesisAccounts []GenesisAccount

	GenesisState struct {
		Params   v034auth.Params `json:"params" yaml:"params"`
		Accounts GenesisAccounts `json:"accounts" yaml:"accounts"`
	}

	BaseAccount struct {
		Address       sdk.AccAddress `json:"address" yaml:"address"`
		Coins         sdk.Coins      `json:"coins" yaml:"coins"`
		PubKey        PubKey         `json:"public_key" yaml:"public_key"`
		AccountNumber string         `json:"account_number" yaml:"account_number"`
		Sequence      string         `json:"sequence" yaml:"sequence"`
	}

	baseAccountPretty struct {
		Address       sdk.AccAddress `json:"address" yaml:"address"`
		Coins         sdk.Coins      `json:"coins" yaml:"coins"`
		PubKey        PubKey         `json:"public_key" yaml:"public_key"`
		AccountNumber string         `json:"account_number" yaml:"account_number"`
		Sequence      string         `json:"sequence" yaml:"sequence"`
	}

	ModuleAccount struct {
		*BaseAccount

		Name        string   `json:"name" yaml:"name"`
		Permissions []string `json:"permissions" yaml:"permissions"`
	}

	moduleAccountPretty struct {
		BaseAccount BaseAccount
		Name        string   `json:"name" yaml:"name"`
		Permissions []string `json:"permissions" yaml:"permissions"`
	}
)

func NewBaseAccount(
	address sdk.AccAddress, coins sdk.Coins, pk PubKey, accountNumber, sequence string,
) *BaseAccount {

	return &BaseAccount{
		Address:       address,
		Coins:         coins,
		PubKey:        pk,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
}

func (acc BaseAccount) GetAddress() sdk.AccAddress {
	return acc.Address
}

func (acc *BaseAccount) GetAccountNumber() string {
	return acc.AccountNumber
}

func (acc *BaseAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

func (acc *BaseAccount) SetCoins(coins sdk.Coins) error {
	acc.Coins = coins
	return nil
}

func (acc BaseAccount) Validate() error {

	return nil
}

func (acc BaseAccount) MarshalJSON() ([]byte, error) {
	alias := baseAccountPretty{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}

	return json.Marshal(alias)
}

// UnmarshalJSON unmarshals raw JSON bytes into a BaseAccount.
func (acc *BaseAccount) UnmarshalJSON(bz []byte) error {
	var alias baseAccountPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		fmt.Println(string(bz))
		return err
	}

	acc.PubKey = alias.PubKey
	acc.Address = alias.Address
	acc.Coins = alias.Coins
	acc.AccountNumber = alias.AccountNumber
	acc.Sequence = alias.Sequence

	return nil
}

func (ma ModuleAccount) Validate() error {
	if err := validatePermissions(ma.Permissions...); err != nil {
		return err
	}

	if strings.TrimSpace(ma.Name) == "" {
		return errors.New("module account name cannot be blank")
	}

	if !ma.Address.Equals(sdk.AccAddress(crypto.AddressHash([]byte(ma.Name)))) {
		return fmt.Errorf("address %s cannot be derived from the module name '%s'", ma.Address, ma.Name)
	}

	return ma.BaseAccount.Validate()
}

// MarshalJSON returns the JSON representation of a ModuleAccount.
func (ma ModuleAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(moduleAccountPretty{
		Name:        ma.Name,
		Permissions: ma.Permissions,
	})
}

// UnmarshalJSON unmarshals raw JSON bytes into a ModuleAccount.
func (ma *ModuleAccount) UnmarshalJSON(bz []byte) error {
	var alias moduleAccountPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	ma.BaseAccount = NewBaseAccount(alias.BaseAccount.Address, alias.BaseAccount.Coins, PubKey{}, alias.BaseAccount.AccountNumber, alias.BaseAccount.Sequence)
	ma.Name = alias.Name
	ma.Permissions = alias.Permissions

	return nil
}

func validatePermissions(permissions ...string) error {
	for _, perm := range permissions {
		if strings.TrimSpace(perm) == "" {
			return fmt.Errorf("module permission is empty")
		}
	}

	return nil
}

func sanitizeGenesisAccounts(genAccounts v038auth.GenesisAccounts) v038auth.GenesisAccounts {
	sort.Slice(genAccounts, func(i, j int) bool {
		return genAccounts[i].GetAccountNumber() < genAccounts[j].GetAccountNumber()
	})

	for _, acc := range genAccounts {
		if err := acc.SetCoins(acc.GetCoins().Sort()); err != nil {
			panic(err)
		}
	}

	return genAccounts
}

func validateGenAccounts(genAccounts v038auth.GenesisAccounts) error {
	addrMap := make(map[string]bool, len(genAccounts))
	for _, acc := range genAccounts {

		// check for duplicated accounts
		addrStr := acc.GetAddress().String()
		if _, ok := addrMap[addrStr]; ok {
			return fmt.Errorf("duplicate account found in genesis state; address: %s", addrStr)
		}

		addrMap[addrStr] = true

		// check account specific validation
		if err := acc.Validate(); err != nil {
			return fmt.Errorf("invalid account found in genesis state; address: %s, error: %s", addrStr, err.Error())
		}
	}

	return nil
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "cosmos-sdk/Account", nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "cosmos-sdk/ModuleAccount", nil)
}
