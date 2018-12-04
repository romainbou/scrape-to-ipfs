package lib

import (
	"bytes"
	"testing"
)

func TestScrape(t *testing.T) {

	helloURL := "https://ipfs.io/ipfs/QmZ4tDuvesekSs4qM5ZBKpXiZGun7S2CYtEZRB3DYXkjGx"
	content, allLinks := Scrape(helloURL, 1)
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
					<img src="https://golang.org/doc/gopher/frontpage.png" />
				</body>
			</html>
			`

	reader := bytes.NewReader([]byte(html))

	links := findLinkedAsset(reader)

	if len(links) != 3 {
		t.Errorf("The HTML should contains 3 links, %d received", len(links))
	}

}
func TestFindAssetLinks2(t *testing.T) {

	html := `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Example</title>
        <meta name="author" content="Example" />
    </head>
    <body>
        <h1>Example</h1>
        <p>
            <a href="https://twitter.com/golang" title="Twitter"><img src="img/twitter.svg" alt="Twitter" title="Twitter" /></a>
            <a href="https://github.com/romainbou" title="Github"><img  src="img/github.svg" alt="Github" title="Github" /></a>
        </p>
    </body>
    <style>
    body {
        text-align: center;
        font-family: 'Lato', Tahoma, 'Trebuchet MS', serif;
    }
    </style>
</html>
			`

	reader := bytes.NewReader([]byte(html))

	links := findLinkedAsset(reader)

	if len(links) != 2 {
		t.Errorf("The HTML should contains 2 links, %d received", len(links))
	}

}
