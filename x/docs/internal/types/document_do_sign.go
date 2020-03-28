package types

type DocumentDoSign struct {
	StorageUri         string `json:"storage_uri"`
	SignerInstance     string `json:"signer_instance"`
	VcrId              string `json:"vcr_id"`
	CertificateProfile string `json:"certificate_profile"`
}
