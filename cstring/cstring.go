package cstring

func Int8ToString(value []int8) string {
	return string(Int8ToBytes(value))
}

func Int8ToBytes(value []int8) []byte {
	length := calculateLen(value)
	output := make([]byte, length)
	for index := 0; index < length; index++ {
		output[index] = byte(value[index])
	}
	return output
}

func calculateLen(value []int8) int {
	length := 0
	maxLength := len(value)
	for length = 0; length < maxLength && value[length] != 0; length++ {
	}
	return length
}
