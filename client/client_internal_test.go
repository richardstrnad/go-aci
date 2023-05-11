package client

import (
	"net/http"
	"testing"
)

const host = "apic.test.com"
const username = "admin"

func TestClientBasic(t *testing.T) {
	client := New(
		host,
		username,
		WithPassword("password"),
	)
	t.Run("Test baseURL", func(t *testing.T) {
		got := client.baseURL.String()
		want := "https://apic.test.com"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("Test Verify correctly set", func(t *testing.T) {
		got := client.httpClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify
		want := false
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})

}

func TestClientHttp(t *testing.T) {
	client := New(
		host,
		username,
		WithPassword("password"),
		WithHttp(),
	)
	t.Run("Test HTTP", func(t *testing.T) {
		got := client.scheme
		want := "http"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestClientNoVerify(t *testing.T) {
	client := New(
		host,
		username,
		WithPassword("password"),
		WithNoVerify(),
	)
	t.Run("Test NoVerify", func(t *testing.T) {
		got := client.verify
		want := false
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})
	t.Run("Test NoVerify correctly set", func(t *testing.T) {
		got := client.httpClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify
		want := true
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})
}
