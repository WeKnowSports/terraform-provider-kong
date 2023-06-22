package helper

func ConvertInterfaceArrToStrings(strs []interface{}) []string {
	arr := make([]string, len(strs))
	for i, str := range strs {
		arr[i] = str.(string)
	}
	return arr
}

func SetObjectID(ID string) map[string]string {
	results := map[string]string{}

	if ID == "" {
		return nil
	}

	results["id"] = ID
	return results
}

func SetConsumerID(ID, username string) map[string]string {
	results := map[string]string{}

	if ID == "" && username == "" {
		return nil
	} else if ID != "" {
		results["id"] = ID
	} else if username != "" {
		results["username"] = username
	}

	return results
}
