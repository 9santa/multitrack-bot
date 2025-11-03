package domain

import "time"

type TrackingResult struct {
	Number      string       `json:"number"`
	Courier     string       `json:"courier"`
	Status      string       `json:"status"`
	Description string       `json:"description"`
	Checkpoints []Checkpoint `json:"checkpoints"`
	LastUpdated time.Time    `json:"last_updated"`
}

type Checkpoint struct {
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
}

type RawTrackingResult struct {
	RawData    interface{} `json:"raw_data"`
	Courier    string      `json:"courier"`
	Successful bool        `json:"successful"`
	Error      string      `json:"error,omitempty"`
}

type HistoryRecord struct {
	Barcode  string `xml:"ItemParameters>Barcode"`
	OperDate string `xml:"OperationParameters>OperDate"`
	OperType string `xml:"OperationParameters>OperType>Name"`
	OperAttr string `xml:"OperationParameters>OperAttr>Name"`
	Address  string `xml:"AddressParameters>OperationAddress>Description"`
}
