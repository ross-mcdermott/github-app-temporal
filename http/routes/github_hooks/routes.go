package github_hooks

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-github/github"
	temporal "go.temporal.io/sdk/client"
)

func RegisterRoutes(router chi.Router, logger *slog.Logger, client temporal.Client) {

	var route = "/hooks/github"

	logger.Debug(fmt.Sprintf("Registering '%s' route.", route))

	router.Group(func(r chi.Router) {

		r.Post(route, func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			r = r.WithContext(ctx)

			payload, err := github.ValidatePayload(r, []byte("0695679902"))
			if err != nil {
				logger.Error("Unable to validate payload")
				return
			}
			event, err := github.ParseWebHook(github.WebHookType(r), payload)
			if err != nil {
				logger.Error("Unable to parse payload")
				return
			}

			childLogger := logger.With(
				slog.String("delivery_id", r.Header.Get("X-GitHub-Delivery")),
			)

			// jsonData, err := json.Marshal(event)

			// println(string(jsonData))

			switch event := event.(type) {
			case *github.CheckSuiteEvent:
				HandleCheckSuiteEvent(ctx, childLogger, event, client)
			case *github.CheckRunEvent:
				HandleCheckRunEvent(ctx, childLogger, event, client)
			case *github.PullRequestEvent:
				HandlePullRequestEvent(ctx, childLogger, event)
			default:
				UnknownEvent(ctx, childLogger, event)
			}
		})
	})

}

type status struct {
	Ok bool `json:"ok"`
}
