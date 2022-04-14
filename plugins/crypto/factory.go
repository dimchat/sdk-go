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
)

/**
 *  Symmetric Key
 */
type SymmetricKeyGenerator func() SymmetricKey
type SymmetricKeyParser func(map[string]interface{}) SymmetricKey

// symmetric key factory
type GeneralSymmetricKeyFactory struct {

	_generate SymmetricKeyGenerator
	_parse SymmetricKeyParser
}

func NewSymmetricKeyFactory(g SymmetricKeyGenerator, p SymmetricKeyParser) SymmetricKeyFactory {
	factory := new(GeneralSymmetricKeyFactory)
	factory.Init(g, p)
	return factory
}

func (factory *GeneralSymmetricKeyFactory) Init(g SymmetricKeyGenerator, p SymmetricKeyParser) SymmetricKeyFactory {
	factory._generate = g
	factory._parse = p
	return factory
}

//-------- ISymmetricKeyFactory

func (factory *GeneralSymmetricKeyFactory) GenerateSymmetricKey() SymmetricKey {
	return factory._generate()
}

func (factory *GeneralSymmetricKeyFactory) ParseSymmetricKey(key map[string]interface{}) SymmetricKey {
	return factory._parse(key)
}

func RegisterSymmetricKeyFactories() {
	// AES
	SymmetricKeySetFactory(AES, NewSymmetricKeyFactory(func() SymmetricKey {
		dict := make(map[string]interface{})
		dict["algorithm"] = AES
		return NewAESKey(dict)
	}, func(dict map[string]interface{}) SymmetricKey {
		return NewAESKey(dict)
	}))
	// PLAIN
	SymmetricKeySetFactory(PLAIN, NewSymmetricKeyFactory(func() SymmetricKey {
		return GetPlainKey()
	}, func(dict map[string]interface{}) SymmetricKey {
		return GetPlainKey()
	}))
}

/**
 *  Asymmetric Key
 */
type PrivateKeyGenerator func() PrivateKey
type PrivateKeyParser func(map[string]interface{}) PrivateKey
type PublicKeyParser func(map[string]interface{}) PublicKey

// private key factory
type GeneralPrivateKeyFactory struct {

	_generate PrivateKeyGenerator
	_parse PrivateKeyParser
}

func NewPrivateKeyFactory(g PrivateKeyGenerator, p PrivateKeyParser) PrivateKeyFactory {
	factory := new(GeneralPrivateKeyFactory)
	factory.Init(g, p)
	return factory
}

func (factory *GeneralPrivateKeyFactory) Init(g PrivateKeyGenerator, p PrivateKeyParser) PrivateKeyFactory {
	factory._generate = g
	factory._parse = p
	return factory
}

//-------- IPrivateKeyFactory

func (factory *GeneralPrivateKeyFactory) GeneratePrivateKey() PrivateKey {
	return factory._generate()
}

func (factory *GeneralPrivateKeyFactory) ParsePrivateKey(key map[string]interface{}) PrivateKey {
	return factory._parse(key)
}

// public key factory
type GeneralPublicKeyFactory struct {

	_parse PublicKeyParser
}

func NewPublicKeyFactory(p PublicKeyParser) PublicKeyFactory {
	factory := new(GeneralPublicKeyFactory)
	factory.Init(p)
	return factory
}

func (factory *GeneralPublicKeyFactory) Init(p PublicKeyParser) PublicKeyFactory {
	factory._parse = p
	return factory
}

//-------- IPublicKeyFactory

func (factory *GeneralPublicKeyFactory) ParsePublicKey(key map[string]interface{}) PublicKey {
	return factory._parse(key)
}

func RegisterPublicKeyFactories() {
	// RSA
	PublicKeySetFactory(RSA, NewPublicKeyFactory(func(dict map[string]interface{}) PublicKey {
		return NewRSAPublicKey(dict)
	}))
	// ECC
	PublicKeySetFactory(ECC, NewPublicKeyFactory(func(dict map[string]interface{}) PublicKey {
		return NewECCPublicKey(dict)
	}))
}

func RegisterPrivateKeyFactories() {
	// RSA
	PrivateKeySetFactory(RSA, NewPrivateKeyFactory(func() PrivateKey {
		dict := make(map[string]interface{})
		dict["algorithm"] = RSA
		return NewRSAPrivateKey(dict)
	}, func(dict map[string]interface{}) PrivateKey {
		return NewRSAPrivateKey(dict)
	}))
	// ECC
	PrivateKeySetFactory(ECC, NewPrivateKeyFactory(func() PrivateKey {
		dict := make(map[string]interface{})
		dict["algorithm"] = ECC
		return NewECCPrivateKey(dict)
	}, func(dict map[string]interface{}) PrivateKey {
		return NewECCPrivateKey(dict)
	}))
}
