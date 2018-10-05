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
)

func main() {

	r := gin.Default()
	//r.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.LoadHTMLGlob("templates/*")
	// r.GET("/", func(c *gin.Context) {
	// 	indexHandler(c)
	// })

	r.GET("/*trail", func(c *gin.Context) {
		argHandler(c)
	})

	r.Run(":8000")
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
	outputFilename := "output.html"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create output file
	outFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Copy data from HTTP response to file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return outputFilename
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
	c.Redirect(http.StatusMovedPermanently, "https://gateway.ipfs.io/ipfs/"+hash)
}

func addFileToIPFS(filename string) string {
	ipfsClient := IpfsApi.NewShell("localhost:5001")
	file, err := ioutil.ReadFile(filename)
	reader := bytes.NewReader(file)
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
