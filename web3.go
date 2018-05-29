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
