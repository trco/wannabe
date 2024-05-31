package actions

func PrepareUrl(host string, path string, query string) string {
	if path == "/" {
		path = ""
	}
	return host + path + query
}
