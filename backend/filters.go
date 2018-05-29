package backend

import (
	"context"
	"math/big"
	"sync"

	"github.com/cryptape/go-web3/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	newLogFilterMethod     = "eth_newFilter"
	newBlockFilterMethod   = "eth_newBlockFilter"
	uninstallFilterMethod  = "eth_uninstallFilter"
	getFilterChangesMethod = "eth_getFilterChanges"
	getFilterLogsMethod    = "eth_getFilterLogs"
)

type BlockConsumer func(hashs []common.Hash) (bool, error)
type LogConsumer func(logs []types.Log) (bool, error)

// SubscribeFilterLogs creates a background log filtering operation, returning
// a subscription immediately, which can be used to stream the found events.
func (b *backend) SubscribeLogFilter(ctx context.Context, query ethereum.FilterQuery, consumer LogConsumer) (Subscription, error) {
	id, err := b.newLogFilter(query)
	if err != nil {
		return nil, err
	}

	s := newSub(
		func() (bool, error) {
			logs, err := b.getFilterLogs(id)
			if err != nil {
				return false, err
			}

			if quit, err := consumer(logs); err != nil {
				return false, err
			} else if quit {
				return true, nil
			}

			return false, nil
		},
		func() error {
			_, err := b.uninstallFilter(id)
			return err
		},
	)
	return s, nil
}

func (b *backend) SubscribeBlockFilter(ctx context.Context, consumer BlockConsumer) (Subscription, error) {
	id, err := b.newBlockFilter()
	if err != nil {
		return nil, err
	}

	s := newSub(
		func() (bool, error) {
			hashs, err := b.getFilterChanges(id)
			if err != nil {
				return false, err
			}

			if quit, err := consumer(hashs); err != nil {
				return false, err
			} else if quit {
				return true, nil
			}

			return false, nil
		},
		func() error {
			_, err := b.uninstallFilter(id)
			return err
		},
	)

	return s, nil
}

func (b *backend) newBlockFilter() (string, error) {
	resp, err := b.provider.SendRequest(newBlockFilterMethod)
	if err != nil {
		return "", err
	}

	return resp.GetString()
}

func (b *backend) newLogFilter(query ethereum.FilterQuery) (string, error) {
	resp, err := b.provider.SendRequest(newLogFilterMethod, toFilterArg(query))
	if err != nil {
		return "", err
	}

	return resp.GetString()
}

func (b *backend) uninstallFilter(id string) (bool, error) {
	resp, err := b.provider.SendRequest(uninstallFilterMethod, id)
	if err != nil {
		return false, err
	}

	return resp.GetBool()
}

func (b *backend) getFilterChanges(id string) ([]common.Hash, error) {
	resp, err := b.provider.SendRequest(getFilterChangesMethod, id)
	if err != nil {
		return nil, err
	}

	var hashs []common.Hash
	if err := resp.GetObject(&hashs); err != nil {
		return nil, err
	}

	return hashs, nil
}

func (b *backend) getFilterLogs(id string) ([]types.Log, error) {
	resp, err := b.provider.SendRequest(getFilterLogsMethod, id)
	if err != nil {
		return nil, err
	}

	var logs []types.Log
	if err := resp.GetObject(&logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func toFilterArg(q ethereum.FilterQuery) interface{} {
	arg := map[string]interface{}{
		"fromBlock": toBlockNumArg(q.FromBlock),
		"toBlock":   toBlockNumArg(q.ToBlock),
		"address":   q.Addresses,
		"topics":    q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	}
	return arg
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

type Subscription interface {
	Quit() error
	Unsubscription()
}

func newSub(executor func() (bool, error), clear func() error) Subscription {
	s := &subscription{
		executor: executor,
		clear:    clear,
		quit:     make(chan error),
		unsub:    make(chan struct{}),
	}

	go func(s *subscription) {
		defer s.clear()

		for {
			select {
			case <-s.unsub:
				close(s.unsub)
				return
			default:
				quit, err := s.executor()
				if err != nil {
					s.quit <- err
					return
				} else if quit {
					s.quit <- nil
					return
				}
			}
		}
	}(s)

	return s
}

type subscription struct {
	executor func() (bool, error)
	clear    func() error
	quit     chan error
	unsub    chan struct{}

	unsubscripted bool
	mu            sync.Mutex
}

func (s *subscription) Quit() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.unsubscripted {
		return nil
	}

	return <-s.quit
}

func (s *subscription) Unsubscription() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.unsubscripted {
		return
	}

	s.unsubscripted = true
	s.unsub <- struct{}{}
}
