/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/sdk-go/plugins/types"
)

/**
 *  AES Key
 *
 *      keyInfo format: {
 *          algorithm: "AES",
 *          keySize  : 32,                // optional
 *          data     : "{BASE64_ENCODE}}" // password data
 *          iv       : "{BASE64_ENCODE}", // initialization vector
 *      }
 */
type AESKey struct {
	BaseSymmetricKey

	_data []byte
	_iv []byte
}

func NewAESKey(dict map[string]interface{}) *AESKey {
	return new(AESKey).Init(dict)
}

func (key *AESKey) Init(dict map[string]interface{}) *AESKey {
	if key.BaseSymmetricKey.Init(dict) != nil {
		// TODO: check algorithm parameters
		// 1. check mode = 'CBC'
		// 2. check padding = 'PKCS7Padding'
	}
	return key
}

func (key *AESKey) keySize() uint {
	// TODO: get from key data

	size, ok := key.Get("keySize").(uint)
	if ok {
		return size
	} else {
		return 32
	}
}

func (key *AESKey) blockSize() uint {
	// TODO: get from iv data

	size, ok := key.Get("blockSize").(uint)
	if ok {
		return size
	} else {
		return aes.BlockSize
	}
}

func (key *AESKey) initVector() []byte {
	if key._iv == nil {
		iv := key.Get("iv")
		if iv == nil {
			iv = key.Get("I")
		}
		if iv == nil {
			// zero iv
			zeros := make([]byte, key.blockSize())
			key.Set("iv", Base64Encode(zeros))
			key._iv = zeros
		} else {
			key._iv = Base64Decode(iv.(string))
		}
	}
	return key._iv
}

//-------- ICryptographyKey

func (key *AESKey) Data() []byte {
	if key._data == nil {
		data := key.Get("data")
		if data == nil {
			data = key.Get("D")
		}
		if data == nil {
			//
			// key data empty? generate new key info
			//
			pw := RandomBytes(key.keySize())
			iv := RandomBytes(key.blockSize())
			key.Set("data", Base64Encode(pw))
			key.Set("iv", Base64Encode(iv))
			// other parameters
			//key.Set("mode", "CBC");
			//key.Set("padding", "PKCS7");
			key._data = pw
			key._iv = iv
		} else {
			key._data = Base64Decode(data.(string))
		}
	}
	return key._data
}

//-------- ISymmetricKey(IEncryptKey)

func (key *AESKey) Encrypt(plaintext []byte) []byte {
	block, err := aes.NewCipher(key.Data())
	if err != nil {
		panic(err)
	}
	blockMode := cipher.NewCBCEncrypter(block, key.initVector())
	padded := PKCS5Padding(plaintext, key.blockSize())
	ciphertext := make([]byte, len(padded))
	blockMode.CryptBlocks(ciphertext, padded)
	return ciphertext
}

//-------- ISymmetricKey(IDecryptKey)

func (key *AESKey) Decrypt(ciphertext []byte) []byte {
	block, err := aes.NewCipher(key.Data())
	if err != nil {
		panic(err)
	}
	blockMode := cipher.NewCBCDecrypter(block, key.initVector())
	plaintext := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(plaintext, ciphertext)
	return PKCS5UnPadding(plaintext)
}

//
//  PKCS5
//

func PKCS5Padding(src []byte, blockSize uint) []byte {
	padding := int(blockSize) - len(src) % int(blockSize)
	tail := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, tail...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	count := int(src[length-1])
	return src[:(length - count)]
}
