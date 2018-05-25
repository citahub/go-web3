package providers

import (
	"fmt"

	"github.com/cryptape/go-web3/errors"
	"github.com/ybbus/jsonrpc"
)

// Interface is providers interfaces.
type Interface interface {
	SendRequest(method string, params ...interface{}) (*jsonrpc.RPCResponse, error)
	Close() error
}

// NewHTTPProviders returns an instance of http provider.
func NewHTTPProviders(address string) Interface {
	return &httpProvider{
		address: address,
		client:  jsonrpc.NewClient(address),
	}
}

type httpProvider struct {
	address string
	client  jsonrpc.RPCClient
}

func (hp *httpProvider) SendRequest(method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	resp, err := hp.client.Call(method, params)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Code, resp.Error.Message)
	}

	if resp.Result == nil {
		return nil, errors.New(0, fmt.Sprintf("the result of method %s is null", method))
	}

	return resp, nil
}

func (hp *httpProvider) Close() error {
	return nil
}
