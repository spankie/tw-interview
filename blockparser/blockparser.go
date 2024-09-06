package blockparser

import (
	"sync/atomic"
	"time"

	"github.com/spankie/tw-interview/blockchain"
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
	GetKeys() []string
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
	if !blockchain.IsValidEthereumAddress(address) {
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
