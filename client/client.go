package client

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
	"net/url"

	"github.com/richardstrnad/go-aci/models"
)

type Client struct {
	Host       string
	Username   string
	Password   string
	scheme     string
	httpClient *http.Client
	baseURL    *url.URL
	verify     bool
}

type Option func(*Client)

func WithPassword(password string) Option {
	return func(c *Client) {
		c.Password = password
	}
}

// WithHttp sets the scheme to http
func WithHttp() Option {
	return func(c *Client) {
		c.scheme = "http"
	}
}

// WithNoVerify disables TLS certificate verification
func WithNoVerify() Option {
	return func(c *Client) {
		c.verify = false
	}
}

// New returns a new Client
func New(host, username string, options ...Option) *Client {
	client := &Client{
		Host:     host,
		Username: username,
		scheme:   "https",
		verify:   true,
	}
	for _, option := range options {
		option(client)
	}
	// We run a bunch of internal functions to set up the client
	client.setBaseURL()
	client.setHttpClient()
	return client
}

func (c *Client) setBaseURL() {
	c.baseURL = &url.URL{
		Scheme: "https",
		Host:   c.Host,
	}
}

func (c *Client) setHttpClient() {
	c.httpClient = &http.Client{}
	// We set the transport to use TLS and to skip verification if verify is false
	c.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !c.verify,
		},
	}
}

func (c *Client) Login() error {
	u, err := url.JoinPath(c.baseURL.String(), "api", "aaaLogin.json")
	if err != nil {
		return err
	}
	data := models.NewLogin(c.Username, c.Password)
	j, err := data.ToJSON()
	resp, err := c.httpClient.Post(u, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	log.Print(resp.Status)
	for _, cookie := range resp.Cookies() {
		log.Print(cookie.Name)
		log.Print(cookie.Value)
	}
	return nil
}
