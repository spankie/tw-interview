package blockparser

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spankie/tw-interview/blockparser/blockchain"
)

// Logger is an interface for logging.
type Logger interface {
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}

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

type Parser struct {
	// NOTE: using atomic int64 for thread safety
	lastScannedBlock  atomic.Int64
	datastore         DataStore
	blockchainQuerier BlockchainQuerier
	scanningInterval  time.Duration
	logger            Logger
}

// NewBlockParser creates a new parser and starts the block transactions scanning.
// the returned parser must be used to start scanning for transactions by calling `parser.StartBlockScanning(ctx)`.
func NewBlockParser(cfgOpts ...ConfigOptionResolver) *Parser {
	cfg := Config{}

	for _, opt := range cfgOpts {
		opt(&cfg)
	}

	LoadDefaultConfig(&cfg)

	return newBlockParserWithConfig(cfg)
}

func newBlockParserWithConfig(cfg Config) *Parser {
	parser := &Parser{
		datastore:         cfg.datastore,
		scanningInterval:  cfg.scanningInterval,
		blockchainQuerier: cfg.blockchainQuerier,
		logger:            cfg.logger,
	}

	return parser
}

// last parsed block.
func (p *Parser) GetCurrentBlock() int {
	return int(p.lastScannedBlock.Load())
}

// add address to observer.
func (p *Parser) Subscribe(address string) bool {
	if !isValidEthereumAddress(address) {
		return false
	}

	// check if the address is already subscribed (already in the db).
	if _, ok := p.datastore.Get(address); ok {
		return false
	}

	// if the address is not in the db, add it so it can be observed when
	// scanning the blockchain.
	if err := p.datastore.Add(address, []blockchain.Transaction{}); err != nil {
		return false
	}

	return true
}

// list of inbound or outbound transactions for an address.
func (p *Parser) GetTransactions(address string) []blockchain.Transaction {
	transactions, _ := p.datastore.Get(address)
	return transactions
}

// StartBlockScanning runs a task every minute to find inbound/outbound
// transactions for subscribed address.
func (p *Parser) StartBlockScanning(ctx context.Context) { //nolint: contextcheck
	if ctx == nil {
		ctx = context.Background()
	}

	err := p.initScannedBlockNumber()
	if err != nil {
		p.logger.Error(fmt.Sprintf("could not start block scanning: %v", err))
		return
	}

	if p.scanningInterval == 0 {
		p.scanningInterval = defaultScanningInterval
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				p.logger.Info("block scanning stopped")
				return
			default:
				time.Sleep(p.scanningInterval)
				p.querySubscribedAddressTransactions(ctx)
			}
		}
	}()
}

// getLatestBlockNumber fetches the latest block number from the blockchain.
func (p *Parser) getLatestBlockNumber() (int64, error) {
	blockNumberStr, err := p.blockchainQuerier.GetLatestBlock()
	if err != nil {
		return 0, fmt.Errorf("error fetching latest block: %w", err)
	}

	return blockchain.ConvertHexToInt(blockNumberStr), nil
}

// initScannedBlockNumber initializes the last scanned block with the current block number.
// this is done when the parser is started to ensure that the parser starts scanning
// from the latest block.
func (p *Parser) initScannedBlockNumber() error {
	blockNumber, err := p.getLatestBlockNumber()
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
func (p *Parser) querySubscribedAddressTransactions(ctx context.Context) {
	latestBlockNumber, err := p.getLatestBlockNumber()
	if err != nil {
		p.logger.Error(fmt.Sprintf("error getting latest block: %v", err))
		return
	}

	// start scanning from the last scanned block to the latest block on the blockchain.
	for blockNumber := p.lastScannedBlock.Load() + 1; blockNumber <= latestBlockNumber; blockNumber++ {
		select {
		case <-ctx.Done():
			p.logger.Info(fmt.Sprintf("scanning block %d stopped", blockNumber))
			return
		default:
			p.saveSubscribedAddressTransactions(p.getTransactionsInBlock(blockNumber))
			p.lastScannedBlock.Store(blockNumber)
		}
	}
}

// saveSubscribedAddressTransactions finds and stores all transaction done by subscribed address.
func (p *Parser) saveSubscribedAddressTransactions(blockTransactions []blockchain.Transaction) {
	for _, transaction := range blockTransactions {
		if _, ok := p.datastore.Get(transaction.From); ok {
			err := p.datastore.Add(transaction.From, []blockchain.Transaction{transaction})
			if err != nil {
				p.logger.Error(fmt.Sprintf(
					"error storing transaction %s for address %s %v",
					transaction.String(), transaction.From, err))
			}
		}

		if _, ok := p.datastore.Get(transaction.To); ok {
			err := p.datastore.Add(transaction.To, []blockchain.Transaction{transaction})
			if err != nil {
				p.logger.Error(fmt.Sprintf(
					"error storing transaction %s for address %s %v",
					transaction.String(), transaction.To, err))
			}
		}
	}
}

// getBlock queries the etheruem blockchain to the block identified by the blockNumber
// represented in hex.
func (p *Parser) getBlock(blockNumber string) (*blockchain.Block, error) {
	block, err := p.blockchainQuerier.GetBlock(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("could not get block: %w", err)
	}

	return block, nil
}

// getTransactionsInBlock requires the address and block number.
func (p *Parser) getTransactionsInBlock(blockNumber int64) []blockchain.Transaction {
	block, err := p.getBlock(fmt.Sprintf("0x%x", blockNumber))
	if err != nil {
		return []blockchain.Transaction{}
	}

	return block.Transactions
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
