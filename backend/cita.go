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
	"errors"
	"math/big"

	"github.com/cryptape/go-web3/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	getBlockNumberMethod        = "cita_blockNumber"
	getBlockMetadataMethod      = "cita_getMetaData"
	getTransactionReceiptMethod = "eth_getTransactionReceipt"
)

func (b *backend) GetBlockNumber() (*big.Int, error) {
	resp, err := b.provider.SendRequest(getBlockNumberMethod)
	if err != nil {
		return nil, err
	}

	hexNumber, err := resp.GetString()
	if err != nil {
		return nil, err
	}

	return hexutil.DecodeBig(hexNumber)
}

func (b *backend) GetBlockMetadata(blockNumber *big.Int) (*types.BlockMetadata, error) {
	hexBlockNumber := hexutil.EncodeBig(blockNumber)
	resp, err := b.provider.SendRequest(getBlockMetadataMethod, hexBlockNumber)
	if err != nil {
		return nil, err
	}

	var meta types.BlockMetadata
	if err := resp.GetObject(&meta); err != nil {
		return nil, err
	}

	return &meta, nil
}

func (b *backend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	resp, err := b.provider.SendRequest(getTransactionReceiptMethod, txHash.Hex())
	if err != nil {
		return nil, err
	}

	var r types.Receipt
	if err := resp.GetObject(&r); err != nil {
		return nil, err
	}

	if r.ErrorMessage != "" {
		return nil, errors.New(r.ErrorMessage)
	}

	return &r, nil
}
