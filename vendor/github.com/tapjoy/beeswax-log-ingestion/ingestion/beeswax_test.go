package ingestion

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/tapjoy/beeswax-log-ingestion/beeswax"
)

const auctionID = "1234567890123456.123456789.12345.tapjoy" // Auction ID, Format: <auctionid.timestamp>.<auctionid.hostid>.<auctionid.tid>.<buzz_key>

func TestAuctionWon(t *testing.T) {
	defer resetClient(httpClient)
	defer resetDIS(identifyDevice)
	defer resetRecordError(logError)

	log.SetOutput(ioutil.Discard) // maybe put in a TestMain or init()

	const deviceID = "123"
	identifyDevice = func(ID string) (string, error) {
		if ID == deviceID {
			return "tjid." + deviceID, nil
		}
		return "", ErrNoDeviceIDFound
	}

	expParams := url.Values{
		"ad_view_id":      []string{auctionID},
		"dsp_creative_id": []string{"1234"},
		"bid_price":       []string{"1"},
		"clearing_price":  []string{"1"},
		"currency":        []string{"USD"},
		"udid":            []string{"tjid." + deviceID},
		"bid_time":        []string{"1234567890"},
	}

	// FIXME: may want to handle this differently
	// TODO: need to verify whether we rely to much on a beeswax model or expose too much
	validParser := func(data []byte) (beeswax.ImpressionMessage, error) {
		return beeswax.ImpressionMessage{
			AuctionID:     auctionID,
			CreativeID:    1234,
			BidPrice:      1,
			ClearingPrice: 1,
			Currency:      "USD",
			TJID:          "tjid." + deviceID,
			AuctionTime:   1234567890,
		}, nil
	}
	notImpressionParser := func(data []byte) (beeswax.ImpressionMessage, error) {
		return beeswax.ImpressionMessage{}, beeswax.ErrNotImpressionMessage
	}

	// maybe this case struct is doing too much
	cases := []struct {
		tag       string
		src       string
		method    func(string, []byte) *http.Request
		fwd       string
		parser    func([]byte) (beeswax.ImpressionMessage, error)
		expParams url.Values
		expInfo   string
		expErr    error
	}{
		{"happy path", "http://beeswax/test", newPostRequest, "http://cpm/test", validParser, expParams, "", nil},
		{"not impression message", "http://beeswax/test", newPostRequest, "http://cpm/test", notImpressionParser, url.Values{}, beeswax.ErrNotImpressionMessage.Error(), nil},

		// TODO: Add test cases
	}

	for _, c := range cases {
		t.Run(c.tag, func(t *testing.T) {
			r := c.method(c.src, []byte{})
			spy := urlSpy{c.expParams, nil}
			httpClient = newTestSpyClient(&spy)

			var recordedError error
			logError = func(msg string, w http.ResponseWriter, err error) {
				recordedError = err
			}
			var recordedInfo string
			logInfo = func(ctx string, r *http.Request, fwd string, w http.ResponseWriter) {
				recordedInfo = ctx
			}

			parser = c.parser
			auctionWonHandler(c.fwd, nil, r)

			if spy.queryErr != nil {
				t.Error(spy.queryErr)
			}
			if c.expInfo != recordedInfo {
				t.Errorf("got '%s', expected '%s'", recordedInfo, c.expInfo)
			}
			if c.expErr != recordedError {
				t.Errorf("got '%s', expected '%s'", recordedError, c.expErr)
			}
		})
	}
}
