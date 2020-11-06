package utils

func Includes(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}
func Splice(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func IndexOf(data []string, element string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}
