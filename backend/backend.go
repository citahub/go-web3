// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	"context"
	"errors"
	"math/big"

	"github.com/cryptape/go-web3/providers"
	"github.com/cryptape/go-web3/types"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// ErrNoCode is returned by call and transact operations for which the requested
	// recipient contract to operate on does not exist in the state db or does not
	// have any code associated with it (i.e. suicided).
	ErrNoCode = errors.New("no contract code at given address")

	// This error is raised when attempting to perform a pending state action
	// on a backend that doesn't implement PendingContractCaller.
	ErrNoPendingState = errors.New("backend does not support pending state")

	// This error is returned by WaitDeployed if contract creation leaves an
	// empty contract behind.
	ErrNoCodeAfterDeploy = errors.New("no contract code after deployment")
)

// BoundContract is the base wrapper object that reflects a contract on the
// Ethereum network. It contains a collection of methods that are used by the
// higher level contract bindings to operate.
type BoundContract struct {
	Address common.Address // Deployment address of the contract on the Ethereum blockchain
	Abi     abi.ABI        // Reflect based ABI to access the correct Ethereum methods
	// caller     ContractCaller     // Read interface to interact with the blockchain
	// transactor ContractTransactor // Write interface to interact with the blockchain
	// filterer   ContractFilterer   // Event filtering to interact with the blockchain
}

// ContractCaller defines the methods needed to allow operating with contract on a read
// only basis.
type ContractCaller interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (string, error)
	AbiAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (abi.ABI, error)
	// ContractCall executes an Ethereum contract call with the specified data as the
	// input.
	CallContract(ctx context.Context, result interface{}, call CallMsg, blockNumber *big.Int) error
}

// ContractTransactor defines the methods needed to allow operating with contract
// on a write only basis. Beside the transacting method, the remainder are helpers
// used when the user does not provide some needed values, but rather leaves it up
// to the transactor to decide.
type ContractTransactor interface {
	// SendTransaction injects the transaction into the pending pool for execution.
	SendTransaction(ctx context.Context, tx *types.Transaction, hexPrivateKey string) (common.Hash, error)
	// DeployContract deploys a contract onto the Ethereum blockchain and binds the
	// deployment address with a Go wrapper.
	DeployContract(ctx context.Context, params *types.TransactParams, abi, code string) (common.Hash, *BoundContract, error)
}

// ContractFilterer defines the methods needed to access log events using one-off
// queries or continuous event subscriptions.
type ContractFilterer interface {
	// SubscribeFilterLogs creates a background log filtering operation, returning
	// a subscription immediately, which can be used to stream the found events.
	SubscribeLogFilter(ctx context.Context, query ethereum.FilterQuery, consumer LogConsumer) (Subscription, error)
	SubscribeBlockFilter(ctx context.Context, consumer BlockConsumer) (Subscription, error)
}

// See https://cryptape.github.io/cita/zh/usage-guide/rpc/index.html
type Cita interface {
	GetBlockNumber() (*big.Int, error)
	GetBlockMetadata(blockNumber *big.Int) (*types.BlockMetadata, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

// Interface defines the methods needed to work with contracts on a read-write basis.
type Interface interface {
	ContractCaller
	ContractTransactor
	ContractFilterer
	Cita
}

func New(provider providers.Interface) Interface {
	return &backend{
		provider: provider,
	}
}

type backend struct {
	provider providers.Interface
}
