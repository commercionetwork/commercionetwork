package types

type ServiceSlice []*Service

func (slice ServiceSlice) hasDuplicate() bool {
	if len(slice) == 0 {
		return false
	}

	for i, s := range slice[:len(slice)-1] {
		i++
		if contains(s.ID, slice[i:]) {
			return true
		}
	}

	return false
}

func contains(str string, slice []*Service) bool {
	for _, c := range slice {
		if str == c.ID {
			return true
		}
	}

	return false
}
