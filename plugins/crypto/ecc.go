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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/digest"
	. "github.com/dimchat/mkm-go/format"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
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
	BasePublicKey

	_data []byte
}

func NewECCPublicKey(dict map[string]interface{}) *ECCPublicKey {
	return new(ECCPublicKey).Init(dict)
}

func (key *ECCPublicKey) Init(dict map[string]interface{}) *ECCPublicKey {
	if key.BasePublicKey.Init(dict) != nil {
		// lazy load
		key._data = nil
	}
	return key
}

func (key *ECCPublicKey) Data() []byte {
	if key._data == nil {
		data := key.Get("data")
		if data == nil {
			return nil
		}
		str, ok := data.(string)
		if ok {
			// check for raw data (33/65 bytes)
			length := len(str)
			if length == 66 || length == 130 {
				// Hex format
				key._data = HexDecode(str)
			}
			// TODO: PEM format
		}
	}
	return key._data
}

func (key *ECCPublicKey) Verify(data []byte, signature []byte) bool {
	pub := key.Data()
	if pub == nil {
		return false
	}
	return secp256k1.VerifySignature(pub, SHA256(data), signature)
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
	BasePrivateKey

	_data []byte

	_pri *ecdsa.PrivateKey
	_publicKey PublicKey
}

func NewECCPrivateKey(dict map[string]interface{}) *ECCPrivateKey {
	return new(ECCPrivateKey).Init(dict)
}

func (key *ECCPrivateKey) Init(dict map[string]interface{}) *ECCPrivateKey {
	if key.BasePrivateKey.Init(dict) != nil {
		// lazy load
		key._data = nil
		key._pri = nil
		key._publicKey = nil
	}
	return key
}

func (key *ECCPrivateKey) getCurve() *secp256k1.BitCurve {
	return secp256k1.S256()
}

func (key *ECCPrivateKey) generate() []byte {
	pri, err := ecdsa.GenerateKey(key.getCurve(), rand.Reader)
	if err != nil {
		panic(err)
	}
	key._pri = pri
	data := make([]byte, 32)
	blob := pri.D.Bytes()
	copy(data[32-len(blob):], blob)
	return data
}

func (key *ECCPrivateKey) getPrivateKey() *ecdsa.PrivateKey {
	if key._pri == nil {
		data := key.Data()
		pri, err := crypto.ToECDSA(data)
		if err != nil {
			panic(err)
		}
		key._pri = pri
	}
	return key._pri
}

func (key *ECCPrivateKey) Data() []byte {
	if key._data == nil {
		data := key.Get("data")
		if data == nil {
			// generate key
			key._data = key.generate()
		} else {
			// parse PEM file content
			str, ok := data.(string)
			if ok {
				// check for raw data (32 bytes)
				length := len(str)
				if length == 64 {
					// Hex format
					key._data = HexDecode(str)
				}
				// TODO: PEM format
			}
		}
	}
	return key._data
}

func (key *ECCPrivateKey) Sign(data []byte) []byte {
	pub := key.Data()
	if pub == nil {
		return nil
	}
	sig, err := secp256k1.Sign(SHA256(data), pub)
	if err !=  nil {
		panic(err)
	}
	return sig
}

func (key *ECCPrivateKey) PublicKey() PublicKey {
	if key._publicKey == nil {
		pri := key.getPrivateKey()
		if pri == nil {
			return nil
		}
		pub := elliptic.Marshal(key.getCurve(), pri.X, pri.Y)
		key._publicKey = PublicKeyParse(map[string]interface{}{
			"algorithm": ECC,
			"data": UTF8Decode(pub),
			"curve": "secp256k1",
			"digest": "SHA256",
		})
	}
	return key._publicKey
}
