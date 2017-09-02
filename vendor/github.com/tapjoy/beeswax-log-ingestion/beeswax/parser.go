package beeswax

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/tapjoy/beeswax-log-ingestion/device"
	beeswaxLogs "github.com/tapjoy/beeswax-log-ingestion/protos/beeswax/logs"
)

// Errors
var (
	ErrFailedToParseImpressionMessage = errors.New("Failed parse proto data")
	ErrNotImpressionMessage           = errors.New("Beeswax message was not impression message")
)

// ImpressionMessage ...
// TODO: this should be part of ingestion
type ImpressionMessage struct {
	AuctionID     string `json:"ad_view_id"`
	CreativeID    uint64 `json:"dsp_creative_id"`
	BidPrice      uint64 `json:"bid_price"`      // in microUSD
	ClearingPrice uint64 `json:"clearing_price"` // in microUSD
	Currency      string `json:"currency"`
	TJID          string `json:"tjid"`
	AuctionTime   uint64 `json:"time"` // epoch time
}

// Parse takes a beeswax protobuff log response, parsing it
// for user ID and win price and adding the Tapjoy internal ID
func Parse(data []byte) (ImpressionMessage, error) {
	beeswaxLogMessage := beeswaxLogs.AdLogMessage{}

	if err := proto.Unmarshal(data, &beeswaxLogMessage); err != nil {
		return ImpressionMessage{}, ErrFailedToParseImpressionMessage
	}

	log.Printf("Beeswax Ad Log Message: %+v", beeswaxLogMessage)

	if beeswaxLogMessage.Impression == nil {
		return ImpressionMessage{}, ErrNotImpressionMessage
	}

	tjid, err := convertToTapjoyInernalID(*beeswaxLogMessage.Impression.UserId)
	if err != nil {
		return ImpressionMessage{}, err
	}

	auctionTime, err := convertToRealEpoch(beeswaxLogMessage.Impression.GetRxTimestampUsecs())
	if err != nil {
		return ImpressionMessage{}, err
	}

	return ImpressionMessage{
		AuctionID:     beeswaxLogMessage.Impression.GetAuctionidStr(),
		CreativeID:    beeswaxLogMessage.Impression.GetCreativeId(),
		BidPrice:      beeswaxLogMessage.Impression.GetBidPriceMicrosUsd(),
		ClearingPrice: beeswaxLogMessage.Impression.GetWinPriceMicrosUsd(),
		Currency:      "USD",
		TJID:          tjid,
		AuctionTime:   auctionTime,
	}, nil
}

func convertToTapjoyInernalID(deviceID string) (string, error) {
	// Ignore string prefix before '.'
	if strings.Contains(deviceID, ".") {
		deviceID = strings.Split(deviceID, ".")[1]
	}

	deviceID, err := device.Identify(deviceID)

	if err != nil {
		return "", err
	}

	return deviceID, nil
}

// convertToRealEpoch is used to fix the epoch time we recieve from Beeswax.
// This function recalculates the timestamp to a correct Epoch time.
func convertToRealEpoch(unixTimestamp uint64) (uint64, error) {
	t := time.Unix(0, int64(unixTimestamp)*1000)
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return unixTimestamp, err
	}
	t2 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	l, err := time.LoadLocation("UTC")
	if err != nil {
		return unixTimestamp, err
	}
	realEpoch := t2.In(l).UnixNano() / 1000
	return uint64(realEpoch), nil
}
