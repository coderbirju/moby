package client // import "github.com/docker/docker/client"

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types/container"
)

// ContainerUpdate updates the resources of a container.
func (cli *Client) ContainerUpdate(ctx context.Context, containerID string, updateConfig container.UpdateConfig) (container.UpdateResponse, error) {
	containerID, err := trimID("container", containerID)
	if err != nil {
		return container.UpdateResponse{}, err
	}

	serverResp, err := cli.post(ctx, "/containers/"+containerID+"/update", nil, updateConfig, nil)
	defer ensureReaderClosed(serverResp)
	if err != nil {
		return container.UpdateResponse{}, err
	}

	var response container.UpdateResponse
	err = json.NewDecoder(serverResp.body).Decode(&response)
	return response, err
}
