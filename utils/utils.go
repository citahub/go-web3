package utils

import "strconv"

func ParseHexToUint64(s string) (uint64, error) {
	num, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return 0, err
	}

	return uint64(num), nil
}

func ConvUint64ToHex(num uint64) string {
	return strconv.FormatUint(num, 16)
}
