package stickers

import (
	"strings"

	"google.golang.org/api/gmail/v1"
)

// Filter is a custom gmail.Filter type.
type Filter gmail.Filter

// ToFilter will take a gmail.Filter and convert to a local Filter.
func ToFilter(f *gmail.Filter) *Filter {
	filter := Filter(*f)
	return &filter
}

// ToApi will convert a local label to an API payload.
func (f *Filter) ToApi() *gmail.Filter {
	filter := gmail.Filter(*f)
	return &filter
}

// NewGmailFilter creates a new Gmail filter.
func NewGmailFilter() *gmail.Filter {
	f := &gmail.Filter{
		Criteria: &gmail.FilterCriteria{},
		Action:   &gmail.FilterAction{},
	}
	return f
}

// GetFiltersForUser will return filters for the user.
func (s GmailService) GetFiltersForUser(user *User) ([]*Filter, error) {
	var filters []*Filter
	r, err := s.SettingsSvc.Filters.List(user.PrimaryEmail).Do()
	if err != nil {
		return nil, ErrFetchingFiltersForUser{user.PrimaryEmail}
	}

	for _, l := range r.Filter {
		filters = append(filters, ToFilter(l))
	}
	return filters, nil
}

// GetLabelForUser gets a specific label for the provided user.
func (s *GmailService) GetFilterForUser(filterKey string, user *User) (*Filter, error) {
	r, err := s.SettingsSvc.Filters.List(user.PrimaryEmail).Do()
	if err != nil {
		return nil, ErrFetchingFiltersForUser{user.PrimaryEmail}
	}

	for _, f := range r.Filter {
		if strings.Contains(f.Criteria.Query, filterKey) {
			return ToFilter(f), nil
		}
	}

	return nil, ErrFilterNotFound{User: user.PrimaryEmail, Filter: filterKey}
}

// CompareFilters compares an upstream and local filters.
func (s *GmailService) CompareFilters(upstream, local *Filter, labelId string) error {
	if upstream.Action != nil {
		var found bool
		for _, label := range upstream.Action.AddLabelIds {
			if label == labelId {
				found = true
			}
		}
		if !found {
			return ErrFilterMismatch{Filter: upstream.Id}
		}
	}

	if upstream.Criteria != nil {
		if upstream.Criteria.Query != local.Criteria.Query {
			return ErrFilterMismatch{Filter: upstream.Id}
		}
	}

	return nil
}

// CreateFilterForUser creates a filter for the user.
func (s *GmailService) CreateFilterForUser(f *gmail.Filter, user *User) error {
	_, err := s.SettingsSvc.Filters.Create(user.PrimaryEmail, f).Do()
	return err
}

// DeleteFilterForUser deletes a filter for the user.
func (s *GmailService) DeleteFilterForUser(f *gmail.Filter, user *User) error {
	err := s.SettingsSvc.Filters.Delete(user.PrimaryEmail, f.Id).Do()
	return err
}
