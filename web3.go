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

package web3

import (
	"github.com/cryptape/go-web3/backend"
	"github.com/cryptape/go-web3/net"
	"github.com/cryptape/go-web3/providers"
)

type Web3 struct {
	Net     net.Interface
	Backend backend.Interface
}

func New(provider providers.Interface) *Web3 {
	return &Web3{
		Net:     net.New(provider),
		Backend: backend.New(provider),
	}
}
