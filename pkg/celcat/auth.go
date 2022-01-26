package celcat

import (
	"bufio"
	"fmt"
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
	fmt.Println("1")
	session.location = celcatUrl

	// Init the header to be sent
	formData := url.Values{
		"Name":                       {username},
		"Password":                   {password},
		"__RequestVerificationToken": {session.token},
	}
	fmt.Println("2")

	// Setup a login request
	req, err := http.NewRequest("POST", celcatUrl.String()+"/LdapLogin/Logon", strings.NewReader(formData.Encode()))
	if err != nil {
		log.Fatal("Could not create request.", err)
		os.Exit(1)
	}
	fmt.Println("3")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the login request
	resp, err := client.Do(req)
	fmt.Println("4")

	if err != nil {
		log.Fatal("Could not login to : ", celcatUrl.String(), err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Could not login to : ", celcatUrl.String())
		os.Exit(1)
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	responseUrl, err := resp.Request.Response.Location()
	fmt.Println("6")

	if err != nil {
		log.Fatal("Could not parse response.", err)
		os.Exit(1)
	}
	fmt.Println("7")

	session.federationId = responseUrl.Query().Get("FederationIds")
	return session
}
