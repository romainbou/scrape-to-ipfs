package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Scrape downloads the content of the given URL and returns it as a string
func Scrape(URL string, depth int) (io.Reader, []string) {
	response, err := http.Get(URL)
	allLinks := []string{}
	if err != nil {
		if depth == 1 {
			log.Fatal(err)
		}
		return nil, allLinks
	}

	if response.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		if depth == 1 {
			log.Fatal(bodyString)
		}
		return nil, allLinks
	}

	bodyReader := response.Body

	if depth == 1 {
		var buf bytes.Buffer
		tee := io.TeeReader(bodyReader, &buf)
		allLinks = findLinkedAsset(tee)
		return &buf, allLinks
	}
	return bodyReader, allLinks

}

// ReplaceLinks replaces the links by their hash in the content
func ReplaceLinks(reader io.Reader, links []string, hashes []string) string {

	gatewayURL := "https://gateway.ipfs.io/ipfs/"

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	content := string(bodyBytes)

	for key, link := range links {
		fmt.Println("Replace ", link)
		replacer := strings.NewReplacer(link, gatewayURL+hashes[key])
		content = replacer.Replace(content)
	}
	return content
}

func getTokenLink(token html.Token) (linkValue string, err error) {
	if token.Data == "link" {
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
			linksArray = appendTokenLink(t, linksArray)
		case html.SelfClosingTagToken:
			t := z.Token()
			linksArray = appendTokenLink(t, linksArray)
		}
	}
}

func appendTokenLink(token html.Token, linksArray []string) []string {
	linkValue, err := getTokenLink(token)
	if err == nil {
		// Link found
		linksArray = append(linksArray, linkValue)
	}
	return linksArray
}

// PrependBaseURL add the base URL if the scheme is empty
func PrependBaseURL(link string, baseURL string) string {
	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}
	if u.Scheme == "" {
		return baseURL + link
	}
	return link
}
