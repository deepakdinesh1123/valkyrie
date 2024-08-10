package docker

import (
	"github.com/docker/docker/client"
)

// newClient creates a new Docker client using environment variables and API version negotiation.
//
// Parameters:
// None
// Returns:
// - *client.Client: The newly created Docker client.
// - error: Any error that occurred during client creation.
func newClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
