package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mikemackintosh/stickers"
)

func main() {
	var port = ":8080"
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}

	if len(os.Getenv("STICKERS_CONFIG")) == 0 {
		log.Println(stickers.ErrMissingConfig{})
	}

	b, err := ioutil.ReadFile(os.Getenv("STICKERS_CONFIG"))
	if err != nil {
		log.Println(err)
	}

	// Load the configuration
	if err := stickers.LoadConfig(b); err != nil {
		log.Println(err)
		return
	}

	// Set the service account file
	if err := stickers.SetServiceAccountFile(os.Getenv("STICKERS_SVCACCT_JSON")); err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Create the client for Google API's
	ctx := context.Background()
	adminConfig, err := stickers.NewClientWithSubject(ctx, stickers.GetImpersonation(), []string{
		stickers.ScopeAdminUserReadOnly,
	})
	if err != nil {
		log.Println(err)
		return
	}

	gmail, err := stickers.NewAdminService(ctx, adminConfig)
	if err != nil {
		log.Println(err)
		return
	}

	// Get all google users
	users, err := gmail.ListUsers()
	if err != nil {
		log.Println(err)
		return
	}

	// Loop through the google users
	for _, user := range users {
		gmailConfig, err := stickers.NewClientWithSubject(ctx, user.PrimaryEmail, []string{
			stickers.ScopeGmailLabels,
			stickers.ScopeGmailSettings,
		})
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, user.PrimaryEmail)
		svc, err := stickers.NewGmailService(ctx, gmailConfig)
		if err != nil {
			log.Println(err)
			return
		}

		// Check for labels
		for _, label := range stickers.GetLabels() {
			l, err := svc.GetLabelForUser(label.Name, user)
			// Label does not exist, needs to be added
			if err != nil {
				log.Println(err)
				// Label does not exist, add it
				err := svc.CreateLabelForUser(label, user)
				if err != nil {
					log.Println(err)
				}
				continue
			}

			// Label exists, lets check the settings
			if err := svc.CompareLabels(l, label); err != nil {
				l.Name = label.Name
				l.Color = label.Color
				if err := svc.UpdateLabelForUser(l, user); err != nil {
					log.Println(err)
				}
			}
		}

		// Check for filters
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

				// Get the updated label
				l, err := svc.GetLabelForUser(filter.Label, user)
				if err != nil {
					log.Println(err)
				}

				// Set the filter action
				f.Action.AddLabelIds = append(f.Action.AddLabelIds, l.Id)

				// F exists, lets check the settings
				if err := svc.CreateFilterForUser(f, user); err != nil {
					log.Println(err)
				}

				continue
			}

			l, err := svc.GetLabelForUser(filter.Label, user)
			if err != nil {
				log.Println(err)
			}

			err = svc.CompareFilters(upstreamFilter, stickers.ToFilter(f), l.Id)
			if err != nil {
				if err := svc.DeleteFilterForUser(upstreamFilter.ToApi(), user); err != nil {
					log.Println(err)
				}

				f.Action.AddLabelIds = append(f.Action.AddLabelIds, l.Id)

				// F exists, lets check the settings
				if err := svc.CreateFilterForUser(f, user); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
