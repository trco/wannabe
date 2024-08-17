package actions

func FilterHeaders(headers map[string][]string, exclude []string) map[string][]string {
	filteredHeaders := make(map[string][]string)

	for key, values := range headers {
		if !contains(exclude, key) {
			filteredHeaders[key] = values
		}
	}

	return filteredHeaders
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
