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
