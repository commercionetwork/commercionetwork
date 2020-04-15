package cli

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"

	"github.com/commercionetwork/commercionetwork/x/docs/types"
)

func TestCliTx(t *testing.T) {
	sender := "cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm"
	recipient, uuid := "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0", "ac33043b-5cb4-4645-a3f9-819140847252"
	contentURIMetadata := "http://thecontentmetadata.com"
	schemaURI, schemaVersion := "theSchemaUri", "theSchemaVersion"
	contentURI := "http://contenturi.com"
	checksumValue, checksumAlgo := "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8", "sha-1"

	type flagValue struct {
		flag  string
		value interface{}
	}

	tests := []struct {
		name        string
		expectedErr error
		params      []string
		flags       []flagValue
	}{
		{
			"happy path",
			fmt.Errorf("no RPC client defined"), // It means it tries to broadcast, so it works.
			[]string{
				recipient,
				uuid,
				contentURI,
				schemaURI,
				schemaVersion,
			},
			[]flagValue{},
		},
		{
			"happy path with do_sign",
			fmt.Errorf("no RPC client defined"), // It means it tries to broadcast, so it works.
			[]string{
				recipient,
				uuid,
				contentURIMetadata,
				schemaURI,
				schemaVersion,
				contentURI,
				checksumValue,
				checksumAlgo,
			},
			[]flagValue{},
		},
		{
			"with invalid sign sdn data",
			fmt.Errorf("sdn_data value \"invalid\" is not supported"), // It means it tries to broadcast, so it works.
			[]string{
				recipient,
				uuid,
				contentURIMetadata,
				schemaURI,
				schemaVersion,
				contentURI,
				checksumValue,
				checksumAlgo,
			},
			[]flagValue{
				{FlagSign, true},
				{FlagSignStorageURI, "http://theSignStorageURI.com"},
				{FlagSignSignerInstance, "theSignerInstance"},
				{FlagSignVcrID, "theVcrId"},
				{FlagSignCertificateProfile, "theCertificateProfile"},
				{FlagSignSdnData, fmt.Sprintf("%s,%s", types.SdnDataCommonName, "invalid")},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cdc := amino.NewCodec()
			cmd := getCmdShareDocument(cdc)

			viper.Set(flags.FlagFrom, sender)
			viper.Set(flags.FlagGenerateOnly, true)

			for _, fl := range tt.flags {
				viper.Set(fl.flag, fl.value)
			}

			err := cmd.RunE(cmd, tt.params)

			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}
