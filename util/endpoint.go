package util

import (
	"net/url"
	"strings"
)

func IsValidEndpoint(endpoint string) bool {
	if !strings.HasPrefix(endpoint, "https://") {
		return false
	}

	// check if the endpoint is a valid URL
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return false
	}

	return true
}
