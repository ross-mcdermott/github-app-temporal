package github_hooks

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, logger *slog.Logger) {

	var route = "/hooks/github"

	logger.Debug(fmt.Sprintf("Registering '%s' route.", route))

	router.Group(func(r chi.Router) {

		r.Post(route, func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			r = r.WithContext(ctx)

			eventType := r.Header.Get("X-GitHub-Event")
			deliveryID := r.Header.Get("X-GitHub-Delivery")

			if eventType == "" {
				logger.Error("Webhook missing 'X-GitHub-Event' header")
			}

			if deliveryID == "" {
				logger.Error("Webhook missing 'X-GitHub-Delivery' header")
			}

			// track the event and delivery id generically.
			logger := slog.New(logger.Handler()).WithGroup("github").With(
				slog.String("event", eventType),
				slog.String("delivery_id", deliveryID),
			)

			logger.Debug("Recieved GitHub WebHook event")

			// read to slice from stream
			b, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}

			if eventType == "check_suite" {
				err = HandleCheckSuiteEvent(ctx, logger, b, deliveryID)

				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			if eventType == "check_run" {
				time.Sleep(5 * time.Second)
				//println(string(b))
				err = HandleCheckRunEvent(ctx, logger, b, deliveryID)

				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			data := status{
				Ok: true,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)

		})
	})

}

type status struct {
	Ok bool `json:"ok"`
}

// func process(w http.ResponseWriter, r *http.Request) {

// 	println("got hook")
// 	data := status{
// 		Ok: true,
// 	}

// 	ctx := r.Context()

// 	r = r.WithContext(ctx)

// 	eventType := r.Header.Get("X-GitHub-Event")
// 	deliveryID := r.Header.Get("X-GitHub-Delivery")

// 	println(fmt.Sprintf("%s - %s", eventType, deliveryID))

// 	if eventType == "" {
// 		// not an event.
// 		return
// 	}

// 	// read to slice so can print as needed.
// 	b, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		return
// 	}

// 	if eventType == "pull_request" {

// 		var p PullRequestWebook
// 		decoder := json.NewDecoder(bytes.NewReader(b))
// 		err = decoder.Decode(&p)

// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		err := processPullRequest(ctx, &p, deliveryID)

// 		if err != nil {
// 			println(err.Error())
// 			return
// 		}
// 	} else if eventType == "check_suite" {

// 		var p CheckSuiteWebhook
// 		decoder := json.NewDecoder(bytes.NewReader(b))
// 		err = decoder.Decode(&p)

// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		err := processCheckSuite(ctx, &p, deliveryID)

// 		if err != nil {
// 			println(err.Error())
// 			return
// 		}

// 	} else {

// 		fmt.Println("Unknown Event Type")
// 		fmt.Println(string(b))

// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(data)
// }

// func processPullRequest(ctx context.Context, p *PullRequestWebook, deliveryID string) error {
// 	println("-- Pull Request --")
// 	println(fmt.Sprintf("Action: %s", p.Action))
// 	println(fmt.Sprintf("Repo: %s", p.Repository.FullName))

// 	if p.Action == "opened" {
// 		println("PR OPENED")
// 	}

// 	if p.Action == "reopened" {
// 		println("PR REOPENED")
// 	}

// 	if p.Action == "closed" {
// 		println("PR CLOSED")
// 	}

// 	return nil
// }

// func processCheckSuite(ctx context.Context, c *CheckSuiteWebhook, deliveryID string) error {
// 	println("-- Check Suite --")
// 	println(fmt.Sprintf("Action: %s", c.Action))
// 	println(fmt.Sprintf("Repo: %s", c.Repository.FullName))

// 	// when requested - need to register the check.
// 	// Action: requested
// 	// kick off a temporal workflow at this point to allow the creation of the suite.

// 	return nil
// }
