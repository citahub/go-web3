/*
Copyright 2016-2017 Cryptape Technologies LLC.

This program is free software: you can redistribute it
and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation,
either version 3 of the License, or (at your option) any
later version.

This program is distributed in the hope that it will be
useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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
