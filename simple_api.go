package inworld

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// All types in this file are generated based on this documentation:
// https://docs.inworld.ai/docs/tutorial-api/reference

// SimpleSendText rpc to send simple text request directly to single character.
func (c Client) SimpleSendText(ctx context.Context, req SimpleSendTextRequest) (Interaction, error) {
	if req.Character == "" {
		return Interaction{}, errors.New("character is required")
	}

	if req.Text == "" {
		return Interaction{}, errors.New("text is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiV1.JoinPath(req.Character+":simpleSendText").String(),
		newReader(req),
	)
	if err != nil {
		return Interaction{}, errors.Wrap(err, "creating request")
	}

	return sendSimpleAPIRequest[Interaction](c, r, req.SessionID)
}

// OpenSession rpc to load world for the interaction session.
func (c Client) OpenSession(ctx context.Context, req OpenSessionRequest) (Session, error) {
	if req.Name == "" {
		return Session{}, errors.New("name is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiV1.JoinPath(req.Name+":openSession").String(),
		newReader(req),
	)
	if err != nil {
		return Session{}, errors.Wrap(err, "creating request")
	}

	return sendSimpleAPIRequest[Session](c, r, "")
}

// SendText rpc to send text to the previously opened session.
func (c Client) SendText(ctx context.Context, req SendTextRequest) (Interaction, error) {
	if req.SessionID == "" {
		return Interaction{}, errors.New("session id is required")
	}

	if req.SessionCharacter == "" {
		return Interaction{}, errors.New("session character is required")
	}

	if req.Text == "" {
		return Interaction{}, errors.New("text is required")
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiV1.JoinPath(req.SessionCharacter+":sendText").String(),
		newReader(req),
	)
	if err != nil {
		return Interaction{}, errors.Wrap(err, "creating request")
	}

	return sendSimpleAPIRequest[Interaction](c, r, req.SessionID)
}

// SendTrigger rpc to send trigger event to the previously opened session.
func (c Client) SendTrigger(ctx context.Context, req SendTriggerRequest) (Interaction, error) {
	if req.SessionID == "" {
		return Interaction{}, errors.New("session id is required")
	}

	if req.SessionCharacter == "" {
		return Interaction{}, errors.New("session character is required")
	}

	if req.TriggerEvent.Trigger == "" {
		return Interaction{}, errors.New("trigger is required")
	}

	for _, p := range req.TriggerEvent.Parameters {
		if p.Name == "" {
			return Interaction{}, errors.New("parameter name is required")
		}
		if p.Value == "" {
			return Interaction{}, errors.New("parameter value is required")
		}
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiV1.JoinPath(req.SessionCharacter+":sendTrigger").String(),
		newReader(req),
	)
	if err != nil {
		return Interaction{}, errors.Wrap(err, "creating request")
	}

	return sendSimpleAPIRequest[Interaction](c, r, req.SessionID)
}

// SimpleSendTextRequest request message for
// [Sessions.SimpleSendText][ai.inworld.engine.v1.Sessions.SimpleSendText].
// https://docs.inworld.ai/docs/tutorial-api/reference/#simplesendtextrequest
type SimpleSendTextRequest struct {
	// Full resource name of the character to send text to. Format
	// workspaces/{workspace}/characters/{character}.
	Character string `json:"character"` // Required.
	// Text to send to the character.
	Text string `json:"text"` // Required.
	// Unique id of the session.
	SessionID string `json:"sessionId,omitempty"` // Optional.
	// Globally unique string, id of the end user of the system.
	EndUserID string `json:"endUserId,omitempty"` // Optional.
	// Full display name of the end user, this will be used by character in dialog.
	EndUserFullname string `json:"endUserFullname,omitempty"` // Optional.
}

// OpenSessionRequest request message for
// [Sessions.OpenSession][ai.inworld.engine.v1.Sessions.OpenSession].
// https://docs.inworld.ai/docs/tutorial-api/reference/#opensessionrequest
type OpenSessionRequest struct {
	// Scene or character resource name to be loaded by default. Format:
	// workspaces/{workspace}/scenes/{scene} OR
	// workspaces/{workspace}/characters/{character}.
	Name string `json:"name"` // Required.
	// Configuration of the experience consumer. End User information.
	User EndUserConfig `json:"user,omitempty"` // Optional.
}

// EndUserConfig represents the configuration of the end user of the system.
// https://docs.inworld.ai/docs/tutorial-api/reference/#enduserconfig
type EndUserConfig struct {
	// Globally unique string, id of the end user of the system.
	EndUserID string `json:"endUserId,omitempty"` // Optional.
	// Full name of the user to be used by characters through the dialog.
	GivenName string `json:"givenName,omitempty"` // Optional.
	// Gender of user.
	Gender string `json:"gender,omitempty"` // Optional.
	// Role of user.
	Role string `json:"role,omitempty"` // Optional.
	// Age of user.
	Age int64 `json:"age,string,omitempty"` // Optional.
}

// SendTextRequest request message for
// [Sessions.SendText][ai.inworld.engine.v1.Sessions.SendText].
// https://docs.inworld.ai/docs/tutorial-api/reference/#sendtextrequest
type SendTextRequest struct {
	// Unique id of the session.
	SessionID string `json:"-"` // Required.
	// Full resource name of the addressed session character. Format:
	// workspaces/{workspace}/sessions/{session}/sessionCharacters/{session_character}
	SessionCharacter string `json:"-"` // Required.
	// Text message to send to the character.
	Text string `json:"text"` // Required.
}

// SendTriggerRequest request message for
// [Sessions.SendTrigger][ai.inworld.engine.v1.Sessions.SendTrigger].
// https://docs.inworld.ai/docs/tutorial-api/reference/#sendtriggerrequest
type SendTriggerRequest struct {
	// Unique id of the session.
	SessionID string `json:"-"` // Required.
	// Full resource name of the addressed session character. Format:
	// workspaces/{workspace}/sessions/{session}/sessionCharacters/{session_character}
	SessionCharacter string `json:"-"` // Required.
	// Custom event to send.
	TriggerEvent TriggerEvent `json:"triggerEvent"` // Required.
	// Globally unique string, id of the end user of the system.
	EndUserID string `json:"endUserId,omitempty"` // Optional.
}

// TriggerEvent message describing trigger event structure.
// https://docs.inworld.ai/docs/tutorial-api/reference/#triggerevent
type TriggerEvent struct {
	// Name of the trigger event. Format: workspaces/{workspace}/triggers/{eventId}
	Trigger string `json:"trigger"` // Required.
	// Map of parameters to send along with the custom_event.
	Parameters []Parameter `json:"parameters,omitempty"` // Optional.
}

// Parameter supports string values for now, but can be expanded in future on
// as-needed basis.
// https://docs.inworld.ai/docs/tutorial-api/reference/#parameter
type Parameter struct {
	// Key name for parameter.
	Name string `json:"name"` // Required.
	// Key value for parameter.
	Value string `json:"value"` // Required.
}

// Session response message for LoadScene RPC.
// https://docs.inworld.ai/docs/tutorial-api/reference/#session
type Session struct {
	// Full resource name of the session. Format:
	// workspaces/{workspace}/sessions/{session}
	Name string `json:"name"`
	// List of created agents
	SessionCharacters []SessionCharacter `json:"sessionCharacters"`
	// Full resource name of the default loaded scene. Format:
	// workspaces/{workspace}/scenes/{scene} OR
	// workspaces/{workspace}/characters/{character}.
	LoadedScene string `json:"loadedScene"`
}

// SessionCharacter message describing the runtime instance of the character.
// https://docs.inworld.ai/docs/tutorial-api/reference/#sessionsessioncharacter
type SessionCharacter struct {
	// Full resource name of the session character. Format:
	// workspaces/{workspace}/sessions/{session}/sessionCharacters/{session_character}
	Name string `json:"name"`
	// Full resource name of the referenced character. Format:
	// workspaces/{workspace}/characters/{character}
	Character string `json:"character"`
	// Character display name specified by developer.
	DisplayName string `json:"displayName"`
	// Different assets for the character.
	CharacterAssets CharacterAsset `json:"characterAssets"`
}

// CharacterAsset represents structured list of different assets associated with
// the character.
// https://docs.inworld.ai/docs/tutorial-api/reference/#sessionsessioncharactercharacterassets
type CharacterAsset struct {
	// Link to uploaded and resized avatar for character.
	AvatarImage string `json:"avatarImg"`
	// Link to uploaded original avatar for character.
	AvatarImageOriginal string `json:"avatarImgOriginal"`
}

// Interaction basic response message for a Simple Text request to a character.
// https://docs.inworld.ai/docs/tutorial-api/reference/#interaction
type Interaction struct {
	// Full resource name of the interaction. Format:
	// workspaces/{workspace}/sessions/{session}/interactions/{interaction}
	Name string `json:"name"`
	// List of text responses. This list includes narrated actions in a form of
	// something happened in correct order.
	TextList []string `json:"textList"`
	// Describes emotion of the character on the moment of the end of interaction.
	Emotion Emotion `json:"emotion"`
	// SessionID that the interaction is in.
	SessionID string `json:"sessionId"`
	// Relationship updates.
	RelationshipUpdate RelationshipUpdate `json:"relationshipUpdate"`
	// Active trigger.
	ActiveTriggers []TriggerEvent `json:"activeTriggers"`

	// There is no documentation for these fields.
	CustomEvent struct {
		CustomEvent string `json:"customEvent"`
		Parameters  []any  `json:"parameters"`
	} `json:"customEvent"`
	Parameters map[string]any `json:"parameters"`
}

// Emotion describes emotion of the session character.
// https://docs.inworld.ai/docs/tutorial-api/reference/#interactionemotion
type Emotion struct {
	Behavior SpaffCode `json:"behavior"`
	Strength string    `json:"strength"`
}

// SpaffCode describes character behavior affected by emotions. For more details
// of character emotions, check out Personality and Emotions:
// https://docs.inworld.ai/docs/tutorial-basics/personality-emotion/#emotions
// https://docs.inworld.ai/docs/tutorial-api/reference/#interactionemotionspaffcode
type SpaffCode string

const (
	SpaffCodeUnspecified   SpaffCode = "SPAFF_CODE_UNSPECIFIED" // No emotion (mostly likely an error message).
	ScaffCodeNeutral       SpaffCode = "NEUTRAL"
	ScaffCodeDisgust       SpaffCode = "DISGUST"
	ScaffCodeContempt      SpaffCode = "CONTEMPT"
	ScaffCodeBelligerence  SpaffCode = "BELLIGERENCE"
	ScaffCodeDomineering   SpaffCode = "DOMINEERING"
	ScaffCodeCriticism     SpaffCode = "CRITICISM"
	ScaffCodeAnger         SpaffCode = "ANGER"
	ScaffCodeTension       SpaffCode = "TENSION"
	ScaffCodeTenseHumor    SpaffCode = "TENSE_HUMOR"
	ScaffCodeDefensiveness SpaffCode = "DEFENSIVENESS"
	ScaffCodeWhining       SpaffCode = "WHINING"
	ScaffCodeSadness       SpaffCode = "SADNESS"
	ScaffCodeStonewalling  SpaffCode = "STONEWALLING"
	ScaffCodeInterest      SpaffCode = "INTEREST"
	ScaffCodeValidation    SpaffCode = "VALIDATION"
	ScaffCodeAffection     SpaffCode = "AFFECTION"
	ScaffCodeHumor         SpaffCode = "HUMOR"
	ScaffCodeSurprise      SpaffCode = "SURPRISE"
	ScaffCodeJoy           SpaffCode = "JOY"
)

// Strength describes strength of the emotion.
// https://docs.inworld.ai/docs/tutorial-api/reference/#interactionemotionstrength
type Strength string

const (
	StrengthUnspecified Strength = "STRENGTH_UNSPECIFIED" // No strength (mostly likely an error message).
	StrengthWeak        Strength = "WEAK"
	StrengthStrong      Strength = "STRONG"
	StrengthNormal      Strength = "NORMAL"
)

// RelationshipUpdate shows changes in relationship based on latest interaction.
// https://docs.inworld.ai/docs/tutorial-api/reference/#relationshipupdate
type RelationshipUpdate struct {
	// How much the character trusts the user.
	Trust int `json:"trust"`
	// How much respect the character has for the user.
	Respect int `json:"respect"`
	// How familiar/comfortable the character is with the user.
	Familiar int `json:"familiar"`
	// How flirtatious the character is with the user.
	Flirtatious int `json:"flirtatious"`
	// How attracted the character is to the user.
	Attraction int `json:"attraction"`
}
