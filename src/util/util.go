package util

import (
	"strconv"
)

func RemoveEmpty(values []string) []string {
	var result []string
	for _, str := range values {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

func StringToUint16(value string) (uint16, error) {
	i, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}

func Uint16ToString(value uint16) string {
	return strconv.FormatUint(uint64(value), 10)
}
