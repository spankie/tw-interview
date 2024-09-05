package cloudflareeth

import (
	"errors"
	"fmt"

	"github.com/spankie/tw-interview/blockparser/blockchain"
)

var (
	ErrInvalidBlockResponse = errors.New("invalid block response")
)

const (
	ethBlockNumberMethod      = "eth_blockNumber"
	ethGetBlockByNumberMethod = "eth_getBlockByNumber"

	cloudFlareBaseURL = "https://cloudflare-eth.com"
)

type Client struct {
	client         httpClient
	jsonRPCVersion string
}

// NewCloudflareEthClient creates a new cloudflare eth client.
func NewCloudflareEthClient() *Client {
	return &Client{
		client:         newHTTPClient(cloudFlareBaseURL),
		jsonRPCVersion: "2.0",
	}
}

func (c Client) GetLatestBlock() (string, error) {
	requestBody := rpcRequestBody{
		Jsonrpc: c.jsonRPCVersion,
		Method:  ethBlockNumberMethod,
		Params:  []any{},
		ID:      83,
	}

	var res response

	err := c.client.Post("", requestBody, &res)
	if err != nil {
		return "", fmt.Errorf("error making post request: %w", err)
	}

	blockNumberStr, ok := res.Result.(string)
	if !ok {
		return "", ErrInvalidBlockResponse
	}

	return blockNumberStr, nil
}

// getBlock queries the etheruem blockchain to the block identified by the blockNumber
// represented in hex.
func (c Client) GetBlock(blockNumber string) (*blockchain.Block, error) {
	rpcReq := rpcRequestBody{
		Jsonrpc: c.jsonRPCVersion,
		Method:  ethGetBlockByNumberMethod,
		Params:  []any{blockNumber, true},
		ID:      1,
	}

	res := &response{Result: &blockchain.Block{}}

	err := c.client.Post("", rpcReq, res)
	if err != nil {
		return nil, fmt.Errorf("http error getting block #%s: %w", blockNumber, err)
	}

	block, ok := res.Result.(*blockchain.Block)
	if !ok {
		return nil, fmt.Errorf("type mismatch: %w", err)
	}

	return block, nil
}
