package jwt

import (
	"fmt"
	"github.com/nikogura/jwt-ssh-agent-go/pkg/agentjwt"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// ExtractDomain extracts the domain from a URL-like string.
func ExtractDomain(urlLikeString string) (domain string, err error) {
	urlLikeString = strings.TrimSpace(urlLikeString)

	if regexp.MustCompile(`^https?`).MatchString(urlLikeString) {
		read, _ := url.Parse(urlLikeString)
		urlLikeString = read.Host
	}

	if regexp.MustCompile(`^www\.`).MatchString(urlLikeString) {
		urlLikeString = regexp.MustCompile(`^www\.`).ReplaceAllString(urlLikeString, "")
	}

	domain = regexp.MustCompile(`([a-z0-9\-]+\.)*[a-z0-9\-]+`).FindString(urlLikeString)
	if domain == "" {
		err = errors.New(fmt.Sprintf("failed parsing domain from %s", urlLikeString))
		return domain, err
	}

	return domain, err
}

// LoadPubKey loads a public key from a file path.
func LoadPubKey(path string) (key string, err error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		err = errors.Wrapf(err, "failed reading %s", path)
		return key, err
	}

	key = string(keyBytes)
	key = strings.TrimRight(key, "\n")

	return key, err
}

// MakeToken creates a JWT token using the provided URL, username, and public key.
func MakeToken(url string, username string, pubkey string) (token string, err error) {
	domain, err := ExtractDomain(url)
	if err != nil {
		err = errors.Wrapf(err, "unparsable url")
		return token, err
	}

	// Make JWT
	token, err = agentjwt.SignedJwtToken(username, domain, pubkey)
	if err != nil {
		err = errors.Wrap(err, "failed to create signed token")
		return token, err
	}

	return token, err
}