package types

type DocumentDoSign struct {
	StorageUri         string  `json:"storage_uri"`
	SignerInstance     string  `json:"signer_instance"`
	SdnData            SdnData `json:"sdn_data"`
	VcrId              string  `json:"vcr_id"`
	CertificateProfile string  `json:"certificate_profile"`
}

type SdnData struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Tin          string `json:"tin"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
	Country      string `json:"country"`
}
