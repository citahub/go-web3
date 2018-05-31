# go-web3
>A CITA JSON-RPC implementation for Go. The usage is similar to eth-web3.

## Usage
```go
import (
	"fmt"

	"github.com/cryptape/go-web3"
	"github.com/cryptape/go-web3/providers"
)

func main() {
	client := web3.New(providers.NewHTTPProviders("http://47.75.129.215:1337"))
	txHash, contract, _ := client.Backend.DeployContract(...)
	fmt.Printf("deploy contract tx hash is %s\n", txHash.Hex())
	fmt.Printf("contract address  is %s\n", contract.Address.Hex())
}
```

## Check List
- [ ] JSON RPC
  - [x] block filters
  - [x] log filters
  - [x] call contract
  - [x] deploy contract
  - [ ] other
- [ ] Code generation base on [abigen](https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts)

## Documents
TODO
