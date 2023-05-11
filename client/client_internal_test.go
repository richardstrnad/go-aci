package client

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

const host = "apic.test.com"
const username = "admin"

func TestClientBasic(t *testing.T) {
	t.Run("Test baseURL", func(t *testing.T) {
		client := New(
			host,
			username,
			WithPassword("password"),
		)
		got := client.baseURL.String()
		want := "https://apic.test.com:443"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("Test Verify correctly set", func(t *testing.T) {
		client := New(
			host,
			username,
			WithPassword("password"),
		)
		got := client.httpClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify
		want := false
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})
	t.Run("Test port set correctly", func(t *testing.T) {
		client := New(
			host,
			username,
			WithPassword("password"),
			WithPort(8443),
		)
		got := client.baseURL.String()
		want := "https://apic.test.com:8443"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
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

func TestClientLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "APIC-cookie",
			Value:    "token",
			HttpOnly: true,
			Path:     "/",
		})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"imdata":[{"aaaLogin":{"attributes":{"token":"token"}}}]}`))
	}))

	host := strings.Split(strings.Split(server.URL, ":")[1], "//")[1]
	port, err := strconv.Atoi(strings.Split(server.URL, ":")[2])
	if err != nil {
		t.Error(err)
	}
	client := New(
		host,
		username,
		WithPassword("password"),
		WithHttp(),
		WithPort(port),
	)

	t.Run("Test Login", func(t *testing.T) {
		err := client.Login()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestClientLoginFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"imdata":[{"aaaLogin":{"attributes":{"token":"token"}}}]}`))
	}))

	host := strings.Split(strings.Split(server.URL, ":")[1], "//")[1]
	port, err := strconv.Atoi(strings.Split(server.URL, ":")[2])
	if err != nil {
		t.Error(err)
	}
	client := New(
		host,
		username,
		WithPassword("password"),
		WithHttp(),
		WithPort(port),
	)

	t.Run("Test Login fail missing cookie", func(t *testing.T) {
		err := client.Login()
		if err != ErrMissingCookie {
			t.Errorf("got %q, want %q", err, ErrMissingCookie)
		}
	})
}
