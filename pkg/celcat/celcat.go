package celcat

import "net/url"

type SessionData struct {
	token         string
	federationIds string
	location      url.URL
}
