package cloudflareeth

import (
	"fmt"

	"github.com/spankie/tw-interview/blockparser/blockchain"
)

const (
	ethBlockNumberMethod      = "eth_blockNumber"
	ethGetBlockByNumberMethod = "eth_getBlockByNumber"

	cloudFlareBaseUrl = "https://cloudflare-eth.com"
)

// TODO: maybe rename this type
type cloudflareEth struct {
	client         httpClient
	jsonRpcVersion string
}

// NewCloudflareEth creates a new cloudflare eth client
func NewCloudflareEth() *cloudflareEth {
	return &cloudflareEth{
		client:         newHTTPClient(cloudFlareBaseUrl),
		jsonRpcVersion: "2.0",
	}
}

func (c cloudflareEth) GetLatestBlock() (string, error) {
	requestBody := rpcRequestBody{
		Jsonrpc: c.jsonRpcVersion,
		Method:  ethBlockNumberMethod,
		Params:  []any{},
		ID:      83,
	}

	res := &response{}

	err := c.client.Post("", requestBody, res)
	if err != nil {
		return "", fmt.Errorf("error making post request: %w", err)
	}

	blockNumberStr, ok := res.Result.(string)
	if !ok {
		return "", fmt.Errorf("get block response is invalid")
	}

	return blockNumberStr, nil
}

// getBlock queries the etheruem blockchain to the block identified by the blockNumber
// represented in hex
func (c cloudflareEth) GetBlock(blockNumber string) (*blockchain.Block, error) {
	rpcReq := rpcRequestBody{
		Jsonrpc: c.jsonRpcVersion,
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
