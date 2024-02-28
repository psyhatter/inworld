package inworld

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// All types in this file are generated based on this documentation:
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge

// CheckDeploymentStatus retrieves the current status of a long-running
// operation. Clients can use it to periodically poll for the operation's
// result, following the interval recommendations provided by the API service.
// This approach is useful for tracking the progress and completion status of
// operations that are not immediately finalized.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes/#check-deployment-status
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#check-deployment-status
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge/#check-deployment-status
func (c Client) CheckDeploymentStatus(ctx context.Context, operationID string) (CheckDeploymentStatusResponse, error) {
	if operationID == "" {
		return CheckDeploymentStatusResponse{}, errors.New("operation id cannot be empty")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		apiStudioV1.JoinPath(operationID).String(),
		http.NoBody,
	)
	if err != nil {
		return CheckDeploymentStatusResponse{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[CheckDeploymentStatusResponse](c, r)
}

// CheckDeploymentStatusResponse represents the result of checking the
// deployment status. This object has no documentation.
// There is no documentation for this object.
type CheckDeploymentStatusResponse struct {
	Name     string `json:"name"`
	Metadata struct {
		Type string `json:"@type"`
	} `json:"metadata"`
	Done     bool `json:"done"`
	Response struct {
		Type string `json:"@type"`
	} `json:"response"`
}

// DeploymentResponse represents the result of the deployment.
// This object has no documentation.
type DeploymentResponse struct {
	// Format:
	// workspaces/{workspace_id}/scenes/{scene_name}/operations/{operation_id}
	// or
	// workspaces/{workspace_id}/characters/{character_name}/operations/{operation_id}
	// or
	// workspaces/{workspace_id}/common-knowledge/{common_knowledge_id}/operations/{operation_id}
	Name     string `json:"name"`
	Metadata struct {
		Type string `json:"@type"`
	} `json:"metadata"`
	Done bool `json:"done"`
}
