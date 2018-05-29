package backend

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/cryptape/go-web3/providers"
	"github.com/cryptape/go-web3/types"
	"github.com/cryptape/go-web3/utils"

	"github.com/ethereum/go-ethereum/common"
)

var callerBE = New(providers.NewHTTPProviders("http://47.75.129.215:1337"))

func TestCodeAt(t *testing.T) {
	_, bc, err := deployContract(t)
	if err != nil {
		t.Fatal(err)
	}

	abi, err := callerBE.AbiAt(context.Background(), bc.Address, nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(abi)
}

func TestCaller(t *testing.T) {
	_, bc, err := deployContract(t)
	if err != nil {
		t.Fatal(err)
	}

	from, err := utils.PrivateKeyToAddress(hexPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	abiB, err := bc.Abi.Pack("getTest")
	if err != nil {
		t.Fatal(err)
	}

	msg := CallMsg{
		From: from,
		To:   &bc.Address,
		Data: abiB,
	}

	result := new(string)
	if err = callerBE.CallContract(context.Background(), result, msg, nil); err != nil {
		t.Fatal(err)
	}

	t.Logf("get test result is %s", *result)
}

func deployContract(t *testing.T) (common.Hash, *BoundContract, error) {
	blockNumber, err := callerBE.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	meta, err := callerBE.GetBlockMetadata(blockNumber)
	if err != nil {
		t.Fatal(err)
	}

	params := &types.TransactParams{
		HexPrivateKey:   hexPrivateKey,
		ValidUntilBlock: blockNumber.Add(blockNumber, big.NewInt(88)),
		Nonce:           "test",
		ChainID:         meta.ChainID,
		Version:         0,
		Quota:           big.NewInt(10000000),
	}

	txHash, bc, err := callerBE.DeployContract(context.Background(), params, tokenABI, tokenBin)
	if err != nil {
		t.Fatal(err)
	}

	return txHash, bc, nil
}
