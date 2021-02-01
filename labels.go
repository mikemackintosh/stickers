package stickers

import (
	"google.golang.org/api/gmail/v1"
)

// Label is a custom gmail.Label type.
type Label gmail.Label

// ToLabel will take a gmail.Label and convert to a local Label.
func ToLabel(l *gmail.Label) *Label {
	label := Label(*l)
	return &label
}

// ToApi will convert a local label to an API payload.
func (l *Label) ToApi() *gmail.Label {
	label := gmail.Label(*l)
	return &label
}

// GetLabelsForUser will return a list of labels for the specified user
func (s *GmailService) GetLabelsForUser(user *User) ([]*Label, error) {
	var labels []*Label
	r, err := s.Svc.Users.Labels.List(user.PrimaryEmail).Do()
	if err != nil {
		return nil, ErrFetchingLabelsForUser{user.PrimaryEmail}
	}

	for _, l := range r.Labels {
		labels = append(labels, ToLabel(l))
	}
	return labels, nil
}

// GetLabelForUser gets a specific label for the provided user.
func (s *GmailService) GetLabelForUser(name string, user *User) (*Label, error) {
	r, err := s.Svc.Users.Labels.List(user.PrimaryEmail).Do()
	if err != nil {
		return nil, ErrFetchingLabelsForUser{user.PrimaryEmail}
	}

	for _, l := range r.Labels {
		if name == l.Name {
			return ToLabel(l), nil
		}
	}

	return nil, ErrLabelNotFound{User: user.PrimaryEmail, Label: name}
}

// CreateLabelForUser will create a new user label.
func (s *GmailService) CreateLabelForUser(label *Label, user *User) error {
	_, err := s.UsersSvc.Labels.Create(user.PrimaryEmail, label.ToApi()).Do()
	return err
}

// CompareLabels compares an upstream and local label.
func (s *GmailService) CompareLabels(upstream, local *Label) error {
	if upstream.Name != local.Name {
		return ErrLabelMismatch{Label: local.Name}
	}

	// Sometimes color is nil, which means we need to check
	// to compare children fields.
	if upstream.Color != nil && local.Color != nil {
		if upstream.Color.TextColor != local.Color.TextColor {
			return ErrLabelMismatch{Label: local.Name}
		}

		if upstream.Color.BackgroundColor != local.Color.BackgroundColor {
			return ErrLabelMismatch{Label: local.Name}
		}
	}

	return nil
}

// UpdateLabelForUser will perform a patch label.
func (s *GmailService) UpdateLabelForUser(label *Label, user *User) error {
	_, err := s.UsersSvc.Labels.Patch(user.PrimaryEmail, label.Id, label.ToApi()).Do()
	return err
}
