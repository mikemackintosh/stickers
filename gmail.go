package stickers

import (
	"context"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	ScopeGmailLabels   = gmail.GmailLabelsScope
	ScopeGmailSettings = gmail.GmailSettingsBasicScope
)

type GmailService struct {
	Svc         *gmail.Service
	UsersSvc    *gmail.UsersService
	SettingsSvc *gmail.UsersSettingsService
}

// NewGmailService creates a new collection of required gmail services.
func NewGmailService(ctx context.Context, options ...option.ClientOption) (*GmailService, error) {
	svc, err := gmail.NewService(ctx, options...)
	s := GmailService{
		Svc:         svc,
		UsersSvc:    gmail.NewUsersService(svc),
		SettingsSvc: gmail.NewUsersSettingsService(svc),
	}
	return &s, err
}
