package genesis

import (
	"context"
	"genesis/blockchain"
	"genesis/config"
	"log"
	"time"

	"github.com/ninja-software/terror"
)

// SystemTicker contains info on the system's ticker (for committing to blockchain automatically)
type SystemTicker struct {
	LastTick         time.Time
	TickInterval     time.Duration
	transactionStore TransactionStorer
	manifestStore    ManifestStorer
	blk              *blockchain.Service
	ticker           *time.Ticker
	stopTicker       chan bool
}

// NewTicker creates and starts a new SystemTicker (for committing to blockchain automatically)
func NewTicker(
	ctx context.Context,
	config *config.PlatformConfig,
	transactionStore TransactionStorer,
	manifestStore ManifestStorer,
	blk *blockchain.Service,
) *SystemTicker {
	interval := time.Duration(config.Blockchain.FlushPendingActionsInterval) * time.Hour
	SystemTicker := &SystemTicker{
		LastTick:         time.Now().UTC(),
		TickInterval:     interval,
		transactionStore: transactionStore,
		manifestStore:    manifestStore,
		blk:              blk,
	}
	SystemTicker.Start(ctx)

	return SystemTicker
}

func (s *SystemTicker) tick(
	ctx context.Context,
	t time.Time,
) {
	log.Printf("[TICKER] Committing... pending transactions to blockchain (%v)\n", t)

	s.LastTick = t.UTC()

	success, err := FlushPendingTransactions(ctx, s.transactionStore, s.manifestStore, s.blk)
	if err != nil {
		log.Println("[TICKER] Failed to commit to blockchain")
		terror.Echo(err)
		return
	}

	if !success {
		log.Println("[TICKER] Failed to commit to blockchain")
		return
	}

	log.Println("[TICKER] Committed pending transactions to blockchain")

	// Check balance
	err = s.blk.CheckBalance(ctx)
	if err != nil {
		log.Println("[TICKER] Failed to check current ETH account balance")
		terror.Echo(err)
	}
}

// Start starts are new ticker
func (s *SystemTicker) Start(ctx context.Context) {
	ticker := time.NewTicker(s.TickInterval)
	s.ticker = ticker
	s.stopTicker = make(chan bool)

	go func() {
		for {
			select {
			case <-s.stopTicker:
				return
			case t := <-s.ticker.C:
				s.tick(ctx, t)
			}
		}
	}()
}

// Reset stops current ticker and starts a new one + updates LastTick to now
func (s *SystemTicker) Reset(ctx context.Context) {
	s.LastTick = time.Now().UTC()
	s.stopTicker <- true
	s.Start(ctx)
}
