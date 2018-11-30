package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	IpfsApi "github.com/ipfs/go-ipfs-api"
)

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

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

	if isValidHTTPURL(url) {
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
	if isValidHTTPURL(trimmedURL) {
		log.Print("call scape and server", trimmedURL)
		srapeAndServe(trimmedURL, c)
	} else {
		indexHandler(c)
	}
}

func srapeAndServe(url string, c *gin.Context) {
	log.Print("scrape")
	filename := scrape(url)
	log.Print("finished scrapting")

	hash := addFileToIPFS(filename)

	redirectToGateway(hash, c)

	// serveFile(filename, c)
}

func scrape(url string) string {
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

func isValidHTTPURL(rawURL string) bool {
	URL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		log.Print(err)
		return false
	} else if strings.HasPrefix(URL.Scheme, "http") {
		return true
	} else {
		return false
	}
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
