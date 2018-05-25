package types

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
