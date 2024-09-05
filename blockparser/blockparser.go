package blockparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spankie/tw-interview/blockparser/blockchain"
	"github.com/spankie/tw-interview/blockparser/cloudflareeth"
)

const (
	defaultScanningInterval = 1 * time.Minute
)

// DataStore is an interface for storing data and querying data.
type DataStore interface {
	Add(key string, value []blockchain.Transaction) error
	Get(key string) ([]blockchain.Transaction, bool)
}

// BlockchainQuerier is an interface for querying the blockchain.
type BlockchainQuerier interface {
	GetLatestBlock() (string, error)
	GetBlock(blockNumber string) (*blockchain.Block, error)
}

type BlockParser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []blockchain.Transaction
}

type parser struct {
	// NOTE: using atomic int64 for thread safety
	lastScannedBlock  atomic.Int64
	datastore         DataStore
	blockchainQuerier BlockchainQuerier
	scanningInterval  time.Duration
}

// NewBlockParser creates a new parser and starts the block transactions scanning.
// TODO: probably need to use function options to set the neccessary dependencies
// TODO: add logger for better logging control.
func NewBlockParser(ctx context.Context, datastore DataStore,
	blockchainQuerier BlockchainQuerier, scanningInterval time.Duration) BlockParser {

	newParser := &parser{
		datastore:         datastore,
		scanningInterval:  scanningInterval,
		blockchainQuerier: blockchainQuerier,
	}

	newParser.setDefaultOptions()

	go newParser.startBlockScanning(ctx)

	return newParser
}

// last parsed block.
func (p *parser) GetCurrentBlock() int {
	return int(p.lastScannedBlock.Load())
}

// add address to observer.
func (p *parser) Subscribe(address string) bool {
	if !isValidEthereumAddress(address) {
		return false
	}

	// check if the address is already subscribed (already in the db)
	if _, ok := p.datastore.Get(address); ok {
		return false
	}

	// if the address is not in the db, add it so it can be observed when
	// scanning the blockchain
	if err := p.datastore.Add(address, []blockchain.Transaction{}); err != nil {
		return false
	}

	return true
}

// list of inbound or outbound transactions for an address.
func (p *parser) GetTransactions(address string) []blockchain.Transaction {
	transactions, _ := p.datastore.Get(address)
	return transactions
}

func (p *parser) setDefaultOptions() {
	if p.datastore == nil {
		p.datastore = newMemoryDataStore[blockchain.Transaction]()
	}

	if p.blockchainQuerier == nil {
		p.blockchainQuerier = cloudflareeth.NewCloudflareEthClient()
	}

	if p.scanningInterval == 0 {
		p.scanningInterval = defaultScanningInterval
	}
}

func (p *parser) getLatestBlock() (int64, error) {
	blockNumberStr, err := p.blockchainQuerier.GetLatestBlock()
	if err != nil {
		return 0, fmt.Errorf("error fetching latest block: %w", err)
	}

	return blockchain.ConvertHexToInt(blockNumberStr), nil
}

// start scanning from the latest block.
func (p *parser) initScannedBlock() error {
	blockNumber, err := p.getLatestBlock()
	if err != nil {
		return err
	}

	p.lastScannedBlock.Store(blockNumber)

	return nil
}

// querySubscribedAddressTransactions scans the blockchain for transactions
// it starts from the last scanned blocked to the current block
// for each block, it filters out transactions done by subscribed addresses
// and saves them in the datastore.
func (p *parser) querySubscribedAddressTransactions(ctx context.Context) {
	currentBlock, err := p.getLatestBlock()
	if err != nil {
		slog.Error(fmt.Sprintf("error getting latest block: %v", err))
		return
	}

	for blockNumber := p.lastScannedBlock.Load() + 1; blockNumber <= currentBlock; blockNumber++ {
		select {
		case <-ctx.Done():
			slog.Info(fmt.Sprintf("scanning block %d stopped", blockNumber))
			return
		default:
			p.saveSubscribedAddressTransactions(p.getTransactionsInBlock(blockNumber))
			p.lastScannedBlock.Store(blockNumber)
		}
	}
}

// saveSubscribedAddressTransactions finds and stores all transaction done by subscribed address.
func (p *parser) saveSubscribedAddressTransactions(transactions []blockchain.Transaction) {
	for _, transaction := range transactions {
		if addr, ok := p.datastore.Get(transaction.From); ok {
			err := p.datastore.Add(transaction.From, append(addr, transaction))
			if err != nil {
				// if this happens we probably want to do some more than logging it.
				slog.Error(fmt.Sprintf("error storing transaction %s for address %s %v", transaction.String(), addr, err))
			}
		}

		if addr, ok := p.datastore.Get(transaction.To); ok {
			err := p.datastore.Add(transaction.To, append(addr, transaction))
			if err != nil {
				// if this happens we probably want to do some more than logging it.
				slog.Error(fmt.Sprintf("error storing transaction %s for address %s %v", transaction.String(), addr, err))
			}
		}
	}
}

// getBlock queries the etheruem blockchain to the block identified by the blockNumber
// represented in hex.
func (p *parser) getBlock(blockNumber string) (*blockchain.Block, error) {
	block, err := p.blockchainQuerier.GetBlock(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("could not get block: %w", err)
	}

	return block, nil
}

// getTransactionsInBlock requires the address and block number.
func (p *parser) getTransactionsInBlock(blockNumber int64) []blockchain.Transaction {
	block, err := p.getBlock(fmt.Sprintf("0x%x", blockNumber))
	if err != nil {
		return []blockchain.Transaction{}
	}

	return block.Transactions
}

// startBlockScanning runs a task every minute to find inbound/outbound
// transactions for subscribed address.
func (p *parser) startBlockScanning(ctx context.Context) { //nolint: contextcheck
	if ctx == nil {
		ctx = context.Background()
	}

	err := p.initScannedBlock()
	if err != nil {
		slog.Error(fmt.Sprintf("could not start block scanning: %v", err))
		return
	}

	if p.scanningInterval == 0 {
		p.scanningInterval = defaultScanningInterval
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("block scanning stopped")
			return
		default:
			time.Sleep(p.scanningInterval)
			p.querySubscribedAddressTransactions(ctx)
		}
	}
}

// isValidEthereumAddress validates an ethereum address.
func isValidEthereumAddress(address string) bool {
	// Check if the address starts with "0x".
	if !strings.HasPrefix(address, "0x") {
		return false
	}

	// Remove the "0x" prefix.
	address = strings.TrimPrefix(address, "0x")

	// Check if the address has exactly 40 characters.
	if len(address) != 40 {
		return false
	}

	// Check if the address contains only hexadecimal characters.
	for _, char := range address {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}

	return true
}
