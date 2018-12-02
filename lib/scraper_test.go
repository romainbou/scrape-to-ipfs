package lib

import (
	"bytes"
	"testing"
)

func TestScrape(t *testing.T) {

	helloURL := "https://ipfs.io/ipfs/QmZ4tDuvesekSs4qM5ZBKpXiZGun7S2CYtEZRB3DYXkjGx"
	content, allLinks := Scrape(helloURL)
	if content != "hello worlds\n" {
		t.Errorf("URL: %s should contains 'hello worlds\\n', received: %s", helloURL, content)
	}
	if len(allLinks) != 0 {
		t.Errorf("URL: %s should contains 0 links received: %d", helloURL, len(allLinks))
	}

}

func TestFindAssetLinks(t *testing.T) {

	html := `
			<!DOCTYPE html>
			<html lang="en-GB">
				<head>
					<link rel="stylesheet" href="https://golang.org/lib/godoc/style.css">
					<script src="https://golang.org/lib/godoc/godocs.js" defer></script>
				<head>
				<body>
					<a href="https://golang.org/doc">Documents</a>
				</body>
			</html>
			`

	reader := bytes.NewReader([]byte(html))

	links := findLinkedAsset(reader)

	if len(links) != 3 {
		t.Errorf("The HTML should contains 3 links, %d received", len(links))
	}

}
