package task

import (
	"context"
	"math/big"
	"time"

	"github.com/GMWalletApp/epusdt/config"
	"github.com/GMWalletApp/epusdt/util/log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type evmLogHandler func(*ethclient.Client, types.Log, int64)

// runEvmWsLogListener connects to wsURL, subscribes to Transfer logs,
// and dispatches each log to handleLog. It retries on transient errors
// with exponential backoff. The ctx lets the caller trigger a clean
// exit — e.g. when admin disables the chain, the caller cancels the
// context and the function returns instead of reconnecting forever.
func runEvmWsLogListener(ctx context.Context, logPrefix, wsURL string, query ethereum.FilterQuery, handleLog evmLogHandler) {
	const (
		minBackoff = 2 * time.Second
		maxBackoff = 60 * time.Second
		rejoinWait = 3 * time.Second
	)
	failWait := minBackoff

	for {
		if ctx.Err() != nil {
			return
		}

		client, err := ethclient.Dial(wsURL)
		if err != nil {
			log.Sugar.Warnf("%s dial: %v, retry in %s", logPrefix, err, failWait)
			if !sleepOrDone(ctx, failWait) {
				return
			}
			failWait = nextBackoff(failWait, maxBackoff)
			continue
		}

		logsCh := make(chan types.Log)
		sub, err := client.SubscribeFilterLogs(ctx, query, logsCh)
		if err != nil {
			client.Close()
			log.Sugar.Warnf("%s subscribe: %v, retry in %s", logPrefix, err, failWait)
			if !sleepOrDone(ctx, failWait) {
				return
			}
			failWait = nextBackoff(failWait, maxBackoff)
			continue
		}
		failWait = minBackoff

		log.Sugar.Infof("%s connected, subscribed to Transfer logs", logPrefix)

		recvLoop(ctx, client, sub, logsCh, logPrefix, handleLog)

		if ctx.Err() != nil {
			return
		}
		if !sleepOrDone(ctx, rejoinWait) {
			return
		}
	}
}

func recvLoop(ctx context.Context, client *ethclient.Client, sub ethereum.Subscription, logsCh <-chan types.Log, logPrefix string, handleLog evmLogHandler) {
	defer func() {
		sub.Unsubscribe()
		client.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Sugar.Infof("%s context cancelled, stopping", logPrefix)
			return
		case err := <-sub.Err():
			if err != nil {
				log.Sugar.Warnf("%s subscription error: %v, reconnecting", logPrefix, err)
			} else {
				log.Sugar.Warnf("%s subscription closed, reconnecting", logPrefix)
			}
			return
		case vLog, ok := <-logsCh:
			if !ok {
				log.Sugar.Warnf("%s log channel closed, reconnecting", logPrefix)
				return
			}
			blockTsMs, ok := waitForConfirmedEVMLog(ctx, client, vLog, logPrefix)
			if !ok {
				continue
			}
			handleLog(client, vLog, blockTsMs)
		}
	}
}

func waitForConfirmedEVMLog(ctx context.Context, client *ethclient.Client, vLog types.Log, logPrefix string) (int64, bool) {
	if vLog.Removed {
		log.Sugar.Warnf("%s dropped removed log tx=%s block=%d", logPrefix, vLog.TxHash.Hex(), vLog.BlockNumber)
		return 0, false
	}

	blockHeader, err := client.HeaderByHash(ctx, vLog.BlockHash)
	if err != nil {
		log.Sugar.Warnf("%s HeaderByHash block=%s: %v", logPrefix, vLog.BlockHash.Hex(), err)
		return 0, false
	}
	if blockHeader == nil || blockHeader.Number == nil {
		log.Sugar.Warnf("%s missing block header for tx=%s block=%d", logPrefix, vLog.TxHash.Hex(), vLog.BlockNumber)
		return 0, false
	}

	confirmations := config.GetEVMConfirmations()
	blockNumber := blockHeader.Number.Uint64()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		latest, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			log.Sugar.Warnf("%s HeaderByNumber latest: %v", logPrefix, err)
		} else if latest != nil && latest.Number != nil && latest.Number.Uint64() >= blockNumber {
			if latest.Number.Uint64()-blockNumber+1 >= confirmations {
				if !isCanonicalEVMBlock(ctx, client, vLog, blockHeader.Number, logPrefix) {
					return 0, false
				}
				return int64(blockHeader.Time) * 1000, true
			}
		}

		select {
		case <-ctx.Done():
			return 0, false
		case <-ticker.C:
		}
	}
}

func isCanonicalEVMBlock(ctx context.Context, client *ethclient.Client, vLog types.Log, blockNumber *big.Int, logPrefix string) bool {
	canonical, err := client.HeaderByNumber(ctx, new(big.Int).Set(blockNumber))
	if err != nil {
		log.Sugar.Warnf("%s canonical block check failed block=%d: %v", logPrefix, vLog.BlockNumber, err)
		return false
	}
	if canonical == nil || canonical.Hash() != vLog.BlockHash {
		log.Sugar.Warnf("%s dropped non-canonical log tx=%s block=%d", logPrefix, vLog.TxHash.Hex(), vLog.BlockNumber)
		return false
	}
	return true
}

// sleepOrDone waits for d or for ctx cancellation, whichever comes
// first. Returns true if the sleep completed normally, false if ctx
// was cancelled (caller should exit).
func sleepOrDone(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func nextBackoff(cur, max time.Duration) time.Duration {
	n := cur * 2
	if n > max {
		return max
	}
	return n
}
