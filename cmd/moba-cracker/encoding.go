package main

import (
	"math/big"

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
	debug.SetTitle("VariantBase64Encode")
	debug.Print("Input:", string(input))
	var result []byte
	base := big.NewInt(65)              // Base 65
	num := new(big.Int).SetBytes(input) // Convert input to a big integer

	// Perform the base-65 encoding
	for num.Cmp(big.NewInt(0)) > 0 {
		remainder := new(big.Int)
		num.DivMod(num, base, remainder) // num = num / 65, remainder = num % 65
		result = append([]byte{VariantBase64Table[remainder.Int64()]}, result...)
	}

	debug.ResetTitle()
	debug.SetStacktrace(false)
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
