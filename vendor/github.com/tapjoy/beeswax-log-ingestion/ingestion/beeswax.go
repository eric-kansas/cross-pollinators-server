package ingestion

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tapjoy/beeswax-log-ingestion/beeswax"
)

// Errors
var (
	ErrNotPostRequest = errors.New("Request method should be POST")
)

var parser = beeswax.Parse

// Handle beeswax log ingestion
// TODO: this handler should probably be moved into the beeswax package
func auctionWonHandler(endpoint string, w http.ResponseWriter, r *http.Request) {
	logRequest("auctionWonHandler", r, endpoint, w)

	if r.Method != http.MethodPost {
		logError("HTTP beacon request method was "+r.Method+" not GET", w, ErrNotPostRequest)
		return
	}

	// Parse body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logError("Failed to get data from http request body", w, err)
		return
	}

	// Parse beeswax impression protobuff
	auctionData, err := parser(data)

	if err == beeswax.ErrNotImpressionMessage {
		// TODO: add detail as to what message was returned (eg, Activity, Request, etc.)
		msg := fmt.Sprintf("%s", err)
		logInfo(msg, r, endpoint, w)
		return
	}
	if err != nil {
		logError("Failed to parse beeswax message", w, err)
		return
	}

	uri, err := buildAuctionWonURI(endpoint, auctionData)
	if err != nil {
		logError("Error building acution won uri", w, err)
		return
	}

	// pass on auction won message
	err = send(uri.String(), w, bytes.NewReader([]byte{}))
	if err != nil {
		logError("Error sending acution won messge to CPM", w, err)
	}
}

func buildAuctionWonURI(endpoint string, msg beeswax.ImpressionMessage) (*url.URL, error) {
	uri, err := url.Parse(endpoint)
	if err != nil {
		return &url.URL{}, err
	}

	query := url.Values{}
	query.Add("ad_view_id", msg.AuctionID)
	query.Add("dsp_creative_id", fmt.Sprintf("%d", msg.CreativeID))
	query.Add("bid_price", fmt.Sprintf("%d", msg.BidPrice))
	query.Add("clearing_price", fmt.Sprintf("%d", msg.ClearingPrice))
	query.Add("currency", msg.Currency)
	query.Add("udid", msg.TJID)
	query.Add("bid_time", strconv.FormatUint(msg.AuctionTime, 10))
	uri.RawQuery = query.Encode()

	return uri, nil
}
