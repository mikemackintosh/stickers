package stickers

import (
	"context"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

// AdminService is a wrapper for admin.Service
type AdminService struct {
	Svc *admin.Service
}

// NewAdminService creates a new AdminService instance.
func NewAdminService(ctx context.Context, options ...option.ClientOption) (*AdminService, error) {
	svc, err := admin.NewService(ctx, options...)
	s := &AdminService{
		Svc: svc,
	}
	return s, err
}
