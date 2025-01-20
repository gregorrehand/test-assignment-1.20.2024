package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     logrus.FieldLogger
}

func NewAPIClient(baseURL string, logger logrus.FieldLogger) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Logger: logger,
	}
}

func (api *APIClient) Get(ctx context.Context, endpoint string, token string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api.BaseURL+endpoint, nil)
	if err != nil {
		api.Logger.Error("Failed to create GET request: ", err)
		return err
	}

	api.addHeaders(req, token)

	return api.handleRequest(req, result)
}

func (api *APIClient) Post(ctx context.Context, endpoint string, token string, payload interface{}, result interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		api.Logger.Error("Failed to marshal POST payload: ", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api.BaseURL+endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		api.Logger.Error("Failed to create POST request: ", err)
		return err
	}
	api.addHeaders(req, token)
	return api.handleRequest(req, result)
}

func (api *APIClient) PostWithoutToken(ctx context.Context, endpoint string, payload interface{}, result interface{}) error {
	return api.Post(ctx, endpoint, "", payload, result)
}

func (api *APIClient) addHeaders(req *http.Request, token string) {
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")
}

func (api *APIClient) handleRequest(req *http.Request, result interface{}) error {
	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		api.Logger.Error("Request failed: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		api.Logger.Errorf("%s Request to %s failed with status: %s", req.Method, req.URL, resp.Status)
		return errors.New("Request failed")
	}

	return api.parseResponse(resp, result)
}

func (api *APIClient) parseResponse(resp *http.Response, result interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if result == nil {
		return nil // Void response expected
	}

	if err := json.Unmarshal(body, result); err != nil {
		return err
	}

	return nil
}
