package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

// Account represents a GitHub account with its access token.
type Account struct {
	Alias string
	Token string
}

// MultiClient manages multiple GitHub account clients.
type MultiClient struct {
	clients map[string]*github.Client
}

// NewMultiClient creates clients for each configured account.
func NewMultiClient(accounts []Account) *MultiClient {
	mc := &MultiClient{clients: make(map[string]*github.Client)}
	for _, acc := range accounts {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: acc.Token})
		tc := oauth2.NewClient(context.Background(), ts)
		mc.clients[acc.Alias] = github.NewClient(tc)
	}
	return mc
}

// For returns the GitHub client for a given account alias.
func (mc *MultiClient) For(alias string) (*github.Client, error) {
	c, ok := mc.clients[alias]
	if !ok {
		return nil, fmt.Errorf("account '%s' not configured", alias)
	}
	return c, nil
}
