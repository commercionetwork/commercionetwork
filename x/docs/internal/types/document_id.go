package types

// DocumentIDs represents a list of documents' UUIDs
type DocumentIDs []string

// AppendIfMissing allows to add to this list of documentIds
// the given id, if it isn't already present
func (documentIds DocumentIDs) AppendIfMissing(id string) (DocumentIDs, bool) {
	for _, ele := range documentIds {
		if ele == id {
			return documentIds, false
		}
	}
	return append(documentIds, id), true
}

// Empty tells if the list of document ids is empty or not
func (documentIds DocumentIDs) Empty() bool {
	return len(documentIds) == 0
}
