package inworld

import (
	"context"
	stderrors "errors"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// All types in this file are generated based on this documentation:
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters

// CreateCharacter initiates the creation of a character that requires
// subsequent deployment for activation. The character can't be used in
// conversation until it is deployed.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#create-character
func (c Client) CreateCharacter(ctx context.Context, workspaceID string, ch Character) (Character, error) {
	if workspaceID == "" {
		return Character{}, stderrors.New("workspace id is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiStudioV1.JoinPath("workspaces", workspaceID, "characters").String(),
		newReader(ch),
	)
	if err != nil {
		return Character{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[Character](c, r)
}

// GetCharacter returns a specific character within a workspace.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#get-character
func (c Client) GetCharacter(
	ctx context.Context,
	characterName string,
	// Specifies what information to include in the response.
	view CharacterItemView,
) (Character, error) {
	if characterName == "" {
		return Character{}, stderrors.New("character name is required")
	}

	url := apiStudioV1.JoinPath(characterName)
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
		return Character{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[Character](c, r)
}

// DeployCharacter asynchronously deploys the character. The deployment process
// is managed as a long-running operation (LRO). The progress and result of this
// operation should be monitored using the returned LRO object. Upon successful
// completion, all characters will reflect the updates during their
// interactions. In the event of a failure, the operation should be retried.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#deploy-character
func (c Client) DeployCharacter(
	ctx context.Context,
	characterName string,
) (DeploymentResponse, error) {
	if characterName == "" {
		return DeploymentResponse{}, stderrors.New("character name is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiStudioV1.JoinPath(characterName+":deploy").String(),
		http.NoBody,
	)
	if err != nil {
		return DeploymentResponse{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[DeploymentResponse](c, r)
}

// GetCharacters returns a list of characters that can be filtered by several
// criteria. When using pagination, ensure that all other parameters provided
// initially remain unchanged.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#list-characters
func (c Client) GetCharacters(ctx context.Context, req GetCharactersRequest) (GetCharactersResponse, error) {
	url := apiStudioV1.JoinPath("workspaces", req.WorkspaceID, "characters")
	q := url.Query()
	if req.View != "" {
		q.Add("view", string(req.View))
	}
	if req.PageSize > 0 {
		q.Add("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageToken != "" {
		q.Add("pageToken", req.PageToken)
	}
	if req.Filter != "" {
		q.Add("filter", req.Filter)
	}

	url.RawQuery = q.Encode()

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url.String(),
		http.NoBody,
	)
	if err != nil {
		return GetCharactersResponse{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[GetCharactersResponse](c, r)
}

// UpdateCharacter updates the specified character. Changes to the character are
// not reflected in conversation until the character is deployed.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#update-character
func (c Client) UpdateCharacter(ctx context.Context, characterName string, upd Character) (Character, error) {
	if characterName == "" {
		return Character{}, stderrors.New("character name cannot be empty")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		apiStudioV1.JoinPath(characterName).String(),
		newReader(upd),
	)
	if err != nil {
		return Character{}, errors.WithStack(err)
	}

	return sendStudioAPIRequest[Character](c, r)
}

// DeleteCharacter deletes a specific character within a workspace.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#delete-character
func (c Client) DeleteCharacter(ctx context.Context, characterName string) error {
	if characterName == "" {
		return stderrors.New("character name cannot be empty")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		apiStudioV1.JoinPath(characterName).String(),
		http.NoBody,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = sendStudioAPIRequest[struct{}](c, r)
	return err
}

// GetCharactersRequest represents a request for retrieving characters.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#request-body-2
type GetCharactersRequest struct {
	WorkspaceID string // Required.
	// Max number of items to retrieve per page. Default is 50.
	PageSize int32 // Optional.
	// A page token received from a previous GetCharactersResponse. Provide this to
	// retrieve the subsequent page. When paginating, all other parameters provided
	// to GetCharacters must remain the same.
	PageToken string // Optional.
	// Specifies whether standard or with scenes character list will be returned in
	// the response.
	View CharacterView // Optional.
	// Filters can be applied to GetCharactersRequest. Filters follow the Google
	// AIP-160 guidelines.
	//
	// 	- Single Filter: This filter accepts only the full resource name of the
	// 	character. For example, to filter a specific character, use
	// 	character.name=workspaces/{workspace_id}/character/{uuid}.
	// 	- Multiple
	// 	Filters: You can also combine multiple filters using the 'OR' condition. For
	// 	instance, to filter by two different characters, use
	// 	character.name=workspaces/{workspace_id}/character/{uuid1} OR
	// 	character.name=workspaces/{workspace_id}/character/{uuid2}.
	Filter string // Optional.
}

// GetCharactersResponse represents the response object for the GetCharacters
// API.
// There is no documentation for this object.
type GetCharactersResponse struct {
	Characters    []Character `json:"characters"`
	NextPageToken string      `json:"nextPageToken"`
}

// Character represents a character with various properties and configurations.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#character
type Character struct {
	// Name is an automatically generated resource name by the system based on
	// given_name property. This field can't be set or changed via API.
	Name string `json:"name,omitempty"` // Optional.
	// DefaultCharacterDescription provides structured default character
	// description.
	DefaultCharacterDescription CharacterDescription `json:"defaultCharacterDescription,omitempty"` // Optional.
	// DefaultCharacterAssets provides structured list of different assets
	// associated with the character.
	DefaultCharacterAssets CharacterAssets `json:"defaultCharacterAssets,omitempty"` // Optional.
	// InitialMood defines the character's initial mood.
	InitialMood CharacterInitialMood `json:"initialMood,omitempty"` // Optional.
	// Personality provides structured description of the character's personality.
	Personality CharacterPersonality `json:"personality,omitempty"` // Optional.
	// CommonKnowledge is the list of assigned common knowledge references.
	CommonKnowledge []string `json:"commonKnowledge"` // Optional.
	// SocialRank is the character's social rank - the insecure/confident slider’s
	// setting. Affects character conversation. For more details:
	// https://docs.inworld.ai/docs/tutorial-basics/personality-emotion/#mood-and-personality-sliders
	SocialRank float32 `json:"socialRank"` // Optional.
	// Scenes represent the list of linked scenes to the character.
	// his field is propagated only for list characters method when CharacterView::WITH_SCENE is set.
	// There is no documentation for the field.
	Scenes []any `json:"scenes"` // Optional.
	// PersonalKnowledge represents information that character knows about.
	PersonalKnowledge *PersonalKnowledge `json:"personalKnowledge"` // Optional.
	// Meta is an immutable field that contains meta character information. This
	// field can't be set or changed via API.
	Meta *Meta `json:"meta"` // Optional.
	// EmotionalFluidity represents the degree of emotional fluidity in the
	// character's conversation.
	EmotionalFluidity float32 `json:"emotionalFluidity,omitempty"` // Optional.
	// SharePortalInfo is an immutable field that contains character portal/arcade
	// info. This field can't be set or changed via API.
	// There is no documentation for the field.
	SharePortalInfo any `json:"sharePortalInfo"` // Optional.
	// YamlConfig used for defining goals and actions v2.
	YamlConfig string `json:"yamlConfig"` // Optional.
	// SafetyConfig represents a list of safety configs.
	SafetyConfig SafetyConfigEntry `json:"safetyConfig,omitempty"` // Optional.
	// Relationship describes the character's relationships.
	Relationship Relationship `json:"relationship,omitempty"` // Optional.
	// CognitiveControl describes the degree of cognitive control.
	CognitiveControl CognitiveControl `json:"cognitiveControl,omitempty"` // Optional.

	// There is no documentation for this field.
	Language string `json:"language,omitempty"`
	// There is no documentation for this field.
	BehavioralContexts []interface{} `json:"behavioralContexts"`
	// There is no documentation for this field. Apparently this is an enum, one of
	// the values is FOURTH_WALL_CONTAINMENT_DISABLED.
	FourthWall string `json:"fourthWall,omitempty"`
	// There is no documentation for this field.
	InworldTags []any `json:"inworldTags"`
	// There is no documentation for this field.
	UserTags []any `json:"userTags"`
	// There is no documentation for this field.
	LongTermCoherence struct {
		Enabled bool `json:"enabled"`
	} `json:"longTermCoherence,omitempty"`
}

// CharacterAssets holds various assets associated with the character.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#characterassets
type CharacterAssets struct {
	// Voice is a reference to the voice resource used by that character.
	// There is no documentation for the field.
	Voice Voice `json:"voice,omitempty"` // Optional.
	// AvatarImg is a link to the uploaded and resized avatar image for the character provided by the user.
	AvatarImg string `json:"avatarImg"` // Optional.
	// AvatarImgOriginal is a link to the uploaded original avatar image for the character provided by the user.
	AvatarImgOriginal string `json:"avatarImgOriginal"` // Optional.
	// AvatarType defines the avatar type of the character.
	AvatarType AvatarType `json:"avatarType,omitempty"` // Optional.
	// AvatarDisplayImageSource specifies the source from which the 2D avatar image of the character is derived.
	AvatarDisplayImageSource AvatarDisplayImageSource `json:"avatarDisplayImageSource,omitempty"` // Optional.
	// RPMAvatar represents the RPM avatar of the character.
	RPMAvatar RPMAvatar `json:"rpmAvatar,omitempty"` // Optional.
	// InnequinAvatar represents the Innequin avatar of the character.
	InnequinAvatar InnequinAvatar `json:"innequinAvatar,omitempty"` // Optional.
}

// CharacterDescription describes a character with various properties.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#characterdescription
type CharacterDescription struct {
	// Name of the character. Defaults to 'Character' if empty.
	GivenName string `json:"givenName,omitempty"` // Required.
	// Free form description to design character's behavior and purpose. Supposed to
	// be composite prompt including NPC info, history and initial scene's state.
	Description string `json:"description"` // Optional.
	// Pronoun of the character. Defaults to PRONOUN_UNSPECIFIED.
	Pronoun Pronoun `json:"pronoun,omitempty"` // Optional.
	// List of commonly used alternative names of this character.
	Nicknames []string `json:"nicknames"` // Optional.
	// Motivation of the character.
	Motivation string `json:"motivation" json:"motivation,omitempty"` // Optional.
	// URI to wikipedia for well-known character for additional data extraction. For
	// more details: https://docs.inworld.ai/docs/tutorial-basics/identity/
	WikipediaURI string `json:"wikipediaUri" json:"wikipediaURI,omitempty"` // Optional.
	// Example of character dialog. For more details:
	// https://docs.inworld.ai/docs/tutorial-basics/dialog-style/#example-dialogue
	ExampleDialog string `json:"exampleDialog"` // Optional.
	// Enumerated example dialog style. For more details of dialog style:
	// https://docs.inworld.ai/docs/tutorial-basics/dialog-style/
	ExampleDialogStyle ExampleDialogStyle `json:"exampleDialogStyle,omitempty"` // Optional.
	// List of adjectives describing character: "eloquent", "charming", "aggressive".
	PersonalityAdjectives []string `json:"personalityAdjectives"` // Optional.
	// Specify the stage of life for character.
	LifeStage LifeStage `json:"lifeStage,omitempty"` // Optional.
	// Repeated field of strings representing the character's hobbies or interests.
	HobbyOrInterests []string `json:"hobbyOrInterests"` // Optional.
	// Role of the character.
	CharacterRole string `json:"characterRole"` // Optional.
	// Indicates if narrated actions and scenarios are enabled on this character. For more details
	NarrativeActionsEnabled bool `json:"narrativeActionsEnabled"` // Optional.
	// CustomDialogStyle indicates designs are using a new layout for Character pages.
	CustomDialogStyles []CustomDialogStyle `json:"customDialogStyles"` // Optional.
	// Flaws & Challenges.
	Flaws string `json:"flaws"` // Optional.
	// Specify dialog style response length.
	DialogResponseLength DialogResponseLength `json:"dialogResponseLength,omitempty"` // Optional.

	// There is no documentation for this field.
	ExternalDescription string `json:"externalDescription"`
}

// CharacterInitialMood determines the initial mood of a character.
// The mood values range from -100 to 100 for different emotions
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#characterinitialmood
type CharacterInitialMood struct {
	// Determines Sadness(-100)/Joy(100) mood.
	Joy int32 `json:"joy"` // Optional.
	// Determines Anger(-100)/Fear(100) mood.
	Fear int32 `json:"fear"` // Optional.
	// Determines Disgust(-100)/Trust(100) mood.
	Trust int32 `json:"trust"` // Optional.
	// Determines Anticipation(-100)/Surprise(100) mood.
	Surprise int32 `json:"surprise"` // Optional.
}

// CharacterPersonality describes the personality of a character.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#characterpersonality
type CharacterPersonality struct {
	// Positive determines Negative(-100)/Positive(100).
	Positive int32 `json:"positive"` // Optional.
	// Peaceful determines Aggressive(-100)/Peaceful(100).
	Peaceful int32 `json:"peaceful"` // Optional.
	// Open determines Cautious(-100)/Open(100).
	Open int32 `json:"open"` // Optional.
	// Extravert determines Introvert(-100)/Extravert(100).
	Extravert int32 `json:"extravert"` // Optional.
}

// CustomDialogStyle is a type that represents a custom dialog style.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#customdialogstyle
type CustomDialogStyle struct {
	// Universally unique identifier automatically generated by system.
	UUID string `json:"uuid"` // Optional.
	// The name of the custom dialog style.
	DisplayName string `json:"displayName"` // Optional.
	// Is this custom style enabled or not.
	IsActive bool `json:"isActive"` // Optional.
	// The list of available adjectives. Adjectives field should be multi-select (up
	// to 3): "anxious", "contemplative", "curt", "expressive", "hilarious", "inquisitive",
	// "intense", "mischievous", "snarky", "sympathetic", "talkative unflinching" etc.
	Adjectives []string `json:"adjectives"` // Optional.
	// The name of the new dialog style: "business abbreviations", "cowboy lingo",
	// "hackneyed cliches", "Internet meme slang", "medical terminology", "movie quotes",
	// "obscure idioms", "Orwellian newspeak", "sports metaphors", "technical jargon".
	Colloquialism string `json:"colloquialism"` // Optional.
}

// Fact describes a fact.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#fact
type Fact struct {
	// Text data of the fact. Max length is 255 symbols.
	Text string `json:"text,omitempty"` // Optional.
}

// InnequinAvatar describes the avatar.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#innequinavatar
type InnequinAvatar struct {
	// Reference to avatar material asset.
	AvatarBackground string `json:"avatarBackground"` // Optional.
	// Reference to avatar background asset.
	AvatarMaterial string `json:"avatarMaterial"` // Optional.
	// URI to store innequin image.
	AvatarImgURI string `json:"avatarImgUri"` // Optional.

	// There is no documentation for this field.
	AvatarImgOriginalURI string `json:"avatarImgOriginalUri"`
}

// InteractionCountStat describes the interaction count statistics.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#interactioncountstat
type InteractionCountStat struct {
	// Total log counters.
	TotalCount int32 `json:"totalCount,omitempty"` // Optional.
	// Total unread log counters.
	TotalUnreadCount int32 `json:"totalUnreadCount,omitempty"` // Optional.
	// Total read log counters.
	TotalReadCount int32 `json:"totalReadCount,omitempty"` // Optional.
}

// Meta describes the statistics of the character.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/scenes/#meta
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#meta
type Meta struct {
	// Immutable. This field can't be set or changed via API. Indicates the amount
	// of common knowledge created in character or scene.
	TotalCommonKnowledge int32 `json:"totalCommonKnowledge"` // Optional.
	// Immutable. This field can't be set or changed via API. Goals version used by
	// character.
	GoalsVersion int32 `json:"goalsVersion"` // Optional.
	// Immutable. This field can't be set or changed via API.
	// Indicates the number of triggers created in scene.
	TotalTriggers int32 `json:"totalTriggers"` // Optional.
	// Immutable. This field can't be set or changed via API.
	// Indicates the number of characters created in scene.
	TotalCharacters int32 `json:"totalCharacters"` // Optional.
}

// PersonalKnowledge represents personal knowledge of a character.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#personalknowledge
type PersonalKnowledge struct {
	// Unique (throughout the entire service) universally unique identifier
	// automatically generated by system.
	UUID string `json:"uuid,omitempty"` // Optional.
	// List of facts. Max size 10000.
	Facts []Fact `json:"facts,omitempty"` // Optional.
}

// Relationship describes the character's relationship configuration.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#relationship
type Relationship struct {
	// Enables character's relationships.
	Enabled bool `json:"enabled"` // Optional.
	// Determines dynamics of relationship changes (-1, 1), where -1 is slow changes
	// and 1 is fast changes.
	Fluidity float64 `json:"fluidity,omitempty"` // Optional.
}

// RPMAvatar describes the ReadyPlayerMe avatar's details.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#rpmavatar
type RPMAvatar struct {
	// Link to the readyplayerme avatar model. To be used in the companion app.
	// Updated through the separate method.
	RPMModelURI string `json:"rpmModelUri"` // Optional.
	// Link to the rendered readyplayerme avatar image. Used to display avatar image
	// in studio and companion app. Updated through the separate method.
	RPMImageURI string `json:"rpmImageUri"` // Optional.
	// Link to the rendered readyplayerme avatar portrait image. in studio and
	// companion app. Updated through the separate method.
	RPMImageURIPortrait string `json:"rpmImageUriPortrait"` // Optional.
	// Link to rendered readyplayerme avatar posture image. To be used in the in
	// studio and companion app. Updated through the separate method.
	RPMImageURIPosture string `json:"rpmImageUriPosture"` // Optional.
}

// SafetyConfigEntry represents a safety configuration entry.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#safetyconfigentry
//
// Examples of map keys:
//
//   - "TOPIC_ALCOHOL"
//   - "TOPIC_POLITICS"
//   - "TOPIC_VIOLENCE"
//   - "TOPIC_RELIGION"
//   - "TOPIC_PROFANITY"
//   - "TOPIC_ADULT_TOPICS"
//   - "TOPIC_SUBSTANCE_USE"
type SafetyConfigEntry map[string]SafetyLevel

// AvatarDisplayImageSource determines which 2D image will be displayed for
// current character.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#avatardisplayimagesource
type AvatarDisplayImageSource string

const (
	// AvatarDisplayImageSourceUnspecified represents unspecified value.
	AvatarDisplayImageSourceUnspecified AvatarDisplayImageSource = "AVATAR_DISPLAY_IMAGE_SOURCE_UNSPECIFIED"
	// AvatarDisplayImageSourceRPMImage use RPM image, that is created from RPM
	// avatar model.
	AvatarDisplayImageSourceRPMImage AvatarDisplayImageSource = "RPM_IMAGE"
	// AvatarDisplayImageSourceInnequinImage use Innequin image that is created from
	// Innequin avatar model.
	AvatarDisplayImageSourceInnequinImage AvatarDisplayImageSource = "INNEQUIN_IMAGE"
	// AvatarDisplayImageSourceUserProvidedImage use user provided image.
	AvatarDisplayImageSourceUserProvidedImage AvatarDisplayImageSource = "USER_PROVIDED_IMAGE"
)

// AvatarType determines which 3D avatar system is used in UI.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#avatartype
type AvatarType string

const (
	// AvatarTypeUnspecified represents unspecified value.
	AvatarTypeUnspecified AvatarType = "AVATAR_TYPE_UNSPECIFIED"
	// AvatarTypeRPM represents RPM avatar system.
	AvatarTypeRPM AvatarType = "AVATAR_TYPE_RPM"
	// AvatarTypeInWorld represents Inworld/Innequin avatar system.
	AvatarTypeInWorld AvatarType = "AVATAR_TYPE_INWORLD"
)

// CharacterItemView defines what dependency info should be provided in
// response. For now only Meta(s) part is optional.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#characteritemview
type CharacterItemView string

const (
	// CharacterItemViewUnspecified is not specified, equivalent to DEFAULT.
	CharacterItemViewUnspecified CharacterItemView = "CHARACTER_ITEM_VIEW_UNSPECIFIED"
	// CharacterItemViewDefault represents Character information default.
	CharacterItemViewDefault CharacterItemView = "CHARACTER_ITEM_VIEW_DEFAULT"
	// CharacterItemViewWithMeta represents Character representation with meta.
	CharacterItemViewWithMeta CharacterItemView = "CHARACTER_ITEM_VIEW_WITH_META"
)

// CharacterSharePurpose defines the purpose of share the character(where it is
// expected to be used):
//   - in oculus,
//   - in portal simulator (standalone).
//
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#charactersharepurpose
type CharacterSharePurpose string

const (
	// CharacterSharePurposeUnspecified is not specified, equivalent to DEFAULT.
	CharacterSharePurposeUnspecified CharacterSharePurpose = "CHARACTER_SHARE_PURPOSE_UNSPECIFIED"
	// CharacterSharePurposeOculus represents character shared for the oculus
	// device.
	CharacterSharePurposeOculus CharacterSharePurpose = "CHARACTER_SHARE_PURPOSE_OCULUS"
	// CharacterSharePurposePortal represents character shared for the portal
	// simulator.
	CharacterSharePurposePortal CharacterSharePurpose = "CHARACTER_SHARE_PURPOSE_PORTAL"
)

// CharacterView defines what dependency info should be provided in response.
// For now only Scene(s) part is optional.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#characterview
type CharacterView string

const (
	// CharacterViewUnspecified represents not specified, equivalent to DEFAULT.
	CharacterViewUnspecified CharacterView = "CHARACTER_VIEW_UNSPECIFIED"
	// CharacterViewDefault represents character information without scene
	// information.
	CharacterViewDefault CharacterView = "CHARACTER_VIEW_DEFAULT"
	// CharacterViewWithScenes represents character representation with scene.
	CharacterViewWithScenes CharacterView = "CHARACTER_VIEW_WITH_SCENES"
)

// CognitiveControl specifies degree of cognitive control.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#cognitivecontrol
type CognitiveControl string

const (
	// CognitiveControlUnspecified represents missing value.
	CognitiveControlUnspecified CognitiveControl = "COGNITIVE_CONTROL_UNSPECIFIED"
	// CognitiveControlNone represents none control.
	CognitiveControlNone CognitiveControl = "COGNITIVE_CONTROL_NONE"
	// CognitiveControlMild represents mild control.
	CognitiveControlMild CognitiveControl = "COGNITIVE_CONTROL_MILD"
	// CognitiveControlStrict represents strict control.
	CognitiveControlStrict CognitiveControl = "COGNITIVE_CONTROL_STRICT"
)

// DialogResponseLength represents dialog style response length.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#dialogresponselength
type DialogResponseLength string

const (
	// DialogResponseLengthUnspecified represents no response length
	// specified/unknown.
	DialogResponseLengthUnspecified DialogResponseLength = "DIALOG_RESPONSE_LENGTH_UNSPECIFIED"
	// DialogResponseLengthVeryShort represents very short response (-2).
	DialogResponseLengthVeryShort DialogResponseLength = "DIALOG_RESPONSE_LENGTH_VERY_SHORT"
	// DialogResponseLengthShort represents short response (-1).
	DialogResponseLengthShort DialogResponseLength = "DIALOG_RESPONSE_LENGTH_SHORT"
	// DialogResponseLengthMedium represents medium/neutral response (0).
	DialogResponseLengthMedium DialogResponseLength = "DIALOG_RESPONSE_LENGTH_MEDIUM"
	// DialogResponseLengthLong represents long response (1).
	DialogResponseLengthLong DialogResponseLength = "DIALOG_RESPONSE_LENGTH_LONG"
	// DialogResponseLengthVeryLong represents very long response (2).
	DialogResponseLengthVeryLong DialogResponseLength = "DIALOG_RESPONSE_LENGTH_VERY_LONG"
)

// ExampleDialogStyle represents a list of styles for example dialogue.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#exampledialogstyle
type ExampleDialogStyle string

const (
	// ExampleDialogStyleUnspecified represents no style specified / unknown.
	ExampleDialogStyleUnspecified ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_UNSPECIFIED"
	// ExampleDialogStyleDefault means the default character dialog is engaging and
	// realistic.
	ExampleDialogStyleDefault ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_DEFAULT"
	// ExampleDialogStyleBubbly means the dialog style is thoughtful, energetic, and
	// friendly.
	ExampleDialogStyleBubbly ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_BUBBLY"
	// ExampleDialogStyleFormal means the dialog style is formal and matter-of-fact.
	ExampleDialogStyleFormal ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_FORMAL"
	// ExampleDialogStyleBlunt means the dialog style is short, direct, and to the
	// point.
	ExampleDialogStyleBlunt ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_BLUNT"
	// ExampleDialogStyleInquisitive means the dialog style is inquisitive, which
	// will usually include a question.
	ExampleDialogStyleInquisitive ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_INQUISITIVE"
	// ExampleDialogStyleCommanding means the dialog style is intense and
	// determined.
	ExampleDialogStyleCommanding ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_COMMANDING"
	// ExampleDialogStyleEmpathetic means the dialog style is gentle and
	// compassionate.
	ExampleDialogStyleEmpathetic ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_EMPATHETIC"
	// ExampleDialogStyleEntertaining means the dialog style is entertaining,
	// animated, and comical.
	ExampleDialogStyleEntertaining ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_ENTERTAINING"
	// ExampleDialogStyleHypochondriac means the dialog style is neurotic and
	// fearful.
	ExampleDialogStyleHypochondriac ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_HYPOCHONDRIAC"
	// ExampleDialogStyleLaidback means the dialog style is relaxed, informal, and
	// nonchalant.
	ExampleDialogStyleLaidback ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_LAIDBACK"
	// ExampleDialogStyleLongWinded means the dialog style is verbose and
	// meandering.
	ExampleDialogStyleLongWinded ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_LONG_WINDED"
	// ExampleDialogStyleMoral means the dialog style is principled and
	// incorruptible.
	ExampleDialogStyleMoral ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_MORAL"
	// ExampleDialogStyleMysterious means the dialog style is elusive and cryptic.
	ExampleDialogStyleMysterious ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_MYSTERIOUS"
	// ExampleDialogStyleRaconteur means the dialog style is entertaining and
	// descriptive.
	ExampleDialogStyleRaconteur ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_RACONTEUR"
	// ExampleDialogStyleSarcastic means the dialog style is snide and contemptuous.
	ExampleDialogStyleSarcastic ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_SARCASTIC"
	// ExampleDialogStyleTenacious means the dialog style is determined and
	// persistent.
	ExampleDialogStyleTenacious ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_TENACIOUS"
	// ExampleDialogStyleVillainous means the dialog style is sinister and cruel.
	ExampleDialogStyleVillainous ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_VILLAINOUS"
	// ExampleDialogStyleCustom indicates a custom style selected by user and
	// expected to be specified.
	ExampleDialogStyleCustom ExampleDialogStyle = "EXAMPLE_DIALOG_STYLE_CUSTOM"
)

// LifeStage represents list of life stages for character.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#lifestage
type LifeStage string

const (
	// LifeStageUnspecified represents no stage of life specified / unknown.
	LifeStageUnspecified LifeStage = "LIFE_STAGE_UNSPECIFIED"
	// LifeStageAdolescence represents Adolescence.
	LifeStageAdolescence LifeStage = "LIFE_STAGE_ADOLESCENCE"
	// LifeStageYoungAdulthood represents Young adulthood.
	LifeStageYoungAdulthood LifeStage = "LIFE_STAGE_YOUNG_ADULTHOOD"
	// LifeStageMiddleAdulthood represents Middle adulthood.
	LifeStageMiddleAdulthood LifeStage = "LIFE_STAGE_MIDDLE_ADULTHOOD"
	// LifeStageLateAdulthood represents Late adulthood.
	LifeStageLateAdulthood LifeStage = "LIFE_STAGE_LATE_ADULTHOOD"
	// LifeStageChildhood represents Childhood.
	LifeStageChildhood LifeStage = "LIFE_STAGE_CHILDHOOD"
)

// Pronoun describes pronounce for the character.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#pronoun
type Pronoun string

const (
	// PronounUnspecified represents no pronoun specified / unknown.
	PronounUnspecified Pronoun = "PRONOUN_UNSPECIFIED"
	// PronounFemale represents She/Her/Hers.
	PronounFemale Pronoun = "PRONOUN_FEMALE"
	// PronounMale represents He/Him/His.
	PronounMale Pronoun = "PRONOUN_MALE"
	// PronounOther represents They/Them/Theirs.
	PronounOther Pronoun = "PRONOUN_OTHER"
)

// SafetyLevel defines the topic-based safety aggressiveness level config.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/characters/#safetylevel
type SafetyLevel string

const (
	// SafetyLevelUnspecified represents missing value.
	SafetyLevelUnspecified SafetyLevel = "SAFETY_LEVEL_UNSPECIFIED"
	// SafetyLevelNoControl represents no safety filter is applied to the given topic.
	SafetyLevelNoControl SafetyLevel = "SAFETY_LEVEL_NO_CONTROL"
	// SafetyLevelMildControl represents allows 'level 0' of given topic.
	SafetyLevelMildControl SafetyLevel = "SAFETY_LEVEL_MILD_CONTROL"
	// SafetyLevelStrictControl represents all filters applied to the given topic.
	SafetyLevelStrictControl SafetyLevel = "SAFETY_LEVEL_STRICT_CONTROL"
)

// Voice represents a description of the character's voice.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/voice/#voice-1
type Voice struct {
	// Voice base name.
	BaseName string `json:"baseName,omitempty"` // Optional.
	// TTS type.
	TTSType TTSType `json:"ttsType,omitempty"` // Optional.
	// Voice pitch. Default = 1, range = [-10, 10].
	Pitch float64 `json:"pitch"` // Optional.
	// Speaking rate. Default = 0, range = [0, 5].
	SpeakingRate float64 `json:"speakingRate,omitempty"` // Optional.
	// Robotic filter level. Default = 0, range = [0, 5].
	RoboticVoiceFilterLevel float64 `json:"roboticVoiceFilterLevel"` // Optional.
	// TTS platform metadata - tts_metadata type must be compatible with the
	// tts_type specified (enforced by service)
	TtsMetadata *ElevenLabsMetadata `json:"ttsMetadata,omitempty"` // Optional.

	// Apparently this is an enum, one of the values: VOICE_GENDER_MALE
	// There is no documentation for this field.
	Gender string `json:"gender,omitempty"` // Optional.
}

// StudioBaseVoice holds studio voice settings.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/voice/#studiobasevoice
type StudioBaseVoice struct {
	// Voice name.
	Name string `json:"name,omitempty"` // Optional.
	// TTS type.
	TTSType TTSType `json:"ttsType,"` // Optional.
	// TTS platform metadata - tts_metadata type must be compatible with the
	// tts_type specified (enforced by service)
	TtsMetadata ElevenLabsMetadata `json:"ttsMetadata,"` // Optional.

	// Voice gender.
	// There is no documentation for this field.
	Gender any `json:"gender,omitempty"` // Optional.
	// Voice age.
	// There is no documentation for this field.
	Age any `json:"age,omitempty"` // Optional.
}

// ElevenLabsMetadata holds eleven labs metadata.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/voice/#elevenlabsmetadata
type ElevenLabsMetadata struct {
	// Voice ID.
	VoiceID string `json:"voiceId,omitempty"` // Optional.
}

// TTSType specifies the provider of the voice.
// https://docs.inworld.ai/docs/tutorial-basics/studio-api/reference/voice/#ttstype
type TTSType string

const (
	// TTSTypeGoogle represents Google Voice.
	TTSTypeGoogle TTSType = "TTS_TYPE_GOOGLE"
	// TTSTypeInworld represents Inworld Voice.
	TTSTypeInworld TTSType = "TTS_TYPE_INWORLD"
	// TTSTypeElevenLabs represents Eleven Labs Voice.
	TTSTypeElevenLabs TTSType = "TTS_TYPE_ELEVEN_LABS"
	// TTSTypeInworldV2 represents Inworld_V2 Voice.
	TTSTypeInworldV2 TTSType = "TTS_TYPE_INWORLD_V2"
)
