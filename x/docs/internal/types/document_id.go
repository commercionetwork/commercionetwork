package types

// DocumentIds represents a list of documents' UUIDs
type DocumentIds []string

// AppendIfMissing allows to add to this list of documentIds
// the given id, if it isn't already present
func (documentIds DocumentIds) AppendIfMissing(id string) (DocumentIds, bool) {
	for _, ele := range documentIds {
		if ele == id {
			return documentIds, false
		}
	}
	return append(documentIds, id), true
}
