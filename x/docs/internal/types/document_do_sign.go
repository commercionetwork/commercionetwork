package types

import "fmt"

type DocumentDoSign struct {
	StorageURI         string  `json:"storage_uri"`
	SignerInstance     string  `json:"signer_instance"`
	SdnData            SdnData `json:"sdn_data"`
	VcrID              string  `json:"vcr_id"`
	CertificateProfile string  `json:"certificate_profile"`
}

var validSdnData = map[string]struct{}{
	"common_name":   {},
	"surname":       {},
	"serial_number": {},
	"given_name":    {},
	"organization":  {},
	"country":       {},
}

type SdnData []string

func (s SdnData) Validate() error {
	for _, val := range s {
		if _, ok := validSdnData[val]; !ok {
			return fmt.Errorf("sdn_data value %s is not supported", val)
		}
	}

	return nil
}
