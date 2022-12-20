package helper

func ConvertInterfaceArrToStrings(strs []interface{}) []string {
	arr := make([]string, len(strs))
	for i, str := range strs {
		arr[i] = str.(string)
	}
	return arr
}
