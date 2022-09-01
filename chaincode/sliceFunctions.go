package chaincode

// Go is not making development faster...
// Big brains force grug to write collection functions because idealism is cool
// Big brains also fight against generics until 1.18. Waste many grugs time.
func removeStringItem(slice []string, s string) ([]string, bool) {
	foundAt := -1
	for i, v := range slice {
		if v == s {
			foundAt = i
		}
	}

	if foundAt == -1 {
		return slice, false
	} else {
		return removeStringIndex(slice, foundAt), true
	}
}

func removeStringIndex(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
