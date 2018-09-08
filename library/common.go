package library

type Common struct {
}

func (t *Common) InArray(str string, arr []string) bool {
	for _, value := range arr {
		if value == str {
			return true
		}
	}
	return false
}
