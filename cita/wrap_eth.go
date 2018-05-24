package cita

const (
	getTransactionReceiptMethod = "eth_getTransactionReceipt"
)

func (c *cita) GetTransactionReceipt(hash string) (*Receipt, error) {
	resp, err := c.provider.SendRequest(getTransactionReceiptMethod, hash)
	if err != nil {
		return nil, err
	}

	var r Receipt
	if err := resp.GetObject(&r); err != nil {
		return nil, err
	}

	return &r, nil
}
