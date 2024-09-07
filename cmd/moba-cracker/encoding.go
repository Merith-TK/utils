package main

import (
	"github.com/Merith-TK/utils/debug"
)

var VariantBase64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
var VariantBase64Dict = make(map[int]byte)
var VariantBase64ReverseDict = make(map[byte]int)

var LicenseTypes = struct {
	Professional int
	Education    int
	Personal     int
}{
	Professional: 1,
	Education:    3,
	Personal:     4,
}

func init() {
	for i := 0; i < len(VariantBase64Table); i++ {
		VariantBase64Dict[i] = VariantBase64Table[i]
		VariantBase64ReverseDict[VariantBase64Table[i]] = i
	}
	debug.Print("VariantBase64Table:", VariantBase64Table)
	debug.Print("VariantBase64Table Length:", len(VariantBase64Table))
}

func VariantBase64Encode(input []byte) []byte {
	var result []byte
	blocksCount := len(input) / 3
	leftBytes := len(input) % 3

	for i := 0; i < blocksCount; i++ {
		// Convert 3 bytes to an integer (24 bits)
		codingInt := int(input[3*i]) | int(input[3*i+1])<<8 | int(input[3*i+2])<<16
		// Break the 24-bit integer into four 6-bit values and map them using VariantBase64Table
		result = append(result, VariantBase64Table[codingInt&0x3f])
		result = append(result, VariantBase64Table[(codingInt>>6)&0x3f])
		result = append(result, VariantBase64Table[(codingInt>>12)&0x3f])
		result = append(result, VariantBase64Table[(codingInt>>18)&0x3f])
	}

	// Handle leftover bytes (padding)
	if leftBytes == 1 {
		codingInt := int(input[3*blocksCount])
		result = append(result, VariantBase64Table[codingInt&0x3f])
		result = append(result, VariantBase64Table[(codingInt>>6)&0x3f])
	} else if leftBytes == 2 {
		codingInt := int(input[3*blocksCount]) | int(input[3*blocksCount+1])<<8
		result = append(result, VariantBase64Table[codingInt&0x3f])
		result = append(result, VariantBase64Table[(codingInt>>6)&0x3f])
		result = append(result, VariantBase64Table[(codingInt>>12)&0x3f])
	}

	return result
}

func EncryptBytes(key uint64, bs []byte) []byte {
	debug.Print("Encrypting:", string(bs))
	debug.Print("Key:", key)
	debug.Print("Length:", len(bs))
	result := make([]byte, len(bs))
	for i := 0; i < len(bs); i++ {
		result[i] = bs[i] ^ byte((key>>8)&0xff)
		key = uint64(result[i])&key | 0x482D
	}
	debug.Print("Result:", string(result))
	return result
}

func DecryptBytes(key int, bs []byte) []byte {
	debug.Print("Decrypting:", string(bs))
	result := make([]byte, len(bs))
	for i := 0; i < len(bs); i++ {
		result[i] = bs[i] ^ byte((key>>8)&0xff)
		key = int(bs[i])&key | 0x482D
	}
	return result
}
