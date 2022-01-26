package celcat

import "net/url"

type SessionData struct {
	token        string
	federationId string
	location     url.URL
}
