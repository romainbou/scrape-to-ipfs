package lib

import (
	"log"
	"net/url"
	"strings"
)

// IsValidHTTPURL Verifies that the given string is an HTTP of HTTPS URL
func IsValidHTTPURL(rawURL string) bool {
	URL, err := url.Parse(rawURL)

	if err != nil {
		log.Print(err)
		return false
	}

	if URL.Scheme != "http" && URL.Scheme != "https" {
		return false
	}

	if !strings.Contains(URL.Host, ".") {
		return false
	}

	return true
}
