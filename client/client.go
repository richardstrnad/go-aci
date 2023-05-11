package client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/richardstrnad/go-aci/models"
)

var (
	ErrMissingCookie = errors.New("missing cookie in response")
)

type Client struct {
	Host       string
	Port       int
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

// WithPort sets the port
func WithPort(port int) Option {
	return func(c *Client) {
		c.Port = port
	}
}

// New returns a new Client
func New(host, username string, options ...Option) *Client {
	client := &Client{
		Host:     host,
		Port:     443,
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
		Scheme: c.scheme,
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
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
	// We set the cookie jar to use the default cookie jar
	c.httpClient.Jar, _ = cookiejar.New(nil)
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
	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed: got status %s", resp.Status)
	}
	// We want to make sure that the cookie jar has a cookie in it
	cookieName := "APIC-cookie"
	cookies := c.httpClient.Jar.Cookies(c.baseURL)
	if len(cookies) == 0 {
		return ErrMissingCookie
	}
	if cookies[0].Name != cookieName {
		return ErrMissingCookie
	}
	return nil
}
