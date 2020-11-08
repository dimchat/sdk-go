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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
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
	Dictionary
	SymmetricKey
}

func (key *AESKey) Init(dictionary map[string]interface{}) *AESKey {
	if key.Dictionary.Init(dictionary) != nil {
		// init
	}
	return key
}

func (key AESKey) Data() []byte {
	data := key.Get("data")
	if data == nil {
		data = key.Get("D")
	}
	if data != nil {
		return Base64Decode(data.(string))
	}

	//
	// key data empty? generate new key info
	//

	// TODO:
	return nil
}

func (key AESKey) Encrypt(plaintext []byte) []byte {
	// TODO:
	return nil
}

func (key AESKey) Decrypt(ciphertext []byte) []byte {
	// TODO:
	return nil
}
