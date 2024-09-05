package blockparser

import (
	"testing"

	"github.com/spankie/tw-interview/blockparser/blockchain"
)

func TestNewDataStore(t *testing.T) {
	store := newMemoryDataStore[blockchain.Transaction]()
	if store == nil {
		t.Errorf("NewDataStore() = nil, want non-nil")
	}
}

var transactions = []blockchain.Transaction{
	{
		BlockHash:        "0xabd",
		BlockNumber:      "0x2",
		From:             "0xdf0",
		Gas:              "0x3",
		GasPrice:         "0x4",
		Hash:             "1xghi",
		Input:            "0x5",
		Nonce:            "0x6",
		To:               "1xjkl",
		TransactionIndex: "0x7",
		Value:            "0x8",
		V:                "0x9",
		R:                "0xa",
		S:                "0xb",
	},
	{
		BlockHash:        "0xbcd",
		BlockNumber:      "0x2",
		From:             "0xefg",
		Gas:              "0x3",
		GasPrice:         "0x4",
		Hash:             "0xhij",
		Input:            "0x5",
		Nonce:            "0x6",
		To:               "0xklm",
		TransactionIndex: "0x7",
		Value:            "0x8",
		V:                "0x9",
		R:                "0xa",
		S:                "0xb",
	},
	{
		BlockHash:        "0xcde",
		BlockNumber:      "0x3",
		From:             "0xfgh",
		Gas:              "0x4",
		GasPrice:         "0x5",
		Hash:             "0ijk",
		Input:            "0x6",
		Nonce:            "0x7",
		To:               "0xlmn",
		TransactionIndex: "0x8",
		Value:            "0x9",
		V:                "0xa",
		R:                "0xb",
		S:                "0xc",
	},
}

func TestAdd(t *testing.T) {
	store := newMemoryDataStore[blockchain.Transaction]()

	t.Run("test add", func(t *testing.T) {
		err := store.Add("testKey", transactions[:1])
		if err != nil {
			t.Errorf("expected nil err, got: %v", err)
		}

		gotTransactions, ok := store.Get("testKey")
		if !ok {
			t.Errorf("Get() = _, %v, want _, true", ok)
		}

		if gotTransactions == nil {
			t.Errorf("Get() = nil, want non-nil")
		}

		expectedNumTransactions := 1

		if len(gotTransactions) != expectedNumTransactions {
			t.Errorf("len(Get()) = %v, want %d", len(gotTransactions), expectedNumTransactions)
		}

		if gotTransactions[0].From != transactions[0].From {
			t.Errorf("From = %v, want %v", gotTransactions[0].From, transactions[0].From)
		}
	})
}

func TestAppend(t *testing.T) {
	store := newMemoryDataStore[blockchain.Transaction]()

	t.Run("test append", func(t *testing.T) {
		err := store.Add("testKey", transactions[0:1])
		if err != nil {
			t.Errorf("expected nil err, got: %v", err)
		}

		err = store.Add("testKey", transactions[1:2])
		if err != nil {
			t.Errorf("expected nil err, got: %v", err)
		}

		gotTransactions, ok := store.Get("testKey")
		if !ok {
			t.Errorf("Get() = _, %v, want _, true", ok)
		}

		if gotTransactions == nil {
			t.Errorf("Get() = nil, want non-nil")
		}

		expectedNumTransactions := 2
		if len(gotTransactions) != expectedNumTransactions {
			t.Errorf("len(Get()) = %v, want %d", len(gotTransactions), expectedNumTransactions)
		}

		if gotTransactions[1].From != transactions[1].From {
			t.Errorf("From = %v, want %v", gotTransactions[1].From, transactions[1].From)
		}
	})
}

func TestAddEmptyKey(t *testing.T) {
	store := newMemoryDataStore[blockchain.Transaction]()

	t.Run("test empty key", func(t *testing.T) {
		if err := store.Add("", transactions); err == nil {
			t.Errorf("should have returned an error for empty key")
		}

		if err := store.Add("  ", transactions); err == nil {
			t.Errorf("should have returned an error for empty key")
		}
	})
}
