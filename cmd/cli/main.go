package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mikemackintosh/stickers"
)

var (
	flagConfig      string
	flagCredentials string
	flagLocale      string
	flagVerbosity   bool
)

func init() {
	flag.StringVar(&flagConfig, "f", "stickers.yaml", "Stickers configuration file (yaml)")
	flag.StringVar(&flagLocale, "l", "colors", "")
	flag.StringVar(&flagCredentials, "c", "service-account.json", "Service account credentials")
	flag.BoolVar(&flagVerbosity, "v", false, "Verbosity mode")
}

func main() {
	flag.Parse()

	var lang LangPack
	lang, ok := language[flagLocale]
	if !ok {
		fmt.Printf("Please use one of the following locales:\n")
		for k := range language {
			if k == "colors" {
				k = k + " (default)"
			}
			fmt.Printf("\t- %s\n", k)
		}
		os.Exit(1)
	}

	if flagVerbosity {
		fmt.Printf(lang[HEADER])
		fmt.Printf(lang[FLAG_CONFIG], flagConfig)
		fmt.Printf(lang[FLAG_CREDENTIALS], flagCredentials)
		fmt.Printf(lang[FLAG_VERBOSITY], flagVerbosity)
		fmt.Printf(lang[MARKER])
		fmt.Printf(lang[READING_CONFIG], flagConfig)
	}

	// Read the configuration file to []byte
	b, err := ioutil.ReadFile(flagConfig)
	if err != nil {
		fmt.Printf(lang[ERROR], err)
		return
	}

	// Load the configuration
	if err := stickers.LoadConfig(b); err != nil {
		fmt.Printf(lang[ERROR], err)
		return
	}

	if flagVerbosity {
		fmt.Printf(lang[FOUND_CONFIGURATION], stickers.GetConfig())
	}

	// Set the service account file
	if err := stickers.SetServiceAccountFile(flagCredentials); err != nil {
		fmt.Printf(lang[ERROR], err)
		return
	}

	// Create the client for Google API's
	ctx := context.Background()
	adminConfig, err := stickers.NewClientWithSubject(ctx, stickers.GetImpersonation(), []string{
		stickers.ScopeAdminUserReadOnly,
	})
	if err != nil {
		fmt.Printf(lang[ERROR], err)
		return
	}

	gmail, err := stickers.NewAdminService(ctx, adminConfig)
	if err != nil {
		fmt.Printf(lang[ERROR], err)
		return
	}

	// Get all google users
	users, err := gmail.ListUsers()
	if err != nil {
		fmt.Printf(lang[ERROR], err)
		return
	}

	// Loop through the google users
	for _, user := range users {
		gmailConfig, err := stickers.NewClientWithSubject(ctx, user.PrimaryEmail, []string{
			stickers.ScopeGmailLabels,
			stickers.ScopeGmailSettings,
		})
		if err != nil {
			fmt.Printf(lang[ERROR], err)
			return
		}

		fmt.Printf(lang[USER], user.PrimaryEmail)
		svc, err := stickers.NewGmailService(ctx, gmailConfig)
		if err != nil {
			fmt.Printf(lang[ERROR], err)
			return
		}

		// Check for labels
		fmt.Printf(lang[TITLE], "labels")
		for _, label := range stickers.GetLabels() {
			l, err := svc.GetLabelForUser(label.Name, user)
			// Label does not exist, needs to be added
			if err != nil {
				fmt.Printf(lang[MISSING], label.Name)
				// Label does not exist, add it
				err := svc.CreateLabelForUser(label, user)
				if err != nil {
					fmt.Printf(lang[ERROR], err)
				}
				continue
			}

			fmt.Printf(lang[FOUND], label.Name)
			// Label exists, lets check the settings
			if err := svc.CompareLabels(l, label); err != nil {
				fmt.Printf(lang[MISMATCH])
				l.Name = label.Name
				l.Color = label.Color
				if err := svc.UpdateLabelForUser(l, user); err != nil {
					fmt.Printf(lang[ERROR], err)
				}
			}
		}

		// Check for filters
		fmt.Printf(lang[TITLE], "filters")
		for _, filter := range stickers.GetFilters() {

			// Let's create a new gmail filter.
			f := stickers.NewGmailFilter()
			if filter.Query != nil {
				q := fmt.Sprintf("%s OR %s", filter.Label, filter.Query.ToString())
				f.Criteria.Query = q
			}

			// Check for an existing filter.
			upstreamFilter, err := svc.GetFilterForUser(filter.Label, user)

			// If it's not found
			if err != nil {
				fmt.Printf(lang[MISSING], filter.Label)

				// Get the updated label
				l, err := svc.GetLabelForUser(filter.Label, user)
				if err != nil {
					fmt.Printf(lang[ERROR], err)
				}

				// Set the filter action
				f.Action.AddLabelIds = append(f.Action.AddLabelIds, l.Id)

				// F exists, lets check the settings
				if err := svc.CreateFilterForUser(f, user); err != nil {
					fmt.Printf(lang[ERROR], err)
				}

				continue
			}

			fmt.Printf(lang[FOUND], filter.Label)

			l, err := svc.GetLabelForUser(filter.Label, user)
			if err != nil {
				fmt.Printf(lang[ERROR], err)
			}

			if err := svc.CompareFilters(upstreamFilter, stickers.ToFilter(f), l.Id); err != nil {
				fmt.Printf(lang[DELETING], filter.Query.ToString())
				if err := svc.DeleteFilterForUser(upstreamFilter.ToApi(), user); err != nil {
					fmt.Printf(lang[ERROR], err)
				}

				f.Action.AddLabelIds = append(f.Action.AddLabelIds, l.Id)

				// F exists, lets check the settings
				if err := svc.CreateFilterForUser(f, user); err != nil {
					fmt.Printf(lang[ERROR], err)
				} else {
					fmt.Printf(lang[UPDATING], filter.Label)
				}
			}
		}
	}
}
