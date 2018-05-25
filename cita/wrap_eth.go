package cita

import (
	"github.com/cryptape/go-web3/providers"
	"github.com/cryptape/go-web3/types"
)

const (
	getTransactionReceiptMethod = "eth_getTransactionReceipt"
)

type ethInterface interface {
	GetTransactionReceipt(hash string) (*types.Receipt, error)

	NewBlockFilter() (uint64, error)
	UninstallFilter(num uint64) (bool, error)
	GetFilterChanges(id uint64) ([]string, error)
}

func newEth(provider providers.Interface) ethInterface {
	return &eth{
		provider: provider,
	}
}

type eth struct {
	provider providers.Interface
}

func (e *eth) GetTransactionReceipt(hash string) (*types.Receipt, error) {
	resp, err := e.provider.SendRequest(getTransactionReceiptMethod, hash)
	if err != nil {
		return nil, err
	}

	var r types.Receipt
	if err := resp.GetObject(&r); err != nil {
		return nil, err
	}

	return &r, nil
}
