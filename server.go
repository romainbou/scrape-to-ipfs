package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	IpfsApi "github.com/ipfs/go-ipfs-api"

	Scraper "github.com/romainbou/scrape-to-ipfs/lib"
	Validator "github.com/romainbou/scrape-to-ipfs/lib"
)

func main() {

	r := gin.Default()

	templatesFolder := "templates/"
	if tplFolder, ok := os.LookupEnv("TEMPLATES"); ok {
		templatesFolder = tplFolder
	}

	r.LoadHTMLGlob(templatesFolder + "*")

	r.GET("/*trail", func(c *gin.Context) {
		argHandler(c)
	})

	if port, ok := os.LookupEnv("PORT"); ok {
		r.Run(":" + port)
	} else {
		r.Run(":8000")
	}

}

func indexHandler(c *gin.Context) {

	URLParam := "url"

	url := c.DefaultQuery(URLParam, "")

	// TODO: Warn if the IPFS server is unavailable

	if Validator.IsValidHTTPURL(url) {
		srapeAndServe(url, c)
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":    "IPFS forwarder",
		"URLParam": URLParam,
	})
}

func argHandler(c *gin.Context) {
	url := c.Param("trail")
	trimmedURL := url[1:]
	log.Print("call is valid: ", trimmedURL)
	if Validator.IsValidHTTPURL(trimmedURL) {
		log.Print("call scape and server", trimmedURL)
		srapeAndServe(trimmedURL, c)
	} else {
		indexHandler(c)
	}
}

func srapeAndServe(URL string, c *gin.Context) {
	log.Print("scrape")

	u, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	baseURL := u.Scheme + "://" + u.Host + "/"
	main, allLinks := Scraper.Scrape(URL, 1)

	log.Print("finished scrapting")

	linkHashes := []string{}
	for _, link := range allLinks {
		fmt.Println("Link: ", Scraper.PrependBaseURL(link, baseURL))
		fmt.Println("Link: ", link)
		linkContent, _ := Scraper.Scrape(Scraper.PrependBaseURL(link, baseURL), 2)
		hash := addFileToIPFS(linkContent)
		linkHashes = append(linkHashes, hash)
	}

	replacedMain := Scraper.ReplaceLinks(main, allLinks, linkHashes)

	reader := strings.NewReader(replacedMain)
	hash := addFileToIPFS(reader)

	redirectToGateway(hash, c)
	return
}

func serveFile(filename string, c *gin.Context) {
	file, _ := ioutil.ReadFile(filename)
	reader := bytes.NewReader(file)
	contentLength := int64(len(file))
	contentType := "text/html;"

	extraHeaders := map[string]string{}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

}

func redirectToGateway(hash string, c *gin.Context) {
	if gatewayURL, ok := os.LookupEnv("GATEWAY_URL"); ok {
		c.Redirect(http.StatusTemporaryRedirect, gatewayURL+hash)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "https://gateway.ipfs.io/ipfs/"+hash)
	}
	return
}

func addFileToIPFS(reader io.Reader) string {

	ipfsClient := IpfsApi.NewShell("localhost:5001")
	if ipfsDaemonURL, ok := os.LookupEnv("IPFS_DAEMON"); ok {
		ipfsClient = IpfsApi.NewShell(ipfsDaemonURL)
	}
	hash, err := ipfsClient.Add(reader)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("added", hash)
	return hash
}
