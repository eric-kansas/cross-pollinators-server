package ingestion

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestForwardBeacon(t *testing.T) {
	defer resetClient(httpClient)
	defer resetDIS(identifyDevice)
	defer resetRecordError(logError)

	log.SetOutput(ioutil.Discard) // maybe put in a TestMain or init()

	deviceID := "123"
	identifyDevice = func(ID string) (string, error) {
		if ID == deviceID {
			return "tjid." + deviceID, nil
		}
		return "", ErrNoDeviceIDFound
	}

	expParams := url.Values{
		"foo":            []string{"bar"},
		"advertising_id": []string{deviceID},
		"udid":           []string{"tjid." + deviceID},
	}

	cases := []struct {
		tag       string
		src       string
		method    func(string, []byte) *http.Request
		fwd       string
		query     string
		body      []byte
		expParams url.Values
		expErr    error
	}{
		{"happy path", "http://beacon/test", newGetRequest, "http://cpm/test", "foo=bar&advertising_id=123", []byte{}, expParams, nil},

		{"no advertising ID provided", "http://beacon/test", newGetRequest, "http://cpm/test", "foo=bar", []byte{}, expParams, ErrNoDeviceIDFound},
		{"TJID not found", "http://beacon/test", newGetRequest, "http://cpm/test", "foo=bar&advertising_id=456", []byte{}, expParams, ErrNoDeviceIDFound},

		// technically doesn't invalidate against: UPDATE, PATCH, etc.
		{"a POST request", "http://beacon/test", newPostRequest, "http://cpm/test", "foo=bar&advertising_id=123", []byte{}, expParams, ErrNotGetRequest},
	}

	for _, c := range cases {
		t.Run(c.tag, func(t *testing.T) {
			r := c.method(fmt.Sprintf("%s?%s", c.src, c.query), c.body)
			spy := urlSpy{c.expParams, nil}
			httpClient = newTestSpyClient(&spy)

			var recordedError error
			logError = func(msg string, w http.ResponseWriter, err error) {
				recordedError = err
			}

			forwardBeacon(c.fwd, nil, r)

			if spy.queryErr != nil {
				t.Error(spy.queryErr)
			}
			if c.expErr != recordedError {
				t.Errorf("got '%s', expected '%s'", recordedError, c.expErr)
			}
		})
	}
}

func newGetRequest(uri string, body []byte) *http.Request {
	req, _ := http.NewRequest("GET", uri, bytes.NewReader(body))
	return req
}

func newPostRequest(uri string, body []byte) *http.Request {
	req, _ := http.NewRequest("POST", uri, bytes.NewReader(body))
	return req
}

type urlSpy struct {
	expParams url.Values
	queryErr  error
}

func (spy *urlSpy) RoundTrip(r *http.Request) (*http.Response, error) {
	if !cmp.Equal(spy.expParams, r.URL.Query()) {
		spy.queryErr = fmt.Errorf("got %+v, expected %+v", r.URL.Query(), spy.expParams)
		return nil, spy.queryErr
	}

	return &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
}

func newTestSpyClient(spy *urlSpy) *http.Client {
	return &http.Client{
		Transport: spy,
	}
}

func resetClient(orig *http.Client) {
	httpClient = orig
}

func resetDIS(orig func(string) (string, error)) {
	identifyDevice = orig
}

func resetRecordError(orig func(message string, w http.ResponseWriter, err error)) {
	logError = orig
}
