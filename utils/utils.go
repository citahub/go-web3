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
