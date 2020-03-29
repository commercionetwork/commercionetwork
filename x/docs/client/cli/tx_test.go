package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

func TestCliTx_HappyPath(t *testing.T) {
	cdc := amino.NewCodec()

	cmd := getCmdShareDocument(cdc)

	sender := "cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm"
	recipient, uuid := "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0", "ac33043b-5cb4-4645-a3f9-819140847252"
	contentUri := "http://thecontent.com"
	schemaURI, schemaVersion := "theSchemaUri", "theSchemaVersion"

	viper.Set(flags.FlagFrom, sender)
	viper.Set(flags.FlagGenerateOnly, true)

	err := cmd.RunE(cmd, []string{
		recipient,
		uuid,
		contentUri,
		schemaURI,
		schemaVersion,
	})

	// It tries to broadcast, so it means it passes.
	require.EqualError(t, err, "no RPC client defined")
}

func TestGetTxCmd_WithDoSign(t *testing.T) {
	cmd := getCmdShareDocument(amino.NewCodec())

	sender := "cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm"
	recipient, uuid := "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0", "ac33043b-5cb4-4645-a3f9-819140847252"
	contentUriMetadata := "http://thecontentmetadata.com"
	schemaURI, schemaVersion := "theSchemaUri", "theSchemaVersion"
	contentUri := "http://contenturi.com"
	checksumValue, checksumAlgo := "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8", "sha-1"

	// Do Sign Params
	viper.Set(FlagSign, true)
	viper.Set(FlagSignStorageURI, "http://theSignStorageURI.com")
	viper.Set(FlagSignSignerInstance, "theSignerInstance")
	viper.Set(FlagSignVcrID, "theVcrId")
	viper.Set(FlagSignCertificateProfile, "theCertificateProfile")
	viper.Set(FlagSignSdnDataFirstName, "firstName")
	viper.Set(FlagSignSdnDataLastName, "lastName")
	viper.Set(FlagSignSdnDataTin, "tinData")
	viper.Set(FlagSignSdnDataEmail, "theEmail")
	viper.Set(FlagSignSdnDataOrganization, "theOrganization")
	viper.Set(FlagSignSdnDataCountry, "theCountry")

	viper.Set(flags.FlagFrom, sender)
	viper.Set(flags.FlagGenerateOnly, true)

	err := cmd.RunE(cmd, []string{
		recipient,
		uuid,
		contentUriMetadata,
		schemaURI,
		schemaVersion,
		contentUri,
		checksumValue,
		checksumAlgo,
	})

	// It tries to broadcast, so it means it passes.
	require.EqualError(t, err, "no RPC client defined")
}
