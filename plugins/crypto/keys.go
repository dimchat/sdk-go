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

type BaseKey struct {
	Dictionary
}

//func (key *BaseKey) Init(dict map[string]interface{}) CryptographyKey {
//	if key.Dictionary.Init(dict) != nil {
//	}
//	return key
//}

//-------- ICryptographyKey

func (key *BaseKey) Algorithm() string {
	return CryptographyKeyGetAlgorithm(key.Map())
}

func (key *BaseKey) Data() []byte {
	panic("BaseKey::Data() > override me!")
}

/**
 *  Symmetric Key
 *  ~~~~~~~~~~~~~
 */
type BaseSymmetricKey struct {
	BaseKey
}

//func (key *BaseSymmetricKey) Init(dict map[string]interface{}) SymmetricKey {
//	if key.BaseKey.Init(dict) != nil {
//	}
//	return key
//}

//-------- IEncryptKey

//func (key *BaseSymmetricKey) Encrypt(_ []byte) []byte {
//	panic("BaseSymmetricKey::Encrypt(plaintext) > override me!")
//}

//-------- IDecryptKey

func (key *BaseSymmetricKey) Decrypt(_ []byte) []byte {
	panic("BaseSymmetricKey::Decrypt(ciphertext) > override me!")
}

func (key *BaseSymmetricKey) Match(pKey EncryptKey) bool {
	return SymmetricKeysMatch(pKey, key)
}

/**
 *  Asymmetric Public Key
 *  ~~~~~~~~~~~~~~~~~~~~~
 */
type BasePublicKey struct {
	BaseKey
}

//func (key *BasePublicKey) Init(dict map[string]interface{}) PublicKey {
//	if key.BaseKey.Init(dict) != nil {
//	}
//	return key
//}

//-------- IVerifyKey

func (key *BasePublicKey) Verify(_ []byte, _ []byte) bool {
	panic("BasePublicKey::Verify(data, signature) > override me!")
}

func (key *BasePublicKey) Match(sKey SignKey) bool {
	return AsymmetricKeysMatch(sKey, key)
}

/**
 *  Asymmetric Private Key
 *  ~~~~~~~~~~~~~~~~~~~~~~
 */
type BasePrivateKey struct {
	BaseKey
}

//func (key *BasePrivateKey) Init(dict map[string]interface{}) PrivateKey {
//	if key.BaseKey.Init(dict) != nil {
//	}
//	return key
//}

//-------- IPrivateKey

//func (key *BasePrivateKey) PublicKey() PublicKey {
//	panic("BasePrivateKey::PublicKey() > override me!")
//}

//-------- ISignKey

//func (key *BasePrivateKey) Sign(_ []byte) []byte {
//	panic("BasePrivateKey::Sign(data) > override me!")
//}
