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
	connection := web3.New(providers.NewHTTPProviders("http://127.0.0.1:1337"))
	num, err := connection.Cita.GetBlockNumber()
	if err != nil {
		panic(err)
	}

	fmt.Printf("block number is %d\n", num)
}
```

## Check List
- [ ] JSON RPC
  - [x] block filters
  - [x] log filters
  - [ ] call contract
  - [ ] other
- [ ] Code generation base on [abigen](https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts)

## Documents
TODO
