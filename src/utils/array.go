package utils

func Includes(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func IndexOf(data []string,element string) (int) {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1   
 }