package webhooks

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/frain-dev/convoy/pkg/verifier"
	"github.com/go-chi/chi/v5"
	"github.com/gobinary/go/pkg/mod/github.com/frain-dev/convoy@v0.7.0/pkg/verifier"
	log "github.com/sirupsen/logrus"
)

var (
	// GOOGLE_CLOUD_PROJECT is a user-set environment variable.
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

	// Retrieve Paystack webhooks secret.
	paystackSecret = os.Getenv("PAYSTACK_WEBHOOK_SECRET")

	// Retrieve Github webhooks secret.
	githubSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")

	pv = getConvoyVerifier()
)

// WebhookEndpoint is a HTTP Function to receive events from the world.
func WebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	// Build Router.
	router := chi.NewRouter()

	router.Route("/v1", func(v1Router chi.Router) {

		v1Router.Post("/webhooks/{provider}", WebhooksHandler)
	})

	// Serve Request.
	router.ServeHTTP(w, r)
}

func WebhooksHandler(w http.ResponseWriter, r *http.Request) {
	// Read Request.
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithError(err).Error("Bad Request: Could not read payload")
		return
	}

	// Verify Request.
	err = pv.VerifyRequest(r, payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.WithError(err).Error("Bad Request: Could not verify request")
		return
	}

	// Respond.
	w.Write([]byte("Event received."))

	// Perform business logic.
}

func getConvoyVerifier() verifier.Verifier {
	return verifier.NewHmacVerifier(&verifier.HmacOptions{
		Header:   "X-Convoy-Signature",
		Hash:     "SHA256",
		Secret:   convoySecret,
		Encoding: "hex",
	})
}
