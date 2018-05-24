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
