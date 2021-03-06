/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
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
 *  Symmetric key for broadcast message,
 *  which will do nothing when en/decoding message data
 */
var PLAIN = "PLAIN"

type PlainKey struct {
	BaseSymmetricKey

	_data []byte
}

func NewPlainKey() *PlainKey {
	dict := make(map[string]interface{})
	dict["algorithm"] = PLAIN
	key := new(PlainKey).Init(dict)
	ObjectRetain(key)
	return key
}

func (key *PlainKey) Init(dict map[string]interface{}) *PlainKey {
	if key.BaseSymmetricKey.Init(dict) != nil {
		key._data = make([]byte, 0)
	}
	return key
}

func (key *PlainKey) Data() []byte {
	return key._data
}

func (key *PlainKey) Encrypt(plaintext []byte) []byte {
	return plaintext
}

func (key *PlainKey) Decrypt(ciphertext []byte) []byte {
	return ciphertext
}

//
//  Singleton
//
var sharedPlainKey *PlainKey = nil

func GetPlainKey() *PlainKey {
	return sharedPlainKey
}

func init() {
	sharedPlainKey = NewPlainKey()
}
