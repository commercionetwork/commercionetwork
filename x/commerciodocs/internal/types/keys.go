package types

const (
	ModuleName = "commerciodocs"

	// Key of the map { DocumentReference => Address }
	OwnersStoreKey = "docs_owners"
	// Key of the map { DocumentReference => Metadata }
	MetadataStoreKey = "docs_metadata"
	// Key of the map { Did => []Sharing }
	SharingStoreKey = "docs_sharing"

	QuerierRoute = ModuleName
)
