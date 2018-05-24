package net

import (
	"testing"

	"github.com/cryptape/go-web3/providers"
)

var i Interface

func init() {
	p := providers.NewHTTPProviders("http://127.0.0.1:1337")
	i = New(p)
}

func TestPeerCount(t *testing.T) {
	count, err := i.PeerCount()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cita peer count is %s\n", count)
}
