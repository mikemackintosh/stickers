package stickers

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// instance variable
var config = &Config{}

// Labels is a collection of Label.
type Labels []*Label

// LocalFilters is a collection of Filter.
type ConfigFilters []ConfigFilter

// ConfigFilter represents a gmail filter.
type ConfigFilter struct {
	Action string       `yaml:"action"`
	Label  string       `yaml:"label"`
	Query  *FilterQuery `yaml:"query,omitempty"`
}

// FilterQuery is a query object for filters.
type FilterQuery struct {
	From    *[]string `yaml:"from,omitempty"`
	NotFrom *[]string `yaml:"not_from,omitempty"`
	To      *[]string `yaml:"to,omitempty"`
	NotTo   *[]string `yaml:"not_to,omitempty"`
}

// ToString will convert a query string to a compatible filter.
func (f *FilterQuery) ToString() string {
	var q string
	if f.From != nil {
		q = fmt.Sprintf("from:{%s}", strings.Join(*f.From, " "))
	}

	if f.NotFrom != nil {
		q += fmt.Sprintf("-from:{%s}", strings.Join(*f.NotFrom, " "))
	}

	if f.To != nil {
		q += fmt.Sprintf("to:{%s}", strings.Join(*f.To, " "))
	}

	if f.NotTo != nil {
		q += fmt.Sprintf("-to:{%s}", strings.Join(*f.NotTo, " "))
	}

	return q
}

// Config is a standard config struct.
type Config struct {
	Labels      Labels        `yaml:"labels"`
	Filters     ConfigFilters `yaml:"filters"`
	Domain      string        `yaml:"domain"`
	Impersonate string        `yaml:"impersonate"`
}

// ErrMissingConfig is used to
type ErrMissingConfig struct{}

// Error
func (e ErrMissingConfig) Error() string {
	return "config could not be found"
}

// LoadConfig will read the configuration from yaml.
func LoadConfig(b []byte) error {
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		return err
	}

	return nil
}

// GetConfig will return the active config.
func GetConfig() *Config {
	return config
}

// GetDomain returns a string of the configured domain.
func GetDomain() string {
	return config.Domain
}

// GetImpersonation returns a string of the user to impersonate.
func GetImpersonation() string {
	return config.Impersonate
}

// GetLabels returns a collection of labels to configure.
func GetLabels() Labels {
	return config.Labels
}

// GetFilters returns a collection of filters to configure.
func GetFilters() ConfigFilters {
	return config.Filters
}
