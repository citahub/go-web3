/*
Copyright 2016-2017 Cryptape Technologies LLC.

This program is free software: you can redistribute it
and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation,
either version 3 of the License, or (at your option) any
later version.

This program is distributed in the hope that it will be
useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package types

import (
	"math/big"
)

const AddressLength = 20

type Block struct {
	Version int          `json:"version"`
	Hash    string       `json:"hash"`
	Header  *BlockHeader `json:"header"`
	Body    *BlockBody   `json:"body"`
}

type BlockHeader struct {
	Timestamp        int64       `json:"timestamp"`
	PrevHash         string      `json:"prevHash"`
	Proof            *BlockProof `json:"proof"`
	StateRoot        string      `json:"stateRoot"`
	TransactionsRoot string      `json:"transactionsRoot"`
	ReceiptsRoot     string      `json:"receiptsRoot"`
	GasUsed          string      `json:"gasUsed"`
	Number           string      `json:"number"`
	Proposer         string      `json:"proposer"`
}

type BlockProof struct {
	Proposal string   `json:"proposal"`
	Height   uint     `json:"height"`
	Round    uint     `json:"round"`
	Commits  []string `json:"commits"`
}

type BlockBody struct {
	Transactions []*BlockTransaction `json:"transactions"`
}

type BlockTransaction struct {
	Hash    string `json:"hash"`
	Content string `json:"content"`
}

type BlockMetadata struct {
	ChainID          uint32   `json:"chainId"`
	ChainName        string   `json:"chainName"`
	Operator         string   `json:"operator"`
	GenesisTimestamp uint64   `json:"genesisTimestamp"`
	Validators       []string `json:"validators"`
	BlockInterval    uint64   `json:"blockInterval"`
}

type TransactionStatus struct {
	Hash   string `json:"hash"`
	Status string `json:"status"`
}

type Receipt struct {
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	ContractAddress   string `json:"contractAddress"`
	Logs              []*Log `json:"logs"`
	Root              string `json:"root"`
	LogsBloom         string `json:"logsBloom"`
	ErrorMessage      string `json:"errorMessage"`
}

type Log struct {
	Address             string   `json:"address"`
	Topics              []string `json:"topics"`
	Data                string   `json:"data"`
	BlockHash           string   `json:"blockHash"`
	BlockNumber         string   `json:"blockNumber"`
	TransactionHash     string   `json:"transactionHash"`
	TransactionIndex    string   `json:"transactionIndex"`
	LogIndex            string   `json:"logIndex"`
	TransactionLogIndex string   `json:"transactionLogIndex"`
}

// TransactParams is the collection of authorization data required to create a
// valid CITA transaction.
type TransactParams struct {
	HexPrivateKey   string
	To              string
	Nonce           string
	ValidUntilBlock *big.Int

	Value *big.Int
	Quota *big.Int

	ChainID uint32
	Version uint32
}
