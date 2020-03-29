package types

import (
	"fmt"
	"strings"
)

type DocumentDoSign struct {
	StorageURI         string  `json:"storage_uri"`
	SignerInstance     string  `json:"signer_instance"`
	SdnData            SdnData `json:"sdn_data"`
	VcrID              string  `json:"vcr_id"`
	CertificateProfile string  `json:"certificate_profile"`
}

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

type SdnData []string

func (s SdnData) Validate() error {
	for _, val := range s {
		if _, ok := validSdnData[val]; !ok {
			return fmt.Errorf("sdn_data value \"%s\" is not supported", val)
		}
	}

	return nil
}

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
