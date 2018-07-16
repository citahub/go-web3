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
