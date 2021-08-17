package chain

import (
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/crypto-com/chain-indexing/internal/primptr"

	"github.com/crypto-com/chain-indexing/appinterface/tendermint"
	applogger "github.com/crypto-com/chain-indexing/internal/logger"
)

const DEFAULT_POLLING_INTERVAL = 5 * time.Second

type BlockHeightTracker struct {
	logger applogger.Logger
	client tendermint.Client

	pollingInterval time.Duration

	subscriptions []chan<- int64

	latestBlockHeight *int64
	rwMutex           sync.RWMutex
}

func NewBlockHeightTracker(logger applogger.Logger, client tendermint.Client) *BlockHeightTracker {
	tracker := &BlockHeightTracker{
		logger: logger.WithFields(applogger.LogFields{
			"module": "BlockHeightTracker",
		}),
		client: client,

		pollingInterval: DEFAULT_POLLING_INTERVAL,

		subscriptions: make([]chan<- int64, 0),

		latestBlockHeight: primptr.Int64Nil(),
	}

	go tracker.Run()

	return tracker
}

func (tracker *BlockHeightTracker) Run() {
	for {
		operation := func() error {
			height, err := tracker.client.LatestBlockHeight()
			if err != nil {
				return fmt.Errorf("error getting chain latest block height: %v", err)
			}

			for _, subscription := range tracker.subscriptions {
				select {
				case subscription <- height:
				default:
					tracker.logger.Info("block subscription channel is blocked, maybe busy?")
				}
			}

			tracker.rwMutex.Lock()
			tracker.latestBlockHeight = &height
			tracker.rwMutex.Unlock()

			tracker.logger.Infof("updated chain latest block height: %d", height)
			<-time.After(tracker.pollingInterval)
			return nil
		}
		notifyFn := func(opErr error, backoffDuration time.Duration) {
			tracker.logger.Errorf("retrying in %s: %v", backoffDuration.String(), opErr)
		}
		_ = backoff.RetryNotify(operation, backoff.NewExponentialBackOff(), notifyFn)
	}
}

func (tracker *BlockHeightTracker) Subscribe(ch chan<- int64) {
	tracker.subscriptions = append(tracker.subscriptions, ch)
}

func (tracker *BlockHeightTracker) GetLatestBlockHeight() *int64 {
	tracker.rwMutex.RLock()
	defer tracker.rwMutex.RUnlock()

	return tracker.latestBlockHeight
}
