package util

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"crypto/cipher"
)

/*ASE 解密   ----------------------------------------------分割线*/
func AesCBCDecrypt(cipherText string) (ret string, err error) {
	key := []byte(`aes_key_16_bytes`)
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return ``, err
	}
	if len(cipherBytes) < aes.BlockSize {
		return ``, err
	}
	iv := cipherBytes[:aes.BlockSize]
	if len(cipherBytes)%aes.BlockSize != 0 {
		return ``, fmt.Errorf(`cipherBytes is not a multiple of the block size`)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherBytes, cipherBytes)
	cipherBytes = Unpad(cipherBytes[aes.BlockSize:])
	return string(cipherBytes), nil
}

func Unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
