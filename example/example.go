package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/cryptape/go-web3"
	"github.com/cryptape/go-web3/providers"
	"github.com/cryptape/go-web3/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var hexPrivateKey = `0x3b181467f7f7bc58bccd223efd9b5b8afa1176f5ef9b619a4421f8ac2ce1400b`

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	connection := web3.New(providers.NewHTTPProviders("http://121.196.200.225:1337"))
	num, _ := connection.Cita.GetBlockNumber()
	status, err := connection.Cita.CreateContract(TokenBin, hexPrivateKey, strconv.Itoa(r.Int()), 999999, num+10, 0)
	if err != nil {
		panic(err)
	}

	parsedAbi, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		panic(err)
	}

	addr := common.BytesToAddress([]byte("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	abiB, err := parsedAbi.Pack("getBalance", addr)
	if err != nil {
		panic(err)
	}

	// meta, _ := connection.Cita.GetBlockMetadata(num)

	// re, err := connection.Cita.CreateContract(abiB, "936bcd7f107a58ebf3c5f1b378f23741b552747bd0630b52ff53842cae23e41e", "testtest", 99999, num+20, meta.ChainID)
	// if err != nil {
	// 	panic(err)
	// }

	privateK, err := crypto.HexToECDSA(utils.CleanHexPrefix(hexPrivateKey))
	if err != nil {
		panic(err)
	}
	fmt.Println(privateK)
	address := crypto.PubkeyToAddress(privateK.PublicKey)

	if err := connection.Cita.Call(address.Hex(), status.ContractAddress, abiB); err != nil {
		panic(err)
	}

	// fmt.Println(re)

	fmt.Println("abiB", string(abiB))
}
