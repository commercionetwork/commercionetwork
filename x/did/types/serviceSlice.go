package types

type ServiceSlice []*Service

func (slice ServiceSlice) hasDuplicate() bool {
	if len(slice) == 0 {
		return false
	}

	for _, s := range slice[:len(slice)-1] {
		if contains(s.ID, slice[1:]) {
			return true
		}
	}

	return false
}

func contains(str string, slice []*Service) bool {
	for _, c := range slice[1:] {
		if str == c.ID {
			return true
		}
	}

	return false
}
