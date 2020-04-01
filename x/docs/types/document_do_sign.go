package types

import (
	"fmt"
	"strings"
)

const (
	SdnDataCommonName   = "common_name"
	SdnDataSurname      = "surname"
	SdnDataSerialNumber = "serial_number"
	SdnDataGivenName    = "given_name"
	SdnDataOrganization = "organization"
	SdnDataCountry      = "country"

	InputStringSep = ","
)

var validSdnData = map[string]struct{}{
	SdnDataCommonName:   {},
	SdnDataSurname:      {},
	SdnDataSerialNumber: {},
	SdnDataGivenName:    {},
	SdnDataOrganization: {},
	SdnDataCountry:      {},
}

// DocumentDoSign represents the optional DoSign value inside a Document.
type DocumentDoSign struct {
	StorageURI         string  `json:"storage_uri"`
	SignerInstance     string  `json:"signer_instance"`
	SdnData            SdnData `json:"sdn_data"`
	VcrID              string  `json:"vcr_id"`
	CertificateProfile string  `json:"certificate_profile"`
}

// SdnData represents the SdnData value inside a DocumentDoSign struct.
type SdnData []string

// Validate checks that the SdnData is valid, only accepts value included in
// validSdnData.
func (s SdnData) Validate() error {
	for _, val := range s {
		if _, ok := validSdnData[val]; !ok {
			return fmt.Errorf("sdn_data value \"%s\" is not supported", val)
		}
	}

	return nil
}

// NewSdnDataFromString generates a SdnData struct based on the input string.
// The input string expects a comma-separated value as:
// "common_name,surname,serial_number"
// If empty string is provided, a SdnData with default value will be provided. Default : "serial_number".
func NewSdnDataFromString(input string) (SdnData, error) {
	if input == "" {
		return SdnData{SdnDataSerialNumber}, nil
	}

	var split SdnData = strings.Split(input, InputStringSep)
	err := split.Validate()
	if err != nil {
		return SdnData{}, err
	}

	return split, nil
}
