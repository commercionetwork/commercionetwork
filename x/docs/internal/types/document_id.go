package types

// DocumentIds represents a list of documents' UUIDs
type DocumentIds []string

// AppendIfMissing allows to add to this list of documentIds
// the given id, if it isn't already present
// TODO: Test this
func (documentIds DocumentIds) AppendIfMissing(id string) DocumentIds {
	for _, ele := range documentIds {
		if ele == id {
			return documentIds
		}
	}
	return append(documentIds, id)
}
