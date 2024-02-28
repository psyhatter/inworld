package inworld

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// All types in this file are generated based on this documentation:
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge

// CreateCommonKnowledge initiates the creation of common knowledge that
// requires subsequent deployment for activation.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge/#create-common-knowledge
func (c Client) CreateCommonKnowledge(
	ctx context.Context,
	workspaceID string,
	k CommonKnowledge,
) (CommonKnowledge, error) {
	if workspaceID == "" {
		return CommonKnowledge{}, errors.New("workspace id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiStudioV1.JoinPath("workspaces", workspaceID, "common-knowledge").String(),
		newReader(k),
	)
	if err != nil {
		return CommonKnowledge{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[CommonKnowledge](c, r)
}

// GetCommonKnowledge to get a specific common knowledge entry within a
// workspace.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge/#get-common-knowledge
func (c Client) GetCommonKnowledge(ctx context.Context, commonKnowledgeID string) (CommonKnowledge, error) {
	if commonKnowledgeID == "" {
		return CommonKnowledge{}, errors.New("common knowledge id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		apiStudioV1.JoinPath(commonKnowledgeID).String(),
		http.NoBody,
	)
	if err != nil {
		return CommonKnowledge{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[CommonKnowledge](c, r)
}

// DeployCommonKnowledge asynchronously deploys common knowledge. The deployment
// process is managed as a long-running operation (LRO). The progress and result
// of this operation should be monitored using the returned LRO object. Upon
// successful completion, all characters and scenes that incorporate this common
// knowledge will reflect the updates during their interactions. In the event of
// a failure, the operation should be retried.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge/#deploy-common-knowledge
func (c Client) DeployCommonKnowledge(ctx context.Context, commonKnowledgeID string) (DeploymentResponse, error) {
	if commonKnowledgeID == "" {
		return DeploymentResponse{}, errors.New("common knowledge id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiStudioV1.JoinPath(commonKnowledgeID+":deploy").String(),
		http.NoBody,
	)
	if err != nil {
		return DeploymentResponse{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[DeploymentResponse](c, r)
}

// ListCommonKnowledge returns a list of common knowledge that can be filtered
// by several criteria. When using pagination, ensure that all other parameters
// provided initially remain unchanged.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge/#list-common-knowledge
func (c Client) ListCommonKnowledge(
	ctx context.Context,
	req ListCommonKnowledgeRequest,
) (ListCommonKnowledgeResponse, error) {
	if req.WorkspaceID == "" {
		return ListCommonKnowledgeResponse{}, errors.New("workspace id is required")
	}

	url := apiStudioV1.JoinPath("workspaces", req.WorkspaceID, "common-knowledge")
	q := url.Query()

	if req.Filter != "" {
		q.Add("filter", req.Filter)
	}
	if req.PageSize > 0 {
		q.Add("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageToken != "" {
		q.Add("pageToken", req.PageToken)
	}

	url.RawQuery = q.Encode()

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.String(),
		http.NoBody,
	)
	if err != nil {
		return ListCommonKnowledgeResponse{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[ListCommonKnowledgeResponse](c, r)
}

// UpdateCommonKnowledge updates the specified common knowledge. Changes to
// common knowledge are not reflected in conversation until common knowledge is
// deployed.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge/#update-common-knowledge
func (c Client) UpdateCommonKnowledge(
	ctx context.Context,
	commonKnowledgeID string,
	k CommonKnowledge,
) (CommonKnowledge, error) {
	if commonKnowledgeID == "" {
		return CommonKnowledge{}, errors.New("common knowledge id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		apiStudioV1.JoinPath(commonKnowledgeID).String(),
		newReader(k),
	)
	if err != nil {
		return CommonKnowledge{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[CommonKnowledge](c, r)
}

// DeleteCommonKnowledge deletes a specific common knowledge entry within a
// workspace.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge/#delete-common-knowledge
func (c Client) DeleteCommonKnowledge(ctx context.Context, commonKnowledgeID string) error {
	if commonKnowledgeID == "" {
		return errors.New("common knowledge id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		apiStudioV1.JoinPath(commonKnowledgeID).String(),
		http.NoBody,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = sendStudioAPIRequest[struct{}](c, r)
	return err
}

// ListCommonKnowledgeRequest is a struct representing a request to list common knowledge items.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge/#request-body-1
type ListCommonKnowledgeRequest struct {
	WorkspaceID string // Required.
	// Max number of items to retrieve per page. Default is 50.
	PageSize int32 // Optional.
	// A page token, received from a previous ListCommonKnowledge call. Provide this
	// to retrieve the subsequent page. When paginating, all other parameters
	// provided to ListCommonKnowledgeRequest must stay the same.
	PageToken string // Optional.
	// CommonKnowledge filter. Filters follow the Google AIP-160 guidelines.
	//
	//  - common_knowledge.name: This filter accepts only the full resource name of
	//  common_knowledge. For example, to filter a specific common knowledge item,
	//  use common_knowledge.name=workspaces/{workspace_id}/common-knowledge/{uuid}.
	//  - Multiple Filters: You can also combine multiple filters using the 'OR'
	//  condition. For instance, to filter by two different common knowledge items,
	//  use common_knowledge.name=workspaces/{workspace}/common-knowledge/{uuid1} OR
	//  common_knowledge.name=workspaces/{workspace}/common-knowledge/{uuid2}.
	Filter string // Optional.
}

// ListCommonKnowledgeResponse is a struct representing the response from a list
// common knowledge request.
// There is no documentation for this object.
type ListCommonKnowledgeResponse struct {
	CommonKnowledge []CommonKnowledge `json:"commonKnowledge"`
	NextPageToken   string            `json:"nextPageToken"`
}

// CommonKnowledge represents a piece of knowledge in the system.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/common-knowledge/#commonknowledge
type CommonKnowledge struct {
	// Name cannot be set or changed, this field is output only.
	// Format workspaces/{workspaceID}/common-knowledge/{commonKnowledgeID}
	Name string `json:"name,omitempty"`
	// Display name of the common knowledge.
	DisplayName string `json:"displayName,omitempty"` // Optional.
	// Free form description outlining the character's behavior and purpose. See
	// https://docs.inworld.ai/docs/tutorial-basics/common-knowledge
	Description string `json:"description,omitempty"` // Optional.
	// A list of memory records in a flexible free form format. Maximum length of
	// individual record is 255 characters. Maximum amount of records per knowledge
	// is 10000.
	MemoryRecords []string `json:"memoryRecords"` // Optional.

	// There is no documentation for this field.
	InworldTags []any `json:"inworldTags"`
}
