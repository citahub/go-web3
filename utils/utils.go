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

package utils

import (
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func ParseHexToUint64(s string) (uint64, error) {
	num, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return 0, err
	}

	return uint64(num), nil
}

func ConvUint64ToHex(num uint64) string {
	return AddHexPrefix(strconv.FormatUint(num, 16))
}

func CleanHexPrefix(s string) string {
	if strings.HasPrefix(s, "0x") {
		return strings.TrimPrefix(s, "0x")
	}

	return s
}

func AddHexPrefix(s string) string {
	if !strings.HasPrefix(s, "0x") {
		return "0x" + s
	}

	return s
}

func PrivateKeyToAddress(hexPrivateKey string) (common.Address, error) {
	privateK, err := crypto.HexToECDSA(CleanHexPrefix(hexPrivateKey))
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(privateK.PublicKey), nil
}
