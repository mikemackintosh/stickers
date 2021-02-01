package stickers

import (
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var serviceAccountFile []byte

// SetServiceAccountFile sets the service account file.
func SetServiceAccountFile(f string) error {
	var err error
	serviceAccountFile, err = ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	return nil
}

// NewClientWithSubject will impersonate a users' email and create a service
// to query the AdminSDK API.
func NewClientWithSubject(ctx context.Context, subject string, scopes []string) (option.ClientOption, error) {

	if serviceAccountFile == nil {
		return nil, ErrMissingServiceAccountFile{}
	}

	config, err := google.JWTConfigFromJSON(serviceAccountFile, scopes...)
	if err != nil {
		return nil, fmt.Errorf("JWTConfigFromJSON: %v", err)
	}
	config.Subject = subject

	ts := config.TokenSource(ctx)

	return option.WithTokenSource(ts), nil
}

// NewClient will use FindDefaultCredentials to generate a new client for checking
// default application credentials or using a configJSON.
func NewClient(ctx context.Context, scopes []string) (option.ClientOption, error) {

	if serviceAccountFile == nil {
		return nil, ErrMissingServiceAccountFile{}
	}

	config, err := google.JWTConfigFromJSON(serviceAccountFile, scopes...)
	if err != nil {
		return nil, fmt.Errorf("JWTConfigFromJSON: %v", err)
	}

	ts := config.TokenSource(ctx)

	return option.WithTokenSource(ts), nil
}

// NewDefaultCredentialsClient will use FindDefaultCredentials to generate a new client for checking
// default application credentials or using a configJSON.
func NewDefaultCredentialsClient(ctx context.Context, scopes ...string) (option.ClientOption, error) {
	credentials, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		return nil, err
	}

	return option.WithTokenSource(credentials.TokenSource), nil
}
