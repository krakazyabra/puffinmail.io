package generateemail

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"

// EmailGenerator keeps track of hosts and min length needed by email generator methods
type EmailGenerator struct {
	hosts []string
	l     int
}

// NewEmailGenerator returns an email generator that creates emails with the given hosts and minimum length user part
func NewEmailGenerator(h []string, l int) *EmailGenerator {
	rand.Seed(time.Now().UTC().UnixNano())
	return &EmailGenerator{hosts: h, l: l}
}

// GetHosts returns all available hosts
func (eg *EmailGenerator) GetHosts() []string {
	return eg.hosts
}

// HostsContains tells whether host h is in eg.hosts
func (eg *EmailGenerator) HostsContains(h string) bool {
	for _, n := range eg.hosts {
		if h == n {
			return true
		}
	}
	return false
}

// NewRandom generates a new random email address. It is the callers responsibility to check for uniqueness
func (eg *EmailGenerator) NewRandom() string {
	a := []byte(alphabet)
	name := make([]byte, eg.l)

	for i := range name {
		name[i] = a[rand.Intn(len(a))]
	}

	domain := eg.hosts[rand.Intn(len(eg.hosts))]

	return string(name) + "@" + domain
}

// NewFromRouteAndHost generates a new email address from a string and host. It is the callers responsibility to check for uniqueness
func (eg *EmailGenerator) NewFromRouteAndHost(r string, h string) (string, error) {
	if eg.HostsContains(h) {
		return string(r) + "@" + h, nil
	}
	return "", fmt.Errorf("invalid host: %s", h)
}

var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

//VerifyRoute verifies the local part of an email address is between 3 and 64 alphanumeric characters
func (eg *EmailGenerator) VerifyRoute(r string) error {
	if len(r) < 3 {
		return fmt.Errorf("route must be at least three characters: %s", r)
	} else if len(r) > 64 {
		return fmt.Errorf("route must be fewer than 64 characters: %s", r)
	} else if !isAlphaNumeric(r) {
		return fmt.Errorf("route may only contain letters (a-z, A-Z) and numbers (0-9): %s", r)
	} else if r == "webmaster" || r == "admin" || r == "postmaster" || r == "administrator" || r == "root" {
		return fmt.Errorf("route is blacklisted: %s", r)
	}
	return nil
}

//VerifyHost verifies the host part of an email address is not empty and is known to the application
func (eg *EmailGenerator) VerifyHost(h string) error {
	if h == "" {
		return errors.New("host must not be an empty string")
	} else if !eg.HostsContains(h) {
		return fmt.Errorf("host not in list of known hosts: %s", h)
	}
	return nil
}
