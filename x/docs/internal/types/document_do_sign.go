package types

type DocumentDoSign struct {
	StorageUri         string  `json:"storage_uri"`
	SignerInstance     string  `json:"signer_instance"`
	SdnData            SdnData `json:"sdn_data"`
	VcrId              string  `json:"vcr_id"`
	CertificateProfile string  `json:"certificate_profile"`
}

type SdnData struct {
	FirstName    string
	LastName     string
	Tin          string
	Email        string
	Organization string
	Country      string
}
