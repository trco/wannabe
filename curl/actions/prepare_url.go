package actions

func PrepareUrl(host string, path string, query string) string {
	return host + path + query
}
