package stickers

import (
	"fmt"

	admin "google.golang.org/api/admin/directory/v1"
)

const (
	ScopeAdminUserReadOnly = admin.AdminDirectoryUserReadonlyScope
	ScopeAdminUser         = admin.AdminDirectoryUserScope
)

// User is a wrapper for admin.User
type User admin.User

// ToUser
func ToUser(u *admin.User) *User {
	user := User(*u)
	return &user
}

// ToApi
func (u *User) ToApi() *admin.User {
	user := admin.User(*u)
	return &user
}

// ListUsers will read from the users api.
func (s *AdminService) ListUsers() ([]*User, error) {
	var pageToken string
	var users []*User

	for {
		// Set to max results, and page loop
		req := s.Svc.Users.List().
			Domain(GetDomain()).
			Projection("full").
			MaxResults(500).
			OrderBy("email")
		if pageToken != "" {
			req.PageToken(pageToken)
		}

		// Make the request
		r, err := req.Do()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve users: %v", err)
		}

		// Loop through users
		for _, u := range r.Users {
			user := User(*u)
			users = append(users, &user)
		}

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}

	return users, nil
}
