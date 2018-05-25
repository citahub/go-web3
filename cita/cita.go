package cita

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/golang/protobuf/proto"

	"github.com/cryptape/go-web3/errors"
	"github.com/cryptape/go-web3/providers"
	"github.com/cryptape/go-web3/types"
	"github.com/cryptape/go-web3/utils"
)

const (
	getBlockNumberMethod      = "cita_blockNumber"
	getBlockByHashMethod      = "cita_getBlockByHash"
	getBlockByNumberMethod    = "cita_getBlockByNumber"
	getBlockMetadataMethod    = "cita_getMetaData"
	getTransactionProofMethod = "cita_getTransactionProof"
	sendTransactionMethod     = "cita_sendTransaction"
)

type Interface interface {
	GetBlockNumber() (uint64, error)
	GetBlockMetadata(number uint64) (*types.BlockMetadata, error)

	GetBlockByHash(hash string, detail bool) (*types.Block, error)
	GetBlockByNumber(number uint64, detail bool) (*types.Block, error)

	GetTransactionProof(hash string) (string, error)
	GetTransactionReceipt(hash string) (*Receipt, error)

	SendRawTransaction(data string) (*types.TransactionStatus, error)
	CreateContract(code, hexPrivateKey, nonce string, quota, validUntilBlock uint64, chainID uint32) (*types.Receipt, error)
	ethInterface
}

func New(provider providers.Interface) Interface {
	return &cita{
		provider: provider,
	}
}

type cita struct {
	provider providers.Interface
}

func (c *cita) GetBlockNumber() (uint64, error) {
	resp, err := c.provider.SendRequest(getBlockNumberMethod)
	if err != nil {
		return 0, err
	}

	s, err := resp.GetString()
	if err != nil {
		return 0, err
	}

	return utils.ParseHexToUint64(s)
}

func (c *cita) GetBlockByHash(hash string, detail bool) (*types.Block, error) {
	resp, err := c.provider.SendRequest(getBlockByHashMethod, hash, detail)
	if err != nil {
		return nil, err
	}
	var b types.Block
	if err := resp.GetObject(&b); err != nil {
		return nil, err
	}

	return &b, nil
}

func (c *cita) GetBlockByNumber(number uint64, detail bool) (*types.Block, error) {
	resp, err := c.provider.SendRequest(getBlockByNumberMethod, utils.ConvUint64ToHex(number), detail)
	if err != nil {
		return nil, err
	}
	var b types.Block
	if err := resp.GetObject(&b); err != nil {
		return nil, err
	}

	return &b, nil
}

func (c *cita) GetBlockMetadata(number uint64) (*types.BlockMetadata, error) {
	resp, err := c.provider.SendRequest(getBlockMetadataMethod, utils.ConvUint64ToHex(number))
	if err != nil {
		return nil, err
	}

	var meta types.BlockMetadata
	if err := resp.GetObject(&meta); err != nil {
		return nil, err
	}

	return &meta, nil
}

func (c *cita) GetTransactionProof(hash string) (string, error) {
	resp, err := c.provider.SendRequest(getTransactionProofMethod, hash)
	if err != nil {
		return "", err
	}

	return resp.GetString()
}

func (c *cita) CreateContract(code, hexPrivateKey, nonce string, quota, validUntilBlock uint64, chainID uint32) (*types.Receipt, error) {
	codeB, err := hex.DecodeString(utils.CleanHexPrefix(code))
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{
		To:              "",
		Data:            codeB,
		ValidUntilBlock: num + validUntilBlock,
		ChainId:         meta.ChainID,
		Nonce:           nonce,
		Quota:           quota,
	}

	sign, err := genSign(tx, hexPrivateKey)
	if err != nil {
		return nil, err
	}

	unTx := &types.UnverifiedTransaction{
		Transaction: tx,
		Signature:   sign,
		Crypto:      types.Crypto_SECP,
	}

	unTxB, err := proto.Marshal(unTx)
	if err != nil {
		return nil, err
	}

	return c.SendRawTransaction(hex.EncodeToString(unTxB))
}

func (c *cita) SendRawTransaction(data string) (*types.TransactionStatus, error) {
	resp, err := c.provider.SendRequest(sendTransactionMethod, data)
	if err != nil {
		return nil, err
	}

	var txStatus types.TransactionStatus
	if err := resp.GetObject(&txStatus); err != nil {
		return nil, err
	}

	return &txStatus, nil
}

func (c *cita) WaitTransactionOnBlock(txHash string, validUntilBlock uint64) (*types.Receipt, error) {
	id, err := c.NewBlockFilter()
	if err != nil {
		return nil, err
	}

	defer c.UninstallFilter(id)

	ch := make(chan interface{})
	defer close(ch)

	go func() {
		for {
			hashList, err := c.GetFilterChanges(id)
			if err != nil {
				ch <- err
				return
			}

			if len(hashList) == 0 {
				continue
			}

			receipt, err := c.GetTransactionReceipt(txHash)
			if err != nil && !errors.IsNull(err) {
				ch <- err
				return
			}

			// transaction on block
			if receipt != nil {
				ch <- receipt
				return
			}

			block, err := c.GetBlockByHash(hashList[len(hashList)-1], false)
			if err != nil {
				fmt.Println("where error 3")
				ch <- err
				return
			}

			blockHeight, err := utils.ParseHexToUint64(block.Header.Number)
			if err != nil {
				ch <- err
				return
			}

			if validUntilBlock < blockHeight {
				ch <- fmt.Errorf("The transaction %s cannot be completed at block height of %d", txHash, validUntilBlock)
				return
			}
		}
	}()

	result := <-ch
	if err, ok := result.(error); ok {
		return nil, err
	}
	return result.(*types.Receipt), nil
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
