package lib

import (
	"testing"
)

func TestScrape(t *testing.T) {

	helloURL := "https://ipfs.io/ipfs/QmZ4tDuvesekSs4qM5ZBKpXiZGun7S2CYtEZRB3DYXkjGx"
	content := Scrape(helloURL)
	if content != "hello worlds\n" {
		t.Errorf("URL: %s should contains 'hello worlds\\n', received: %s", helloURL, content)
	}

	emptyURL := "https://ipfs.io/ipfs/QmZ4tDuvesekSs4qM5ZBKpXiZGun7S2CYtEZRB3DYXkjGx"
	content := Scrape(helloURL)
	if content != "hello worlds\n" {
		t.Errorf("URL: %s should contains 'hello worlds\\n', received: %s", helloURL, content)
	}

}
