package cita

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/golang/protobuf/proto"

	"github.com/cryptape/go-web3/internal/proto/transaction"
	"github.com/cryptape/go-web3/providers"
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
	GetBlockMetadata(number uint64) (*BlockMetadata, error)

	GetBlockByHash(hash string, detail bool) (*Block, error)
	GetBlockByNumber(number uint64, detail bool) (*Block, error)

	GetTransactionProof(hash string) (string, error)
	GetTransactionReceipt(hash string) (*Receipt, error)

	SendRawTransaction(data string) (*TransactionStatus, error)
	CreateContract(code, hexPrivateKey, nonce string, quota, validUntilBlock uint64) (*TransactionStatus, error)
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

func (c *cita) GetBlockByHash(hash string, detail bool) (*Block, error) {
	resp, err := c.provider.SendRequest(getBlockByHashMethod, hash, detail)
	if err != nil {
		return nil, err
	}
	var b Block
	if err := resp.GetObject(&b); err != nil {
		return nil, err
	}

	return &b, nil
}

func (c *cita) GetBlockByNumber(number uint64, detail bool) (*Block, error) {
	resp, err := c.provider.SendRequest(getBlockByNumberMethod, utils.ConvUint64ToHex(number), detail)
	if err != nil {
		return nil, err
	}
	var b Block
	if err := resp.GetObject(&b); err != nil {
		return nil, err
	}

	return &b, nil
}

func (c *cita) GetBlockMetadata(number uint64) (*BlockMetadata, error) {
	resp, err := c.provider.SendRequest(getBlockMetadataMethod, utils.ConvUint64ToHex(number))
	if err != nil {
		return nil, err
	}

	var meta BlockMetadata
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

func (c *cita) CreateContract(code, hexPrivateKey, nonce string, quota, validUntilBlock uint64) (*TransactionStatus, error) {
	num, err := c.GetBlockNumber()
	if err != nil {
		return nil, err
	}

	meta, err := c.GetBlockMetadata(num)
	if err != nil {
		return nil, err
	}

	codeB, err := hex.DecodeString(utils.CleanHexPrefix(code))
	if err != nil {
		return nil, err
	}
	tx := &transaction.Transaction{
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

	unTx := &transaction.UnverifiedTransaction{
		Transaction: tx,
		Signature:   sign,
		Crypto:      transaction.Crypto_SECP,
	}

	unTxB, err := proto.Marshal(unTx)
	if err != nil {
		return nil, err
	}

	return c.SendRawTransaction(hex.EncodeToString(unTxB))
}

func (c *cita) SendRawTransaction(data string) (*TransactionStatus, error) {
	resp, err := c.provider.SendRequest(sendTransactionMethod, data)
	if err != nil {
		return nil, err
	}

	var txStatus TransactionStatus
	if err := resp.GetObject(&txStatus); err != nil {
		return nil, err
	}

	return &txStatus, nil
}

func genSign(tx *transaction.Transaction, hexPrivateKey string) ([]byte, error) {
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
