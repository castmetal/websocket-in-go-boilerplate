package common

func IsNullOrEmpty(value string) bool {
	return len(value) <= 0
}

func IsNullOrEmptyByte(value []byte) bool {
	return len(value) <= 0
}
