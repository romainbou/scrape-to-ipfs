package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func srapeAndServe(url string, c *gin.Context) {
	log.Print("scrape")
	main, allLinks := Scraper.Scrape(url)

	log.Print("finished scrapting")

	for _, link := range allLinks {
		addFileToIPFS(link)
	}

	hash := addFileToIPFS(main)

	redirectToGateway(hash, c)
}

func serveFile(filename string, c *gin.Context) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		// c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
	reader := bytes.NewReader(file)
	contentLength := int64(len(file))
	contentType := "text/html;"

	extraHeaders := map[string]string{}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

}

func redirectToGateway(hash string, c *gin.Context) {
	if gatewayURL, ok := os.LookupEnv("GATEWAY_URL"); ok {
		c.Redirect(http.StatusMovedPermanently, gatewayURL+hash)
	} else {
		c.Redirect(http.StatusMovedPermanently, "https://gateway.ipfs.io/ipfs/"+hash)
	}
}

func addFileToIPFS(content string) string {

	ipfsClient := IpfsApi.NewShell("localhost:5001")
	if ipfsDaemonURL, ok := os.LookupEnv("IPFS_DAEMON"); ok {
		ipfsClient = IpfsApi.NewShell(ipfsDaemonURL)
	}
	reader := strings.NewReader(content)
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
