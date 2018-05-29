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
