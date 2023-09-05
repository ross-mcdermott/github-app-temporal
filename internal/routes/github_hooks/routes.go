package github_hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// HTTP middleware setting a value on the request context
func Routes(r chi.Router) {
	r.Post("/github", process)

}

type status struct {
	Ok bool `json:"ok"`
}

func process(w http.ResponseWriter, r *http.Request) {

	println("got hook")
	data := status{
		Ok: true,
	}

	ctx := r.Context()

	r = r.WithContext(ctx)

	eventType := r.Header.Get("X-GitHub-Event")
	deliveryID := r.Header.Get("X-GitHub-Delivery")

	println(fmt.Sprintf("%s - %s", eventType, deliveryID))

	if eventType == "" {
		// not an event.
		return
	}

	// read to slice so can print as needed.
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	if eventType == "pull_request" {

		var p PullRequestWebook
		decoder := json.NewDecoder(bytes.NewReader(b))
		err = decoder.Decode(&p)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := processPullRequest(ctx, &p, deliveryID)

		if err != nil {
			println(err.Error())
			return
		}
	} else if eventType == "check_suite" {
		var p CheckSuiteWebhook
		decoder := json.NewDecoder(bytes.NewReader(b))
		err = decoder.Decode(&p)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := processCheckSuite(ctx, &p, deliveryID)

		if err != nil {
			println(err.Error())
			return
		}

	} else {

		fmt.Println("Unknown Event Type")
		fmt.Println(string(b))

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func processPullRequest(ctx context.Context, p *PullRequestWebook, deliveryID string) error {
	println("-- Pull Request --")
	println(fmt.Sprintf("Action: %s", p.Action))
	println(fmt.Sprintf("Repo: %s", p.Repository.FullName))

	if p.Action == "opened" {
		println("PR OPENED")
	}

	if p.Action == "reopened" {
		println("PR REOPENED")
	}

	if p.Action == "closed" {
		println("PR CLOSED")
	}

	return nil
}

func processCheckSuite(ctx context.Context, c *CheckSuiteWebhook, deliveryID string) error {
	println("-- Check Suite --")
	println(fmt.Sprintf("Action: %s", c.Action))
	println(fmt.Sprintf("Repo: %s", c.Repository.FullName))

	return nil
}
