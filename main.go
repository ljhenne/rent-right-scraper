package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	allowedDomain = os.Getenv("ALLOWED_DOMAIN")
	certFile = flag.String("cert", "./cert.pem", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "./private.pem", "A PEM encoded private key file.")
	url = os.Getenv("URL_TO_VISIT")
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}


func main()  {
	fmt.Println("Starting script...")

	// Load client cert
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	check(err)
	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
	)
	c.WithTransport(&http.Transport{TLSClientConfig: tlsConfig})
	c.OnHTML("#search-results", func (e *colly.HTMLElement) {
		numChildren := e.DOM.Children().Length()
		fmt.Printf("Found %s child elements... gonna save the file now", strconv.Itoa(numChildren))

		err := ioutil.WriteFile("site.html", e.Response.Body, 0600)
		check(err)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ... ", r.URL.String())
	})

	err = c.Visit(url)
	check(err)
}
