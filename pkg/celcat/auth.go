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
		log.Fatal("Could not get the Request Verification Token from ", url.String(), err)
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
func Login(client *http.Client, celcatUrl url.URL, username string, password string) SessionData {
	var session SessionData
	session.token = getRequestVerificationToken(client, celcatUrl)
	session.location = celcatUrl

	// Init the header to be sent
	formData := url.Values{
		"Name":                       {username},
		"Password":                   {password},
		"__RequestVerificationToken": {session.token},
	}

	// Setup a login request
	req, err := http.NewRequest("POST", celcatUrl.String()+"/LdapLogin/Logon", strings.NewReader(formData.Encode()))
	if err != nil {
		log.Fatal("Could not create request.", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the login request
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 || len(client.Jar.Cookies(&celcatUrl)) < 2 {
		log.Fatal("Could not login to, check your login and password : ", celcatUrl.String())
		os.Exit(1)
	}
	defer resp.Body.Close()
	responseUrl, err := resp.Request.Response.Location()
	session.federationId = responseUrl.Query().Get("FederationIds")
	return session
}
