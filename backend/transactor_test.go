package backend

import (
	"context"
	"math/big"
	"testing"

	"github.com/cryptape/go-web3/providers"
	"github.com/cryptape/go-web3/types"
)

// TokenABI is the input ABI used to generate the binding from.
const tokenABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTest\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"supplyAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

// TokenBin is the compiled bytecode used for deploying new contracts.
const tokenBin = `0x608060405234801561001057600080fd5b5060405160208061028383398101604090815290513360009081526020819052919091205561023f806100446000396000f3006080604052600436106100615763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166327e235e38114610066578063a8cd0a80146100a6578063a9059cbb146100bb578063f8b2cb4f14610100575b600080fd5b34801561007257600080fd5b5061009473ffffffffffffffffffffffffffffffffffffffff6004351661012e565b60408051918252519081900360200190f35b3480156100b257600080fd5b50610094610140565b3480156100c757600080fd5b506100ec73ffffffffffffffffffffffffffffffffffffffff60043516602435610145565b604080519115158252519081900360200190f35b34801561010c57600080fd5b5061009473ffffffffffffffffffffffffffffffffffffffff600435166101eb565b60006020819052908152604090205481565b607b90565b3360009081526020819052604081205482118015906101645750600082115b156101e157336000818152602081815260408083208054879003905573ffffffffffffffffffffffffffffffffffffffff871680845292819020805487019055805186815290519293927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a35060016101e5565b5060005b92915050565b73ffffffffffffffffffffffffffffffffffffffff16600090815260208190526040902054905600a165627a7a723058203a29026d7b39ad12e051828a190354bac8abd1a710144c33afb3489ab9434e9c0029`

const hexPrivateKey = `936bcd7f107a58ebf3c5f1b378f23741b552747bd0630b52ff53842cae23e41e`

var txBE = New(providers.NewHTTPProviders("http://121.196.200.225:1337"))

func TestDeployContract(t *testing.T) {
	blockNumber, err := txBE.GetBlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	meta, err := txBE.GetBlockMetadata(blockNumber)
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

	txHash, bc, err := txBE.DeployContract(context.Background(), params, tokenABI, tokenBin)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(txHash.Hex(), bc.Address.Hex())
}
