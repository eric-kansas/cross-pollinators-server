package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/tapjoy/beeswax-log-ingestion/configs"
)

const disEndPoint = "api/v2/device_identifiers?"

// Http Client
var (
	timeout    = time.Duration(1 * time.Second)
	httpClient = &http.Client{Timeout: timeout}
)

// DISResponse is used to contain response data returned from
// device identity service
type DISRespose struct {
	InternalIdentifier string `json:"internal_identifier"`
	Key                string `json:"key"`
	Message            string `json:"message"`
}

// Identify makes a call to DIS to try to convert the passed in ID to
// Tapjoy Internal ID. If an error occurs it returns the originally passed in ID.
func Identify(deviceID string) (string, error) {
	data := url.Values{}
	data.Set("advertising_id", deviceID)

	var url bytes.Buffer
	url.WriteString(configs.Data.DISHostname)
	url.WriteString(disEndPoint)
	url.WriteString(data.Encode())

	req, err := http.NewRequest("POST", url.String(), nil)
	if err != nil {
		log.Printf("Error making new request to DIS: %v", err)
		return deviceID, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error failed to request to DIS: %v", err)
		return deviceID, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return deviceID, err
	}
	defer res.Body.Close()

	disRes := DISRespose{}
	json.Unmarshal(body, &disRes)

	if disRes.InternalIdentifier != "" {
		return disRes.InternalIdentifier, nil
	}

	return deviceID, fmt.Errorf("No InternalIdentifier from DIS. Message: %s", disRes.Message)
}
