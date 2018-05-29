package backend

import (
	"context"
	"fmt"
	"testing"

	"github.com/cryptape/go-web3/providers"
	"github.com/cryptape/go-web3/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

var be = New(providers.NewHTTPProviders("http://121.196.200.225:1337"))

func TestBlockFilter(t *testing.T) {
	sub, err := be.SubscribeBlockFilter(context.Background(), func(hashs []common.Hash) (bool, error) {
		for _, hash := range hashs {
			fmt.Println(hash.Hex())
		}

		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := sub.Quit(); err != nil {
		t.Fatal(err)
	}

	t.Log("quit")
}

func TestLogsFilter(t *testing.T) {
	query := ethereum.FilterQuery{}
	sub, err := be.SubscribeLogFilter(context.Background(), query, func(logs []types.Log) (bool, error) {
		t.Log(logs)
		for _, log := range logs {
			fmt.Println(log)
		}

		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := sub.Quit(); err != nil {
		t.Fatal(err)
	}

	t.Log("quit")
}
