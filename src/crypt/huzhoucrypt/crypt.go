package main


/*

*/
import (
	"C"
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"errors"
	"fmt"
)

//export HZEncrypt
func HZEncrypt(key *C.char, value *C.char, blockSize int) *C.char {
	if v, err := AESEncrypt([]byte(C.GoString(key)), []byte(C.GoString(value)), blockSize); err != nil {
		fmt.Println(err)
		return C.CString("")
	} else {
		return C.CString(string(BASE64Encrypt(v)))
	}

}

func padding(value []byte,blockSize int) []byte {
	paddingCount := blockSize - len(value)%blockSize
	if paddingCount == 0 {
		return value
	} else {
		return append(value, bytes.Repeat([]byte{byte(0)}, paddingCount)...)
	}
}

func unpadding(value []byte) []byte {
	for i := len(value) - 1; i > 0; i-- {
		if value[i] != 0 {
			return value[:i+1]
		}
	}
	return value
}

//export AESEncrypt
func AESEncrypt(key []byte, value []byte, blockSize int) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error: NewCipher(%d bytes)=%s", len(key), err))
	}
	src := padding(value,blockSize)
	vlen := len(src)
	ciphertext := make([]byte, vlen)
	tmpData := make([]byte, blockSize)
	for index := 0; index < vlen; index += blockSize {
		block.Encrypt(tmpData, src[index:index+blockSize])
		copy(ciphertext, tmpData)
	}
	return ciphertext, nil
}

func AESDecrypt(key []byte, value string, blockSize int) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decvalue, err := BASE64Decrypt(value)
	if err != nil {
		return nil, err
	}
	vlen := len(decvalue)
	decryptData := make([]byte, vlen)
	tmpData := make([]byte, blockSize)
	for index := 0; index < vlen; index += blockSize {
		block.Decrypt(tmpData, decryptData[index:index+blockSize])
		copy(decryptData, tmpData)
	}
	return unpadding(decryptData), nil
}

//export BASE64Encrypt
func BASE64Encrypt(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

//export BASE64Decrypt
func BASE64Decrypt(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

func main() {
}
