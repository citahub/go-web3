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
