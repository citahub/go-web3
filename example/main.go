package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/cryptape/go-web3"
	"github.com/cryptape/go-web3/backend"
	"github.com/cryptape/go-web3/providers"
	"github.com/cryptape/go-web3/types"
	"github.com/cryptape/go-web3/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var hexPrivateKey = `0x3b181467f7f7bc58bccd223efd9b5b8afa1176f5ef9b619a4421f8ac2ce1400b`

// TokenABI is the input ABI used to generate the binding from.
const TokenABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"supplyAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// TokenBin is the compiled bytecode used for deploying new contracts.
const TokenBin = `0x608060405234801561001057600080fd5b5060405160208061025e83398101604090815290513360009081526020819052919091205561021a806100446000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166327e235e3811461005b578063a9059cbb1461009b578063f8b2cb4f146100e0575b600080fd5b34801561006757600080fd5b5061008973ffffffffffffffffffffffffffffffffffffffff6004351661010e565b60408051918252519081900360200190f35b3480156100a757600080fd5b506100cc73ffffffffffffffffffffffffffffffffffffffff60043516602435610120565b604080519115158252519081900360200190f35b3480156100ec57600080fd5b5061008973ffffffffffffffffffffffffffffffffffffffff600435166101c6565b60006020819052908152604090205481565b33600090815260208190526040812054821180159061013f5750600082115b156101bc57336000818152602081815260408083208054879003905573ffffffffffffffffffffffffffffffffffffffff871680845292819020805487019055805186815290519293927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a35060016101c0565b5060005b92915050565b73ffffffffffffffffffffffffffffffffffffffff16600090815260208190526040902054905600a165627a7a723058206ccefef77f8d04f19aa699995bf84fbe7eaef2a3d57eccbdd06cc570f50b6b730029`

func main() {
	client := web3.New(providers.NewHTTPProviders("http://47.75.129.215:1337"))
	deployContract(client)
}

func deployContract(client *web3.Web3) {
	blockNumber, err := client.Backend.GetBlockNumber()
	if err != nil {
		panic(err)
	}
	meta, err := client.Backend.GetBlockMetadata(blockNumber)
	if err != nil {
		panic(err)
	}

	txParams := &types.TransactParams{
		HexPrivateKey:   hexPrivateKey,
		ValidUntilBlock: blockNumber.Add(blockNumber, big.NewInt(88)),
		Nonce:           "test_deploy_contract",
		ChainID:         meta.ChainID,
		Version:         0,
		Quota:           big.NewInt(10000000),
	}
	txHash, contract, err := client.Backend.DeployContract(context.Background(), txParams, TokenABI, TokenBin)
	if err != nil {
		panic(err)
	}
	fmt.Printf("deploy contract tx hash is %s\n", txHash.Hex())
	fmt.Printf("contract address  is %s\n", contract.Address.Hex())

	callContract(client, &contract.Address)
}

func callContract(client *web3.Web3, address *common.Address) {
	result := new(string)
	from, err := utils.PrivateKeyToAddress(hexPrivateKey)
	if err != nil {
		panic(err)
	}

	abiBuild, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		panic(err)
	}
	abiB, err := abiBuild.Pack("getBalance", from)
	if err != nil {
		panic(err)
	}

	msg := backend.CallMsg{
		From: from,
		To:   address,
		Data: abiB,
	}
	if err := client.Backend.CallContract(context.Background(), result, msg, nil); err != nil {
		panic(err)
	}

	fmt.Printf("balance is %s", *result)
}
