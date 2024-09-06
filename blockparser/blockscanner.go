package blockparser

import (
	"context"
	"fmt"
	"time"

	"github.com/spankie/tw-interview/blockchain"
)

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

// getLatestBlockNumber fetches the latest block number from the blockchain.
func (p *Parser) getLatestBlockNumber() (int64, error) {
	blockNumberStr, err := p.blockchainQuerier.GetLatestBlock()
	if err != nil {
		return 0, fmt.Errorf("error fetching latest block: %w", err)
	}

	return blockchain.ConvertHexToInt(blockNumberStr), nil
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
			if len(p.datastore.GetKeys()) > 0 {
				p.saveSubscribedAddressTransactions(p.getTransactionsInBlock(blockNumber))
			}

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
