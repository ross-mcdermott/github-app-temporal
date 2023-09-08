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
				err = handleCheckSuiteEvent(ctx, childLogger, event, client)
			case *github.CheckRunEvent:
				err = handleCheckRunEvent(ctx, childLogger, event, client)
			case *github.PullRequestEvent:
				err = handlePullRequestEvent(ctx, childLogger, event)
			default:
				err = UnknownEvent(ctx, childLogger, event)
			}

			if err != nil {
				childLogger.Error(err.Error())
			}
		})
	})

}

type status struct {
	Ok bool `json:"ok"`
}
