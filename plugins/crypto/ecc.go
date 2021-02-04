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
 *  ECC Public Key
 *
 *      keyInfo format: {
 *          algorithm    : "ECC",
 *          curve        : "secp256k1",
 *          data         : "..." // base64_encode()
 *      }
 */
type ECCPublicKey struct {
	Dictionary
	PublicKey
}

func NewECCPublicKey(dict map[string]interface{}) *ECCPublicKey {
	return new(ECCPublicKey).Init(dict)
}

func (key *ECCPublicKey) Init(dict map[string]interface{}) *ECCPublicKey {
	if key.Dictionary.Init(dict) != nil {
		// init
	}
	return key
}

func (key ECCPublicKey) Data() []byte {
	// TODO:
	return nil
}

func (key ECCPublicKey) Verify(data []byte, signature []byte) bool {
	// TODO:
	return false
}

func (key ECCPublicKey) Match(sKey SignKey) bool {
	return AsymmetricKeysMatch(sKey, key)
}

/**
 *  ECC Private Key
 *
 *      keyInfo format: {
 *          algorithm    : "ECC",
 *          curve        : "secp256k1",
 *          data         : "..." // base64_encode()
 *      }
 */
type ECCPrivateKey struct {
	Dictionary
	PrivateKey
}

func NewECCPrivateKey(dict map[string]interface{}) *ECCPrivateKey {
	return new(ECCPrivateKey).Init(dict)
}

func (key *ECCPrivateKey) Init(dict map[string]interface{}) *ECCPrivateKey {
	if key.Dictionary.Init(dict) != nil {
		// init
	}
	return key
}

func (key ECCPrivateKey) Data() []byte {
	// TODO:
	return nil
}

func (key ECCPrivateKey) Sign(data []byte) []byte {
	// TODO:
	return nil
}
