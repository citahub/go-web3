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
