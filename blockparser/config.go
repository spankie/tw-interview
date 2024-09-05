package blockparser

import (
	"log/slog"
	"time"

	"github.com/spankie/tw-interview/blockparser/blockchain"
	"github.com/spankie/tw-interview/blockparser/cloudflareeth"
)

const (
	defaultScanningInterval = 1 * time.Minute
)

type ConfigOptionResolver func(*Config)

type Config struct {
	datastore         DataStore
	blockchainQuerier BlockchainQuerier
	scanningInterval  time.Duration
	logger            Logger
}

func LoadDefaultConfig(config *Config) {
	if config == nil {
		config = &Config{}
	}

	if config.datastore == nil {
		config.datastore = newMemoryDataStore[blockchain.Transaction]()
	}

	if config.blockchainQuerier == nil {
		config.blockchainQuerier = cloudflareeth.NewClient()
	}

	if config.scanningInterval == 0 {
		config.scanningInterval = defaultScanningInterval
	}

	if config.logger == nil {
		config.logger = slog.Default()
	}
}

func WithLogger(logger Logger) ConfigOptionResolver {
	return func(c *Config) {
		c.logger = logger
	}
}

func WithDataStore(datastore DataStore) ConfigOptionResolver {
	return func(c *Config) {
		c.datastore = datastore
	}
}

func WithBlockchainQuerier(blockchainQuerier BlockchainQuerier) ConfigOptionResolver {
	return func(c *Config) {
		c.blockchainQuerier = blockchainQuerier
	}
}

func WithScanningInterval(scanningInterval time.Duration) ConfigOptionResolver {
	return func(c *Config) {
		c.scanningInterval = scanningInterval
	}
}
