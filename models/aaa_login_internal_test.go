package models

import (
	"fmt"
	"testing"
)

const name = "admin"
const pwd = "password"

func TestJSON(t *testing.T) {
	login := NewLogin(name, pwd)
	t.Run("Test baseURL", func(t *testing.T) {
		got, err := login.ToJSON()
		if err != nil {
			t.Errorf("got error %q", err)
		}
		want := fmt.Sprintf(`{"aaaUser":{"attributes":{"name":"%s","pwd":"%s"}}}`, name, pwd)
		if string(got) != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
