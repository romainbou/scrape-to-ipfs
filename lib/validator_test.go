package lib

import (
	"testing"
)

func TestIsValidHTTPURL(t *testing.T) {

	validURLs := []string{
		"http://example.com",
		"https://example.com",
	}

	for _, validURL := range validURLs {
		isValid := IsValidHTTPURL(validURL)
		if !isValid {
			t.Errorf("URL: %s should be valid", validURL)
		}
	}

	invalidURLs := []string{
		"example.com",
		"https://com",
		"htp://example.com",
		"ftp://example.com",
		"ipfs://example.com",
		"https:/example.com",
		"https//example.com",
	}

	for _, invalidURL := range invalidURLs {
		isValid := IsValidHTTPURL(invalidURL)
		if isValid {
			t.Errorf("URL: %s shouldn't be valid", invalidURL)
		}
	}

}
