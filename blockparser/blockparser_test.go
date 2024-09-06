package blockparser

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/spankie/tw-interview/blockchain"
)

type MockBlockchainQuerier struct {
	LatestBlock    int
	LatestBlockErr error
	Block          *blockchain.Block
	BlockErr       error
}

func (m *MockBlockchainQuerier) GetLatestBlock() (string, error) {
	if m.LatestBlockErr != nil {
		return "", m.LatestBlockErr
	}

	m.LatestBlock++
	latestBlock := m.LatestBlock

	return strconv.Itoa(latestBlock), nil
}

func (m *MockBlockchainQuerier) GetBlock(_ string) (*blockchain.Block, error) {
	if m.BlockErr != nil {
		return nil, m.BlockErr
	}

	return m.Block, nil
}

var sampleBlock = blockchain.Block{
	Difficulty: "0x4ea3f27bc",
	ExtraData:  "0x476574682f4c5649562f76312e302e302f6c696e75782f676f312e342e32",
	GasLimit:   "0x1388",
	GasUsed:    "0x0",
	Hash:       "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae",
	LogsBloom: "0x0000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
		"00000000000000000000000",
	Miner:           "0xbb7b8287f3f0a933474a79eae42cbca977791171",
	MixHash:         "0x4fffe9ae21f1c9e15207b1f472d5bbdd68c9595d461666602f2be20daf5e7843",
	Nonce:           "0x689056015818adbe",
	Number:          "0x1b4",
	ParentHash:      "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
	ReceiptsRoot:    "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
	Sha3Uncles:      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
	Size:            "0x220",
	StateRoot:       "0xddc8b0234c2e0cad087c8b389aa7ef01f7d79b2570bccb77ce48648aa61c904d",
	Timestamp:       "0x55ba467c",
	TotalDifficulty: "0x78ed983323d",
	Transactions: []blockchain.Transaction{
		{
			BlockHash:        "0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
			BlockNumber:      "0x5daf3b",
			From:             "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
			Gas:              "0xc350",
			GasPrice:         "0x4a817c800",
			Hash:             "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b",
			Input:            "0x68656c6c6f21",
			Nonce:            "0x15",
			To:               "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
			TransactionIndex: "0x41",
			Value:            "0xf3dbb76162000",
			V:                "0x25",
			R:                "0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fea",
			S:                "0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721c",
		},
		{
			BlockHash:        "0x2d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
			BlockNumber:      "0x5daf3c",
			From:             "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
			Gas:              "0xc351",
			GasPrice:         "0x4a817c801",
			Hash:             "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944c",
			Input:            "0x68656c6c6f22",
			Nonce:            "0x16",
			To:               "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bc",
			TransactionIndex: "0x42",
			Value:            "0xf3dbb76162001",
			V:                "0x26",
			R:                "0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5feb",
			S:                "0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721d",
		},
		{
			BlockHash:        "0x3d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
			BlockNumber:      "0x5daf3d",
			From:             "0xa7d9ddbe1f17865597fbd27ec712455208b6b76f",
			Gas:              "0xc352",
			GasPrice:         "0x4a817c802",
			Hash:             "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944d",
			Input:            "0x68656c6c6f23",
			Nonce:            "0x17",
			To:               "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bd",
			TransactionIndex: "0x43",
			Value:            "0xf3dbb76162002",
			V:                "0x27",
			R:                "0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fec",
			S:                "0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721e",
		},
	},
	TransactionsRoot: "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
	Uncles:           []string{},
}

func TestParserSubscription(t *testing.T) {
	datastore := newMemoryDataStore[blockchain.Transaction]()
	block := sampleBlock
	blockchainQuerier := &MockBlockchainQuerier{
		LatestBlock: 0x7b,
		Block:       &block,
	}

	parser := NewBlockParser(WithDataStore(datastore),
		WithBlockchainQuerier(blockchainQuerier), WithScanningInterval(1*time.Second))

	// invalid address
	invalidAddr := "0xabc"
	if ok := parser.Subscribe(invalidAddr); ok {
		t.Errorf("should not subscribe an invalid address (%s). got %v, want false", invalidAddr, ok)
	}

	addr1 := "0x26bce6ecb5b10138e4bf14ac0ffcc8727fef3b2e"
	if ok := parser.Subscribe(addr1); !ok {
		t.Errorf("should subscribe a valid address (%s). got %v, want true", addr1, ok)
	}

	if ok := parser.Subscribe(addr1); ok {
		t.Errorf("should not subscribe an already subscribed address. got %v, want false", ok)
	}

	transactions := parser.GetTransactions(addr1)
	if len(transactions) != 0 {
		t.Errorf("should not have any transactions for address %s. got %d, want 0", addr1, len(transactions))
	}
}

func TestParserGetTransactions(t *testing.T) {
	t.Run("test parser get transactions", func(t *testing.T) {
		datastore := newMemoryDataStore[blockchain.Transaction]()
		block := sampleBlock

		blockchainQuerier := &MockBlockchainQuerier{
			LatestBlock: 0x7b,
			Block:       &block,
		}

		parser := NewBlockParser(WithDataStore(datastore),
			WithBlockchainQuerier(blockchainQuerier), WithScanningInterval(1*time.Second))

		address := blockchainQuerier.Block.Transactions[0].From

		if subscribed := parser.Subscribe(address); !subscribed {
			t.Errorf("should subscribe address %s; got %v, want true", address, subscribed)
		}

		// initialize the scanned block.
		err := parser.initScannedBlockNumber()
		if err != nil {
			t.Errorf("initScannedBlockNumber() = %v, want nil error", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// query the block to find transaction for subscribed address.
		parser.querySubscribedAddressTransactions(ctx)

		// the transactions should be two.
		if transactions := parser.GetTransactions(address); len(transactions) != 2 {
			t.Errorf("should get 2 transactions but got %d", len(transactions))
		}

		// try to query the same block again and make sure the transactions are not duplicated.
		blockchainQuerier.LatestBlock--

		parser.querySubscribedAddressTransactions(ctx)

		if transactions := parser.GetTransactions(address); len(transactions) != 2 {
			t.Errorf("should get 2 transactions, but got: %v", len(transactions))
		}
	})
}

func TestStartBlockScanning(t *testing.T) {
	t.Run("test start block scanning", func(t *testing.T) {
		datastore := newMemoryDataStore[blockchain.Transaction]()
		block := sampleBlock

		blockchainQuerier := &MockBlockchainQuerier{
			LatestBlock: 0x7b,
			Block:       &block,
		}

		parser := NewBlockParser(WithDataStore(datastore),
			WithBlockchainQuerier(blockchainQuerier), WithScanningInterval(500*time.Millisecond))

		address := blockchainQuerier.Block.Transactions[0].From

		if subscribed := parser.Subscribe(address); !subscribed {
			t.Errorf("should subscribe address %s; got %v, want true", address, subscribed)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// start scanning the block.
		parser.StartBlockScanning(ctx)

		time.Sleep(1100 * time.Millisecond)

		// the transactions should be 4.
		expectedNumTransactions := 4

		if transactions := parser.GetTransactions(address); len(transactions) != expectedNumTransactions {
			t.Errorf("should get %d transactions but got %d", expectedNumTransactions, len(transactions))
		}
	})
}
