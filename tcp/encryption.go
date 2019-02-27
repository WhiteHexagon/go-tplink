package tcp

func encrypt(input string) string {
	key := byte(171)
	result := make([]byte, len(input)+4)
	result[3] = byte(len(input)) //not seen any big commands so far...
	//var b byte
	//for i := 0; i < len(input); i++ {
	//b = input[i]
	for i, b := range []byte(input) {
		a := key ^ b
		key = a
		result[i+4] = a
	}
	return string(result)
}

func decrypt(input string) string {
	key := byte(171)
	result := make([]byte, len(input))
	//for i := 0; i < len(input); i++ {
	//	b = input[i]
	for i, b := range []byte(input) {
		a := key ^ b
		key = b
		result[i] = a
	}
	return string(result)
}
