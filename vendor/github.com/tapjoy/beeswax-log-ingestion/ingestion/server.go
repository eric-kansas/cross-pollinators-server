package ingestion

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"strings"

	"github.com/didip/tollbooth"
	newrelic "github.com/newrelic/go-agent"
)

// Errors
var (
	ErrNoHostProvided = errors.New("No host(s) provided")
)

var httpClient = &http.Client{
	Timeout: 1 * time.Second,
}

// TODO: standardize logging interface
var logError = recordError
var logInfo = logRequest

// TODO: http.Client a member of server? or should http.Server be a member of package?
type server struct {
	cpmAddr     string
	httpServer  *http.Server // possibly move into package like http.Client ???
	newrelicApp newrelic.Application
	rateLimit   int64
}

// ListenAndServe ...
func (s *server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// NewServer ...
// unexported server enforces use of constructor
func NewServer(host, cpm string, opts ...Option) (*server, error) {
	if host == "" || cpm == "" {
		return nil, ErrNoHostProvided
	}

	// default server
	srv := server{
		cpmAddr: sanitizeURI(cpm),
		httpServer: &http.Server{
			Addr:         sanitizeURI(host),
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
			IdleTimeout:  1 * time.Second,
		},
		newrelicApp: nil,
		rateLimit:   1,
	}

	// option(s) configuration
	for _, opt := range opts {
		err := opt(&srv)
		if err != nil {
			return nil, err
		}
	}

	// set handler
	srv.httpServer.Handler = registerHandlers(&srv)

	return &srv, nil
}

func registerHandlers(s *server) *http.ServeMux {
	mux := http.NewServeMux()

	registerV1(mux, s)

	// ops
	mux.HandleFunc(healthzHandler("/healthz", s))

	return mux
}

func registerV1(mux *http.ServeMux, s *server) {
	// exchange
	mux.HandleFunc(limitHandler("/exchange/v1/auction_won", "/api/services/v1/cpm/auction_won", s, auctionWonHandler))

	// beacons
	mux.HandleFunc(beaconHandler("/exchange/v1/impression", "/api/services/v1/cpm/impression", s))
	mux.HandleFunc(beaconHandler("/exchange/v1/video_start", "/api/services/v1/cpm/video_start", s))
	mux.HandleFunc(beaconHandler("/exchange/v1/video_tracking", "/api/services/v1/cpm/video_tracking", s))
	mux.HandleFunc(beaconHandler("/exchange/v1/video_complete", "/api/services/v1/cpm/video_complete", s))
	mux.HandleFunc(beaconHandler("/exchange/v1/clickthru", "/api/services/v1/cpm/clickthru", s))
}

// TODO: possibly set method receiver to ingestion.server
func limitHandler(pattern, cpmPattern string, s *server, handler func(string, http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
	cpmAddr := s.cpmAddr + cpmPattern
	f := func(w http.ResponseWriter, r *http.Request) { handler(cpmAddr, w, r) }
	limiter := tollbooth.LimitFuncHandler(tollbooth.NewLimiter(s.rateLimit, time.Nanosecond), f)

	if s.newrelicApp != nil {
		return newrelic.WrapHandleFunc(s.newrelicApp, pattern, limiter.ServeHTTP)
	}
	return pattern, limiter.ServeHTTP
}

// TODO: possibly set method receiver to ingestion.server
func beaconHandler(pattern, cpmPattern string, s *server) (string, func(http.ResponseWriter, *http.Request)) {
	cpmAddr := s.cpmAddr + cpmPattern
	f := func(w http.ResponseWriter, r *http.Request) { forwardBeacon(cpmAddr, w, r) }

	if s.newrelicApp != nil {
		return newrelic.WrapHandleFunc(s.newrelicApp, pattern, f)
	}
	return pattern, f
}

// TODO: possibly set method receiver to ingestion.server
func healthzHandler(pattern string, s *server) (string, func(http.ResponseWriter, *http.Request)) {
	f := func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "Success") }

	if s.newrelicApp != nil {
		return newrelic.WrapHandleFunc(s.newrelicApp, pattern, f)
	}
	return pattern, f
}

// TODO: send() a method of ingestion.server ?
func send(endPoint string, w http.ResponseWriter, body io.Reader) error {
	log.Printf("Send: GET %s", endPoint)
	req, err := http.NewRequest("GET", endPoint, body)
	if err != nil {
		logError("Error making "+endPoint+" request to CPM", w, err)
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logError("Error making "+endPoint+" request to CPM", w, err)
		return err
	}

	// responses are ignored
	log.Printf("Status: %s GET %s\n%v", resp.Status, endPoint, resp)
	return nil
}

func recordError(message string, w http.ResponseWriter, err error) {
	log.Printf("Error: %s %v", message, err)
	if txn, ok := w.(newrelic.Transaction); ok {
		txn.NoticeError(err)
	}
}

func sanitizeURI(addr string) string {
	return strings.TrimSuffix(addr, "/")
}

func logRequest(ctx string, r *http.Request, fwd string, w http.ResponseWriter) {
	dump, err := httputil.DumpRequest(r, false)
	if err != nil {
		logError("failed to debug request", w, err)
	}
	log.Printf("%s\n\t%q\n\tForward to: %s\n", ctx, dump, fwd)
}
