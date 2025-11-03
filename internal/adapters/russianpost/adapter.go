package russianpost

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"multitrack-bot/internal/domain"
	"net/http"
	"time"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	Response OperationHistoryResponse `xml:"getOperationHistoryResponse"`
}

type OperationHistoryResponse struct {
	Data OperationHistoryData `xml:"OperationHistoryData"`
}

type OperationHistoryData struct {
	Records []HistoryRecord `xml:"historyRecord"`
}

type HistoryRecord struct {
	Barcode  string `xml:"ItemParameters>Barcode"`
	OperDate string `xml:"OperationParameters>OperDate"`
	OperType string `xml:"OperationParameters>OperType>Name"`
	OperAttr string `xml:"OperationParameters>OperAttr>Name"`
	Address  string `xml:"AddressParameters>OperationAddress>Description"`
}

type RussianPostAdapter struct {
	client   *http.Client
	login    string
	password string
	endpoint string
}

func NewRussianPostAdapter(login, password string) *RussianPostAdapter {
	return &RussianPostAdapter{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		login:    login,
		password: password,
		endpoint: "https://tracking.russianpost.ru/rtm34",
	}
}

func (a *RussianPostAdapter) Name() string {
	return "russianpost"
}

func (a *RussianPostAdapter) Validate(trackingNumber string) bool {
	// basic validation
	return len(trackingNumber) == 14
}

func (a *RussianPostAdapter) Track(ctx context.Context, trackingNumber string) (*domain.RawTrackingResult, error) {
	// log.Printf("[Track] courier=%s number=%s", a.Name(), trackingNumber)

	requestXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope"
               xmlns:oper="http://russianpost.org/operationhistory"
               xmlns:data="http://russianpost.org/operationhistory/data"
               xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soap:Header/>
   <soap:Body>
      <oper:getOperationHistory>
         <data:OperationHistoryRequest>
            <data:Barcode>%s</data:Barcode>
            <data:MessageType>0</data:MessageType>
            <data:Language>RUS</data:Language>
         </data:OperationHistoryRequest>
         <data:AuthorizationHeader soapenv:mustUnderstand="1">
            <data:login>%s</data:login>
            <data:password>%s</data:password>
         </data:AuthorizationHeader>
      </oper:getOperationHistory>
   </soap:Body>
</soap:Envelope>`, trackingNumber, a.login, a.password)

	req, err := http.NewRequestWithContext(ctx, "POST", a.endpoint, bytes.NewBufferString(requestXML))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// ⚠️ SOAP 1.2 обязательный заголовок (иначе 415)
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var soapResp Envelope
	if err := xml.Unmarshal(body, &soapResp); err != nil {
		log.Printf("[DEBUG SOAP RAW RESPONSE]\n%s\n", string(body))
		return nil, fmt.Errorf("unmarshal SOAP: %w", err)
	}

	records := soapResp.Body.Response.Data.Records
	if len(records) == 0 {
		return &domain.RawTrackingResult{
			Courier:    a.Name(),
			Successful: false,
			Error:      "no tracking history data",
		}, nil
	}

	// Преобразуем в RawTrackingResult
	raw := &domain.RawTrackingResult{
		Courier:    "Почта России",
		RawData:    records,
		Successful: true,
	}

	return raw, nil
}
