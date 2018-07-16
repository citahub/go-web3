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
