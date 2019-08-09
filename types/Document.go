package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type DocumentMetadataSchema struct {
	Uri     string `json:"uri"`
	Version string `json:"version"`
}

type DocumentMetadata struct {
	ContentUri string                 `json:"content_uri"`
	Schema     DocumentMetadataSchema `json:"schema"`
	Proof      string                 `json:"proof"`
}

type DocumentChecksum struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
}

type Document struct {
	Sender     sdk.AccAddress   `json:"sender"`
	Recipient  sdk.AccAddress   `json:"recipient"`
	ContentUri string           `json:"content_uri"`
	Metadata   DocumentMetadata `json:"metadata"`
	Checksum   DocumentChecksum `json:"checksum"`
}
