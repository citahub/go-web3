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

package net

import (
	"github.com/cryptape/go-web3/providers"
)

const (
	peerCountMethod = "net_peerCount"
)

type Interface interface {
	PeerCount() (string, error)
}

func New(provider providers.Interface) Interface {
	return &net{
		provider: provider,
	}
}

type net struct {
	provider providers.Interface
}

func (n *net) PeerCount() (string, error) {
	resp, err := n.provider.SendRequest(peerCountMethod)
	if err != nil {
		return "", err
	}

	return resp.GetString()
}
