package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateSdnData(t *testing.T) {
	tests := []struct {
		name        string
		sdnData     SdnData
		isValid     bool
		expectedErr error
	}{
		{
			"invalid: invalid field empty",
			SdnData{
				SdnDataCommonName,
				SdnDataSurname,
				"",
			},
			false,
			fmt.Errorf("sdn_data value \"%s\" is not supported", ""),
		},
		{
			"invalid: strange value not included in supportedSdnData",
			SdnData{
				SdnDataCommonName,
				SdnDataSurname,
				"age",
			},
			false,
			fmt.Errorf("sdn_data value \"%s\" is not supported", "age"),
		},
		{
			"valid: all supported fields",
			SdnData{
				SdnDataCommonName,
				SdnDataSurname,
				SdnDataSurname,
				SdnDataGivenName,
				SdnDataOrganization,
				SdnDataCountry,
			},
			true,
			nil,
		},
		{
			"valid: subset of supported fields",
			SdnData{
				SdnDataCommonName,
				SdnDataSurname,
				SdnDataGivenName,
				SdnDataCountry,
			},
			true,
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdnData.Validate()

			if tt.isValid {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func TestNewSdnDataFromString(t *testing.T) {
	tests := []struct {
		name            string
		inputString     string
		shouldFail      bool
		expectedErr     error
		expectedSdnData SdnData
	}{
		{
			"empty string returns a default value",
			"",
			false,
			nil,
			SdnData{
				"serial_number",
			},
		},
		{
			"valid fields included",
			fmt.Sprintf("%s,%s,%s,%s", SdnDataGivenName, SdnDataSurname, SdnDataOrganization, SdnDataSerialNumber),
			false,
			nil,
			SdnData{
				SdnDataGivenName,
				SdnDataSurname,
				SdnDataOrganization,
				SdnDataSerialNumber,
			},
		},
		{
			"invalid fields included",
			fmt.Sprintf("%s,%s,%s,%s", SdnDataGivenName, "invalid", SdnDataOrganization, SdnDataSerialNumber),
			true,
			fmt.Errorf("sdn_data value \"invalid\" is not supported"),
			SdnData{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			sdnData, err := NewSdnDataFromString(tt.inputString)
			if tt.shouldFail {
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.Equal(t, tt.expectedSdnData, sdnData)
			}
		})
	}
}
