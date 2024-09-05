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

// TODO: use generics for the result type.
type response struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  any    `json:"result"`
}

type httpClient interface {
	Post(url string, body rpcRequestBody, res any) error
}

type client struct {
	baseURL    string
	httpClient http.Client
}

// newHTTPClient creates a new http client with a base url.
func newHTTPClient(baseURL string) httpClient {
	c := *http.DefaultClient
	c.Timeout = defaultTimeout

	return client{baseURL: baseURL, httpClient: c}
}

func (c client) Post(url string, body rpcRequestBody, res any) error {
	return c.doRequest(http.MethodPost, url, body, res)
}

func (c client) doRequest(method, url string, body rpcRequestBody, dataRes any) error {
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
