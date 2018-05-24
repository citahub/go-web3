package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/cryptape/go-web3"
	"github.com/cryptape/go-web3/providers"
)

var hexPrivateKey = `0x3b181467f7f7bc58bccd223efd9b5b8afa1176f5ef9b619a4421f8ac2ce1400b`

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	connection := web3.New(providers.NewHTTPProviders("http://121.196.200.225:1337"))
	status, err := connection.Cita.CreateContract(SimpleStorageBin, hexPrivateKey, strconv.Itoa(r.Int()), 999999, 10)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%-v\n", status)
}
