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
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  RSA Public Key
 *
 *      keyInfo format: {
 *          algorithm: "RSA",
 *          data: "..."       // base64
 *      }
 */
type RSAPublicKey struct {
	Dictionary
	PublicKey
	EncryptKey
}

func (key *RSAPublicKey) Init(dictionary map[string]interface{}) *RSAPublicKey {
	if key.Dictionary.Init(dictionary) != nil {
		// init
	}
	return key
}

func (key RSAPublicKey) Data() []byte {
	// TODO:
	return nil
}

func (key RSAPublicKey) Verify(data []byte, signature []byte) bool {
	// TODO:
	return false
}

func (key RSAPublicKey) Encrypt(plaintext []byte) []byte {
	// TODO:
	return nil
}

/**
 *  RSA Private Key
 *
 *      keyInfo format: {
 *          algorithm    : "RSA",
 *          keySizeInBits: 1024, // optional
 *          data         : "..." // base64_encode()
 *      }
 */
type RSAPrivateKey struct {
	Dictionary
	PrivateKey
	DecryptKey
}

func (key *RSAPrivateKey) Init(dictionary map[string]interface{}) *RSAPrivateKey {
	if key.Dictionary.Init(dictionary) != nil {
		// init
	}
	return key
}

func (key RSAPrivateKey) Data() []byte {
	// TODO:
	return nil
}

func (key RSAPrivateKey) Sign(data []byte) []byte {
	// TODO:
	return nil
}

func (key RSAPrivateKey) Decrypt(ciphertext []byte) []byte {
	// TODO:
	return nil
}
