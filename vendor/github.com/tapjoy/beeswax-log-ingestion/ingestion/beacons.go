package ingestion

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/tapjoy/beeswax-log-ingestion/device"
)

// Errors
var (
	ErrNotGetRequest   = errors.New("Request method should be GET")
	ErrNoDeviceIDFound = errors.New("No 'advertising_id' found in url parameters")
)

var identifyDevice = device.Identify

func forwardBeacon(endpoint string, w http.ResponseWriter, r *http.Request) {
	logRequest("forwardBeacon", r, endpoint, w)

	if r.Method != http.MethodGet {
		logError("HTTP beacon request method was "+r.Method+" not GET", w, ErrNotGetRequest)
		return
	}

	tjid, err := fetchTJID(r.URL)
	if err != nil {
		logError("Failed to identify device", w, err)
		return
	}

	uri, err := appendUDID(endpoint, tjid, r.URL)
	if err != nil {
		logError("Failed to append TJID", w, err)
		return
	}

	send(uri.String(), w, bytes.NewReader([]byte{}))
}

func fetchTJID(uri *url.URL) (string, error) {
	log.Printf("fetchTJID: Parsing TJID from: %v", uri)
	advertiserID := uri.Query().Get("advertising_id")
	if advertiserID == "" {
		return "", ErrNoDeviceIDFound
	}

	log.Printf("fetchTJID: Fetching TJID for: %v", advertiserID)

	//TODO: https://jira.tapjoy.net/browse/BID-85
	tjid, err := identifyDevice(advertiserID)
	if err != nil {
		return "", err
	}
	log.Printf("fetchTJID: Success - fetched TJID: %v", tjid)
	return tjid, nil
}

func appendUDID(endpoint, tjid string, orig *url.URL) (*url.URL, error) {
	log.Printf("appendTJID: Building uri from end point: %v", endpoint)
	uri, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	q := orig.Query()
	q.Set("udid", tjid)
	uri.RawQuery = q.Encode() // add existing+tjid params to new uri
	log.Printf("appendTJID: Success - appended UDID/TJID: %v", uri)
	return uri, nil
}
