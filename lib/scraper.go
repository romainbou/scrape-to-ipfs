package lib

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

// Scrape downloads the content of the given URL and returns it as a string
func Scrape(url string) (string, []string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	allLinks := []string{}

	if response.StatusCode != http.StatusOK {
		allLinks = findLinkedAsset(response.Body)
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

	// Find asset (css, fonts, JS, images, videos, 3D files, iframe) URLS
	// Add all of them to IPFS
	// Replace them with the Gateway link

	return bodyString, allLinks
}

func getTokenLink(token html.Token) (linkValue string, err error) {
	if token.Data == "a" || token.Data == "link" {
		for _, a := range token.Attr {
			if a.Key == "href" {
				return a.Val, nil
			}
		}
	}
	if token.Data == "img" || token.Data == "script" {
		for _, a := range token.Attr {
			if a.Key == "src" {
				return a.Val, nil
			}
		}
	}
	return "", errors.New("No value found in the link")
}

func findLinkedAsset(reader io.Reader) []string {

	z := html.NewTokenizer(reader)

	linksArray := []string{}

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// End of doc
			return linksArray
		case html.StartTagToken:
			t := z.Token()
			linkValue, err := getTokenLink(t)
			if err == nil {
				// Link found
				linksArray = append(linksArray, linkValue)
			}
		}
	}
}
