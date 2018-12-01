package lib

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Scrape downloads the content of the given URL and returns it as a string
func Scrape(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Fatal(bodyString)
	}

	// TODO detect if it's a web page to recursively scape the assets and replace the links with IPFS

	bodyBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)

	return bodyString
}
