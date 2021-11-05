package config

import (
	"net/url"

	"github.com/rs/zerolog/log"
)

// isValidUrl tests a string to determine if it is a valid URL or not
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	log.Debug().Str("remoteConfig", toTest).Msg("Valid URL for remote config.")
	return true
}
