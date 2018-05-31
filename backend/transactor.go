package backend

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/cryptape/go-web3/errors"
	"github.com/cryptape/go-web3/types"
	"github.com/cryptape/go-web3/utils"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/golang/protobuf/proto"
)

const (
	sendTransactionMethod = "cita_sendTransaction"

	saveAbiAddress = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

// SendTransaction injects the transaction into the pending pool for execution.
func (b *backend) SendTransaction(ctx context.Context, tx *types.Transaction, hexPrivateKey string) (common.Hash, error) {
	sign, err := genSign(tx, hexPrivateKey)
	if err != nil {
		return common.Hash{}, err
	}

	unTx := &types.UnverifiedTransaction{
		Transaction: tx,
		Signature:   sign,
		Crypto:      types.Crypto_SECP,
	}

	data, err := proto.Marshal(unTx)
	if err != nil {
		return common.Hash{}, err
	}

	resp, err := b.provider.SendRequest(sendTransactionMethod, hex.EncodeToString(data))
	if err != nil {
		return common.Hash{}, err
	}

	var status types.TransactionStatus
	if err := resp.GetObject(&status); err != nil {
		return common.Hash{}, err
	}

	if status.Status != "OK" {
		return common.Hash{}, fmt.Errorf("send transaction failed, status is %s", status.Status)
	}
	return common.HexToHash(status.Hash), nil
}

// DeployContract deploys a contract onto the Ethereum blockchain and binds the
// deployment address with a Go wrapper.
func (b *backend) DeployContract(ctx context.Context, params *types.TransactParams, abiStr, code string) (common.Hash, *BoundContract, error) {
	codeB, err := hex.DecodeString(utils.CleanHexPrefix(code))
	if err != nil {
		return common.Hash{}, nil, err
	}

	tx := &types.Transaction{
		To:              "",
		Data:            codeB,
		ValidUntilBlock: params.ValidUntilBlock.Uint64(),
		ChainId:         params.ChainID,
		Nonce:           params.Nonce,
		Quota:           params.Quota.Uint64(),
	}

	txHash, err := b.SendTransaction(ctx, tx, params.HexPrivateKey)
	if err != nil {
		return common.Hash{}, nil, err
	}

	receipt, err := b.waitOnBlock(ctx, txHash)
	if err != nil {
		return common.Hash{}, nil, err
	}

	contractAddress := common.HexToAddress(receipt.ContractAddress)
	if err = b.saveABI(ctx, params.ChainID, params.Version, params.HexPrivateKey, abiStr, contractAddress); err != nil {
		return common.Hash{}, nil, err
	}

	abiParsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return common.Hash{}, nil, err
	}
	bc := &BoundContract{
		Address: contractAddress,
		Abi:     abiParsed,
	}

	return txHash, bc, nil
}

func (b *backend) saveABI(
	ctx context.Context,
	chainID, version uint32,
	hexPrivateKey, abiStr string,
	contractAddress common.Address,
) error {
	blockNumber, err := b.GetBlockNumber()
	if err != nil {
		return err
	}

	abiHex := hex.EncodeToString([]byte(abiStr))
	data, err := hex.DecodeString(utils.CleanHexPrefix(contractAddress.Hex()) + abiHex)
	if err != nil {
		return err
	}
	tx := &types.Transaction{
		To:              saveAbiAddress,
		Data:            data,
		Quota:           99999999,
		ChainId:         chainID,
		Version:         version,
		ValidUntilBlock: blockNumber.Add(blockNumber, big.NewInt(88)).Uint64(),
	}

	txHash, err := b.SendTransaction(ctx, tx, hexPrivateKey)
	if err != nil {
		return err
	}
	_, err = b.waitOnBlock(ctx, txHash)
	return err
}

func (b *backend) waitOnBlock(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	sub, err := b.SubscribeBlockFilter(ctx, func(hashes []common.Hash) (bool, error) {
		if len(hashes) <= 0 {
			return false, nil
		}

		receipt, err = b.TransactionReceipt(ctx, txHash)
		if err != nil && !errors.IsNull(err) {
			return false, err
		}

		if receipt != nil {
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return nil, err
	}

	if err = sub.Quit(); err != nil {
		return nil, err
	}

	return receipt, nil
}

func genSign(tx *types.Transaction, hexPrivateKey string) ([]byte, error) {
	if strings.HasPrefix(hexPrivateKey, "0x") {
		hexPrivateKey = strings.TrimPrefix(hexPrivateKey, "0x")
	}

	txB, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}

	h := sha3.New256()
	h.Write(txB)
	hash := h.Sum(nil)
	return crypto.Sign(hash, privateKey)
}
