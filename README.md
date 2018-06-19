# go-web3
>A CITA JSON-RPC implementation for Go. The usage is similar to eth-web3.

## Check List
- [x] JSON RPC
  - [x] block filters
  - [x] log filters
  - [x] call contract
  - [x] deploy contract
- [ ] Code generation base on [abigen](https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts)

## Usage
这里只介绍合约, filter 相关的用法, 更多接口可以查看 [JSON RPC 文档](https://cryptape.github.io/cita/zh/usage-guide/rpc/index.html#json-rpc)

### 部署合约
`token.sol` 是一个非常简单的代币合约文件，只包含了两个方法，查询某一地址的余额以及向某一地址转账，代码如下:
```solidity
pragma solidity ^0.4.18;

contract Token {
    mapping (address => uint) public balances;

    event Transfer(address indexed _from, address indexed _to, uint256 _value);

    function Token(uint supplyAmount) {
        balances[msg.sender] = supplyAmount;
    }

    function getBalance(address account) constant public returns (uint balance) {
        return balances[account];
    }

    function transfer(address _to, uint256 _value) public returns (bool success) {
        if (balances[msg.sender] >= _value && _value > 0) {
            balances[msg.sender] -= _value;
            balances[_to] += _value;
            Transfer(msg.sender, _to, _value);
            return true;
        } else {
            return false;
        }
    }
}
```
使用 `solc` 编译获得二进制和 ABI
```
$solc --bin --abi token.sol
======= token.sol:Token =======
Binary:
....
Contract JSON ABI:
...
```
构造交易并发送
```go
bytecode := ... #contract bytecode
abi := ... #contract abi
ctx := ... #context.Context
hexPrivateKey := `0x3b181467f7f7bc58bccd223efd9b5b8afa1176f5ef9b619a4421f8ac2ce1400b`

client := web3.New(providers.NewHTTPProviders("http://127.0.0.1:1337"))

blockNumber, _ := client.Backend.GetBlockNumber()
meta, _ := client.Backend.GetBlockMetadata(blockNumber)

txParams := &types.TransactParams{
	HexPrivateKey:   hexPrivateKey,
	ValidUntilBlock: blockNumber.Add(blockNumber, big.NewInt(88)),
	Nonce:           "test_deploy_contract",
	ChainID:         meta.ChainID,
	Version:         0,
	Quota:           big.NewInt(10000000),
}
hash, contract, err := client.Backend.DeployContract(ctx, txParams, abi, bytecode)
fmt.Println(hash) // tx hash
fmt.Println(contract.Address) // contract address
```
部署成功后会返回交易 hash 以及合约地址, ABI 等信息. 并且会自动的把合约代码存储到链上.

### 调用合约
#### eth_call
合约调用需要指定合约地址, 函数, 参数.
```go
from := ... #address
data, _ := contract.Abi.Pack("getBalance", from)

msg := backend.CallMsg{
	From: from,
	To:   contract.Address,
	Data: data,
}
result := new(string)
client.Backend.CallContract(ctx, result, msg, nil)
fmt.Println(*result) // address balance
```

### Filters
#### Block filters
创建一个过滤器, 每当有新的块出现时会发送通知.
```go
client := web3.New(providers.NewHTTPProviders("http://127.0.0.1:1337"))

ctx := ... #context.Context
sub, _ := client.Backend.SubscribeBlockFilter(ctx, func(hashs []common.Hash) (bool, error) {
	for _, hash := range hashs {
		fmt.Println(hash.Hex()) // block hash
	}

	return true, nil // return true stop subscription
})

if err := sub.Quit(); err != nil {
	fmt.Println(err)
}
```

#### Log filters
创建一个日志过滤器
```go
query := ethereum.FilterQuery{
  Addresses: []common.Address{contract.Address},
}
sub, err := be.SubscribeLogFilter(context.Background(), query, func(logs []types.Log) (bool, error) {
  t.Log(logs)
  for _, log := range logs {
    fmt.Println(log)
  }

  return true, nil
})
if err := sub.Quit(); err != nil {
  fmt.Println(err)
}
```
