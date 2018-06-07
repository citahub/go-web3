/*
Copyright 2016-2017 Cryptape Technologies LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package backend

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/cryptape/go-web3/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	getAbiMethod  = "eth_getAbi"
	ethCallMethod = "eth_call"
)

type CallMsg struct {
	From common.Address
	To   *common.Address
	Data []byte
}

func (b *backend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (string, error) {
	blockNumberHex := "latest"
	if blockNumber != nil {
		blockNumberHex = hexutil.EncodeBig(blockNumber)
	}

	resp, err := b.provider.SendRequest(getAbiMethod, contract.Hex(), blockNumberHex)
	if err != nil {
		return "", err
	}
	return resp.GetString()
}

// CodeAt returns the code of the given account. This is needed to differentiate
// between contract internal errors and the local chain being out of sync.
func (b *backend) AbiAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (abi.ABI, error) {
	blockNumberHex := "latest"
	if blockNumber != nil {
		blockNumberHex = hexutil.EncodeBig(blockNumber)
	}

	resp, err := b.provider.SendRequest(getAbiMethod, contract.Hex(), blockNumberHex)
	if err != nil {
		return abi.ABI{}, err
	}
	abiStr, err := resp.GetString()
	if err != nil {
		return abi.ABI{}, err
	}

	return abi.JSON(strings.NewReader(abiStr))
}

// ContractCall executes an Ethereum contract call with the specified data as the
// input.
func (b *backend) CallContract(ctx context.Context, result interface{}, call CallMsg, blockNumber *big.Int) error {
	data := utils.AddHexPrefix(hex.EncodeToString(call.Data))
	params := map[string]string{"from": call.From.Hex(), "data": data, "to": ""}
	if call.To != nil {
		params["to"] = call.To.Hex()
	}

	blockNumberHex := "latest"
	if blockNumber != nil {
		blockNumberHex = hexutil.EncodeBig(blockNumber)
	}

	resp, err := b.provider.SendRequest(ethCallMethod, params, blockNumberHex)
	if err != nil {
		return err
	}

	return resp.GetObject(result)
}
