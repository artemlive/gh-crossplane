package github

import (
	"github.com/cli/go-gh/v2/pkg/api"
)

// GetGitHubClient returns a GitHub API client using the default authentication method.
func GetGitHubClient() (*api.RESTClient, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, err
	}

	return client, nil
}
