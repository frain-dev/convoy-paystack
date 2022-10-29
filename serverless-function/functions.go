package webhooks

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	client "github.com/frain-dev/convoy-go"
	"github.com/frain-dev/convoy/pkg/verifier"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

var (
	// Retrieve convoy webhooks secret.
	convoySecret = os.Getenv("CONVOY_WEBHOOK_SECRET")

	// Retrieve Hash
	hash = os.Getenv("CONVOY_HASH")

	// Retrieve Header
	header = os.Getenv("CONVOY_HEADER")

	// Configure Paystack Verifier
	pv = getPaystackVerifier()

	// Configure Convoy Verifier
	cv = getConvoyVerifier()
)

// WebhookEndpoint is a HTTP Function to receive events from the world.
func WebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	// Build Router.
	router := chi.NewRouter()

	router.Route("/v1", func(v1Router chi.Router) {

		v1Router.Post("/webhooks/simple", SimpleWebhooksHandler)
		v1Router.Post("/webhooks/advanced", AdvancedWebhooksHandler)
	})

	// Serve Request.
	router.ServeHTTP(w, r)
}

func SimpleWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	// Log Request.
	log.Info(r.Header)

	// Read Request.
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	// Verify Request.
	err = pv.VerifyRequest(r, payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	// Respond.
	w.Write([]byte("Event received."))

	// Perform business logic.
}

func AdvancedWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	// Log Request.
	log.Info(r.Header)

	// Read Request.
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	// Verify Request
	cv.Payload = payload
	err = cv.Verify()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	// Respond.
	w.Write([]byte("Event received."))

	// Perform business logic.
}

func getPaystackVerifier() verifier.Verifier {
	return verifier.NewHmacVerifier(&verifier.HmacOptions{
		Header:   header,
		Hash:     hash,
		Secret:   convoySecret,
		Encoding: "hex",
	})
}

func getConvoyVerifier() *client.Webhook {
	opts := client.ConfigOpts{
		SigHeader: header,
		Hash:      hash,
		Secret, convoySecret,
		IsAdvanced: true,
		Encoding:   client.HexEncoding,
		Tolerance:  5 * time.Duration,
	}
	return &client.NewWebhook(&opts)
}
