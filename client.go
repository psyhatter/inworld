// Package inworld provides access to API from inworld.ai.
package inworld

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"io"
	"net/http"
	"net/url"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// Error response format in case of errors. The only documentation that managed
// to find is here:
// https://docs.inworld.ai/docs/tutorial-integrations/node/api/#error-handling
type Error struct {
	Code    codes.Code    `json:"code"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
}

// Error implements error.
func (e *Error) Error() string { return e.Message }

func (e *Error) GRPCStatus() *status.Status {
	s := status.New(e.Code, e.Message)
	if len(e.Details) == 0 {
		return s
	}

	details, err := structpb.NewList(e.Details)
	if err != nil {
		return s
	}

	d := make([]proto.Message, len(details.GetValues()))
	for i, v := range details.GetValues() {
		d[i] = v
	}

	detailed, err := s.WithDetails(d...)
	if err != nil {
		return s
	}

	return detailed
}

// NewClient creates a new instance of the Client struct and initializes its
// fields with the provided values. It takes in two API keys (simpleAPIKey and
// studioAPIKey) as strings, and an http client.
func NewClient(simpleAPIKey, studioAPIKey string, client http.Client) Client {
	return Client{
		simpleAPIKey: simpleAPIKey,
		studioAPIKey: studioAPIKey,
		client:       client,
	}
}

type Client struct {
	simpleAPIKey string
	studioAPIKey string
	client       http.Client
}

var (
	apiV1       = api.JoinPath("v1")
	apiStudioV1 = api.JoinPath("studio/v1")

	api = func() *url.URL {
		u, err := url.Parse("https://api.inworld.ai")
		if err != nil {
			panic(err)
		}
		return u
	}()
)

func sendStudioAPIRequest[T any](c Client, r *http.Request) (T, error) {
	r.Header.Set("Authorization", "Basic "+c.studioAPIKey)
	r.Header.Set("Grpc-Metadata-X-Authorization-Bearer-Type", "studio_api")
	return sendRequest[T](c, r)
}

func sendSimpleAPIRequest[T any](c Client, r *http.Request, sessionID string) (T, error) {
	r.Header.Set("Authorization", "Basic "+c.simpleAPIKey)
	if sessionID != "" {
		r.Header.Set("Grpc-Metadata-Session-Id", sessionID)
	}
	return sendRequest[T](c, r)
}

func sendRequest[T any](c Client, r *http.Request) (response T, err error) {
	resp, err := c.client.Do(r)
	if err != nil {
		return response, errors.WithStack(err)
	}

	defer func() { err = combine(err, errors.WithStack(resp.Body.Close())) }()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, errors.Wrap(err, "reading http body")
	}

	rdr := json.NewDecoder(bytes.NewReader(b))

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var e Error
		if err = rdr.Decode(&e); err != nil || e.Code == codes.OK {
			return response, errors.Errorf(
				"request failed with status %d: %s",
				resp.StatusCode,
				limit(b, 200),
			)
		}
		return response, errors.WithStack(&e)
	}

	if err = rdr.Decode(&response); err != nil {
		return response, errors.Wrapf(err, "json unmarshaling to %T: %s", response, limit(b, 200))
	}

	return response, nil
}

func limit(v []byte, limit int) []byte {
	if len(v) > limit {
		return v[:limit]
	}
	return v
}

func newReader(v any) *jsonReader { return &jsonReader{v: v} }

type jsonReader struct {
	v   any
	buf *bytes.Buffer
}

func (r *jsonReader) Read(p []byte) (n int, err error) {
	if r.buf == nil {
		r.buf = bytes.NewBuffer(nil)
		err = json.NewEncoder(r.buf).Encode(r.v)
		if err != nil {
			return 0, errors.Wrap(err, "marshaling")
		}
	}

	return r.buf.Read(p)
}

func combine(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	return errors.WithStack(stderrors.Join(err1, err2))
}
