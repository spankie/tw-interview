package cloudflareeth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
)

type rpcRequestBody struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type response struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  any    `json:"result"`
}

type httpClient struct {
	baseURL    string
	httpClient http.Client
}

func (c httpClient) Post(url string, body rpcRequestBody, res any) error {
	return c.doRequest(http.MethodPost, url, body, res)
}

func (c httpClient) doRequest(method, url string, body rpcRequestBody, dataRes any) error {
	dataBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("cannot marshal request body to json: %w", err)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.baseURL, url), bytes.NewReader(dataBytes))
	if err != nil {
		return fmt.Errorf("error occurred during request creation: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(dataRes)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	return nil
}
