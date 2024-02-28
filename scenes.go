package inworld

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// All types in this file are generated based on this documentation:
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes

// CreateScene initiates the creation of a scene that requires subsequent
// deployment for activation.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes/#create-scene
func (c Client) CreateScene(ctx context.Context, workspaceID string, scene Scene) (Scene, error) {
	if workspaceID == "" {
		return Scene{}, errors.New("workspace id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiStudioV1.JoinPath("workspaces", workspaceID, "scenes").String(),
		newReader(scene),
	)
	if err != nil {
		return Scene{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[Scene](c, r)
}

// GetScene to get a specific scene within a workspace.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes/#get-scene
func (c Client) GetScene(
	ctx context.Context,
	sceneID string,
	// Specifies standard view or view with meta info in the response.
	view SceneItemView,
) (Scene, error) {
	if sceneID == "" {
		return Scene{}, errors.New("scene id is required")
	}

	url := apiStudioV1.JoinPath(sceneID)
	if view != "" {
		q := url.Query()
		q.Add("view", string(view))
		url.RawQuery = q.Encode()
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.String(),
		http.NoBody,
	)
	if err != nil {
		return Scene{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[Scene](c, r)
}

// DeployScene asynchronously deploys the scene. The deployment process is
// managed as a long-running operation (LRO). The progress and result of this
// operation should be monitored using the returned LRO object. Upon successful
// completion, all characters will reflect the updates during their
// interactions. In the event of a failure, the operation should be retried.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes/#deploy-scene
func (c Client) DeployScene(ctx context.Context, sceneID string) (DeploymentResponse, error) {
	if sceneID == "" {
		return DeploymentResponse{}, errors.New("scene id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiStudioV1.JoinPath(sceneID+":deploy").String(),
		http.NoBody,
	)
	if err != nil {
		return DeploymentResponse{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[DeploymentResponse](c, r)
}

// GetScenes returns a list of scenes that can be filtered by several criteria.
// When using pagination, ensure that all other parameters provided initially
// remain unchanged.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes/#list-scenes
func (c Client) GetScenes(
	ctx context.Context,
	req GetScenesRequest,
) (GetScenesResponse, error) {
	if req.WorkspaceID == "" {
		return GetScenesResponse{}, errors.New("workspace id is required")
	}

	url := apiStudioV1.JoinPath("workspaces", req.WorkspaceID, "scenes")
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
		return GetScenesResponse{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[GetScenesResponse](c, r)
}

// UpdateScene updates the specified character. Changes to the character are not
// reflected in conversation until the character is deployed.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes/#update-scene
func (c Client) UpdateScene(
	ctx context.Context,
	sceneID string,
	k Scene,
) (Scene, error) {
	if sceneID == "" {
		return Scene{}, errors.New("scene id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		apiStudioV1.JoinPath(sceneID).String(),
		newReader(k),
	)
	if err != nil {
		return Scene{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[Scene](c, r)
}

// DeleteScene to delete a specific scene within a workspace.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes/#delete-scene
func (c Client) DeleteScene(ctx context.Context, sceneID string) error {
	if sceneID == "" {
		return errors.New("scene id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		apiStudioV1.JoinPath(sceneID).String(),
		http.NoBody,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = sendStudioAPIRequest[struct{}](c, r)
	return err
}

// GetScenesRequest is a struct representing a request to list scene items.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes/#request-body-2
type GetScenesRequest struct {
	WorkspaceID string // Required.
	// Max number of items to retrieve per page. Default is 50.
	PageSize int32 // Optional.
	// A page token, received from a previous GetScenes call. Provide this
	// to retrieve the subsequent page. When paginating, all other parameters
	// provided to GetScenesRequest must stay the same.
	PageToken string // Optional.
	// Scenes filter. Filters follow the Google AIP-160 guidelines.
	//  - Single Filter: This filter accepts only the full resource name of the
	//  scene. For example, to filter a specific character, use
	//  scene.name=workspaces/{workspace_id}/scenes/{uuid}.
	//  - Multiple Filters: You can also combine multiple filters using the 'OR'
	//  condition. For instance, to filter by two different scenes, use
	//  scene.name=workspaces/{workspace_id}/scenes/{uuid1} OR
	//  scene.name=workspaces/{workspace_id}/scenes/{uuid2}.
	Filter string // Optional.
}

// GetScenesResponse is a struct representing the response from a get
// scenes request.
// There is no documentation for this object.
type GetScenesResponse struct {
	Scenes        []Scene `json:"scenes"`
	NextPageToken string  `json:"nextPageToken"`
}

// Scene represents a description of the Scene.
type Scene struct {
	// Immutable. This field can't be set or changed via API. Automatically
	// generated resource name by the system based on display_name property.
	Name string `json:"name,omitempty"` // Optional.
	// Scene's description. This field should give clear information about the
	// scene.
	Description string `json:"description,omitempty"` // Optional.
	// List of assigned triggers and scenes.
	SceneTriggers []SceneTrigger `json:"sceneTriggers"` // Optional.
	// User specified name.
	DisplayName string `json:"displayName,omitempty"` // Optional.
	// Immutable. This field can't be set or changed via API. Meta character
	// information.
	Meta *Meta `json:"meta"` // Optional.
	// List of assigned common_knowledge (contains resource references).
	CommonKnowledge []string `json:"commonKnowledge"` // Optional.
	// Current time period.
	TimePeriod string `json:"timePeriod"` // Optional.
	// List of references to scene characters.
	Characters []SceneCharacterReference `json:"characters"` // Optional.

	// There is no documentation for this field.
	InworldTags []any `json:"inworldTags"`
	// There is no documentation for this field.
	DefaultSceneAssets struct {
		SceneIMG         string `json:"sceneImg"`
		SceneIMGOriginal string `json:"sceneImgOriginal"`
	} `json:"defaultSceneAssets"`
}

// SceneCharacterReference holds scene character reference.
type SceneCharacterReference struct {
	// The resource name of the referenced character. Format:
	// workspaces/{workspace}/characters/{character}
	Character string `json:"character"` // Required.
	// Immutable. This field can't be set or changed via API. Character display
	// title.
	DisplayTitle string `json:"displayTitle,omitempty"` // Optional.
	// Immutable. This field can't be set or changed via API. Character display
	// image.
	ImageURI string `json:"imageUri,omitempty"` // Optional.
}

// SceneTrigger holds assigned triggers and scenes.
type SceneTrigger struct {
	// Reference to the trigger.
	Trigger string `json:"trigger,omitempty"` // Optional.
	// Scene-trigger description.
	Description string `json:"description,omitempty"` // Optional.
}

// SceneItemView specifies the provider of the information.
type SceneItemView string

const (
	// SceneItemViewUnspecified represents not specified, equivalent to DEFAULT.
	SceneItemViewUnspecified SceneItemView = "SCENE_ITEM_VIEW_UNSPECIFIED"
	// SceneItemViewDefault represents scene information by default.
	SceneItemViewDefault SceneItemView = "SCENE_ITEM_VIEW_DEFAULT"
	// SceneItemViewWithMeta represents scene representation with meta.
	SceneItemViewWithMeta SceneItemView = "SCENE_ITEM_VIEW_WITH_META"
)
