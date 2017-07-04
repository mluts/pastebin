package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	contentType = "application/x-www-form-urlencoded"
)

// Client ...
type Client struct {
	url      string
	insecure bool
}

// NewClient ...
func NewClient(url string, insecure bool) *Client {
	return &Client{url: url, insecure: insecure}
}

// Paste ...
func (c *Client) Paste(body io.Reader) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !c.insecure},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Post(c.url, contentType, body)
	if err != nil {
		log.Printf("error pasting to %s: %s", c.url, err)
		return err
	}

	if res.StatusCode != 200 {
		log.Printf("unexpected response from %s: %d", c.url, res.StatusCode)
		return errors.New("unexpected response")
	}

	fmt.Printf("%s", res.Request.URL.String())

	return nil
}
