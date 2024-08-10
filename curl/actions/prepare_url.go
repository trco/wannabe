package actions

func PrepareUrl(scheme string, host string, path string, query string) string {
	if path == "/" {
		path = ""
	}
	return scheme + "://" + host + path + query
}
