package blockchain

import (
	"math/big"
	"strconv"
)

// Transaction represents a transaction in a block.
type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

func (t Transaction) String() string {
	return strconv.FormatInt(ConvertHexToInt(t.Nonce), 10)
}

func ConvertHexToInt(hex string) int64 {
	bigNumber := new(big.Int)

	bigNumber, ok := bigNumber.SetString(hex, 0)
	if !ok {
		return 0
	}

	return bigNumber.Int64()
}
