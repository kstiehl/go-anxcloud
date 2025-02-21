package zone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type RecordRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	RData  string `json:"rdata"`
	Region string `json:"region"`
	TTL    int    `json:"ttl,omitempty"`
}

// ListRecords API method
func (a api) ListRecords(ctx context.Context, zone string) ([]Record, error) {
	url := fmt.Sprintf(
		"%s%s/%s/records",
		a.client.BaseURL(),
		pathPrefix,
		zone,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create record list request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute record list request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return nil, fmt.Errorf("could not execute record list request, got response %s", httpResponse.Status)
	}

	responsePayload := make([]Record, 0)
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not decode zone list response: %w", err)
	}

	return responsePayload, nil
}

// NewRecord new record API method
func (a api) NewRecord(ctx context.Context, zone string, record RecordRequest) (Zone, error) {
	url := fmt.Sprintf(
		"%s%s/%s/records",
		a.client.BaseURL(),
		pathPrefix,
		zone,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(record); err != nil {
		panic(fmt.Sprintf("could not create request data for create zone: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestData)
	if err != nil {
		return Zone{}, fmt.Errorf("could not create record create request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Zone{}, fmt.Errorf("could not execute record create request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Zone{}, fmt.Errorf("could not execute record create request, got response %s", httpResponse.Status)
	}

	var responsePayload Zone
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Zone{}, fmt.Errorf("could not decode record create response: %w", err)
	}

	return responsePayload, nil
}

// UpdateRecord record API method
func (a api) UpdateRecord(ctx context.Context, zone string, id uuid.UUID, record RecordRequest) (Zone, error) {
	url := fmt.Sprintf(
		"%s%s/%s/records/%s",
		a.client.BaseURL(),
		pathPrefix,
		zone,
		id,
	)

	requestData := bytes.Buffer{}
	if err := json.NewEncoder(&requestData).Encode(record); err != nil {
		panic(fmt.Sprintf("could not create request data for update zone: %v", err))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, &requestData)
	if err != nil {
		return Zone{}, fmt.Errorf("could not create record update request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return Zone{}, fmt.Errorf("could not execute record update request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return Zone{}, fmt.Errorf("could not execute record update request, got response %s", httpResponse.Status)
	}

	var responsePayload Zone
	err = json.NewDecoder(httpResponse.Body).Decode(&responsePayload)
	_ = httpResponse.Body.Close()
	if err != nil {
		return Zone{}, fmt.Errorf("could not decode record update response: %w", err)
	}

	return responsePayload, nil
}

// DeleteRecord record API method
func (a api) DeleteRecord(ctx context.Context, zone string, id uuid.UUID) error {
	url := fmt.Sprintf(
		"%s%s/%s/records/%s",
		a.client.BaseURL(),
		pathPrefix,
		zone,
		id,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("could not create record delete request: %w", err)
	}

	httpResponse, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute record delete request: %w", err)
	}
	if httpResponse.StatusCode >= 500 && httpResponse.StatusCode < 600 {
		return fmt.Errorf("could not execute record delete request, got response %s", httpResponse.Status)
	}

	return httpResponse.Body.Close()
}
