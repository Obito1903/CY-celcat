package celcat

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getRequestVerificationToken(client *http.Client, url url.URL) string {
	var token string = ""
	// Send a request to the celcat app
	resp, err := client.Get(url.String() + "/LdapLogin")
	if err != nil {
		log.Fatal("Could not get the Request Verification Token from ", url, err)
		os.Exit(1)
	}

	// Init the scanner
	scanner := bufio.NewScanner(resp.Body)
	// Scan the response for the __RequestVerificationToken element
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "__RequestVerificationToken") {
			requestTokenTag := scanner.Text()
			posStartValue := strings.Index(requestTokenTag, "value=")
			posEndValue := strings.Index(requestTokenTag, "\" />")
			token = requestTokenTag[posStartValue+7 : posEndValue]
			break
		}
	}
	return token
}

// Log in the specified user into the given client
func Login(client *http.Client, url url.URL, username string, password string) {
	getRequestVerificationToken(client, url)
}
