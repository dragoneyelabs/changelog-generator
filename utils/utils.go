package utils

import (
	"net/url"
	"strings"
)

// InArray checks if the string is in the given slice
func InArray(s string, arr []string) bool {
	for _, e := range arr {
		if e == s {
			return true
		}
	}
	return false
}

// GetPrefix returns the prefix of the string from a list of valid prefixes if found, otherwise returns defaultPrefix
func GetPrefix(s string, validPrefixes []string, defaultPrefix string) string {
	for _, p := range validPrefixes {
		if strings.HasPrefix(s, p) {
			return p
		}
	}
	return defaultPrefix
}

// TrimPrefix returns s without the provided leading prefix string, case insensitive. If s doesn't start with prefix, s
// is returned unchanged.
func TrimPrefix(s, prefix string) string {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// HasPrefix tests whether the string s begins with prefix, case insensitive
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

// Capitalize capitalize the first letter of the string
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	l := strings.ToUpper(string(s[0]))
	if len(s) == 1 {
		return l
	}
	return l + s[1:]
}

// GetHost returns the host from the uri
func GetHost(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

// GetPath returns the path from the uri
func GetPath(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return u.Path, nil
}
