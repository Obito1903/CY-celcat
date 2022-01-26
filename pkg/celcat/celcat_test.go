package celcat

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
)

func TestRequestToken(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		os.Exit(1)
	}

	client := &http.Client{
		Jar: jar,
	}
	url, err := url.Parse("https://services-web.u-cergy.fr/calendar")
	if err != nil {
		os.Exit(1)
	}
	t.Log("Token: " + getRequestVerificationToken(client, *url))
}
