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
