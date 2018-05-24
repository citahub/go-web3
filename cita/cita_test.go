package cita

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/cryptape/go-web3/providers"
	"github.com/cryptape/go-web3/utils"
)

var i Interface

func init() {
	p := providers.NewHTTPProviders("http://127.0.0.1:1337")
	i = New(p)
}

func TestBlockNumber(t *testing.T) {
	number, err := i.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cita block number is %d\n", number)
}

func TestGetBlockByNumber(t *testing.T) {
	number, err := i.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}

	b, err := i.GetBlockByNumber(number, true)
	if err != nil {
		t.Fatal(err)
	}
	number2, err := utils.ParseUint64(b.Header.Number)
	if err != nil {
		t.Fatal(err)
	}

	if number2 != number {
		t.Fatalf("block number want %d, but got %d\n", number, number2)
	}
}

func TestGetBlockByHash(t *testing.T) {
	number, err := i.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}

	b, err := i.GetBlockByNumber(number, true)
	if err != nil {
		t.Fatal(err)
	}

	prevB, err := i.GetBlockByHash(b.Header.PrevHash, true)
	if err != nil {
		t.Fatal(err)
	}

	if prevB.Hash != b.Header.PrevHash {
		t.Fatalf("prev block hash want %s, but got %s\n", b.Header.PrevHash, prevB.Hash)
	}
}

func TestGetBlockMetadata(t *testing.T) {
	number, err := i.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}

	_, err = i.GetBlockMetadata(number)
	if err != nil {
		t.Fatal(err)
	}
}

var code = []byte(`0x608060405234801561001057600080fd5b5060bf8061001f6000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a72305820e75a7b649045ed321bff861cdbb87d42586007247f49df7d5598ff27206beab30029`)
var hexPrivateKey = `936bcd7f107a58ebf3c5f1b378f23741b552747bd0630b52ff53842cae23e41e`

func TestCreateContract(t *testing.T) {
	status, err := i.CreateContract(code, hexPrivateKey, getRand(), 1000000, 10)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(status.Hash)
	if status.Status != "OK" {
		t.Fatalf("CreateContract status want OK, but got %s\n", status.Status)
	}
}

// TODO: implements eth filters
// func TestGetTransactionProof(t *testing.T) {
// 	status, err := i.CreateContract(code, hexPrivateKey, getRand(), 1000000, 10)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	proof, err := i.GetTransactionProof(status.Hash)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	fmt.Println("proof", proof)
// 	t.Logf("proof hash is %s", proof)
// }

func getRand() string {
	return strconv.Itoa(rand.Intn(10000))
}
