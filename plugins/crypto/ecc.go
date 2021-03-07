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
	. "github.com/dimchat/mkm-go/digest"
	. "github.com/dimchat/mkm-go/format"
	"github.com/dimchat/sdk-go/plugins/crypto/secp256k1"
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

//-------- ICryptographyKey

func (key *ECCPublicKey) Data() []byte {
	if key._data == nil {
		data := key.Get("data")
		str := data.(string)
		// check for raw data (33/65 bytes)
		length := len(str)
		if length == 66 || length == 130 {
			// Hex format
			key._data = HexDecode(str)
		}
		// TODO: PEM format
	}
	return key._data
}

//-------- IPublicKey(IVerifyKey)

func (key *ECCPublicKey) Verify(data []byte, signature []byte) bool {
	if len(signature) > 64 {
		signature = secp256k1.SignatureFromDER(signature)
	}
	pub := key.Data()
	if len(pub) == 65 {
		pub = pub[1:]
	}
	return secp256k1.Verify(pub, SHA256(data), signature)
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

	_publicKey PublicKey
}

func NewECCPrivateKey(dict map[string]interface{}) *ECCPrivateKey {
	return new(ECCPrivateKey).Init(dict)
}

func (key *ECCPrivateKey) Init(dict map[string]interface{}) *ECCPrivateKey {
	if key.BasePrivateKey.Init(dict) != nil {
		// lazy load
		key._data = nil
		key._publicKey = nil
	}
	return key
}

//-------- ICryptographyKey

func (key *ECCPrivateKey) Data() []byte {
	if key._data == nil {
		data := key.Get("data")
		if data == nil {
			// generate key
			_, pri := secp256k1.Generate()
			key._data = pri
		} else {
			// parse PEM file content
			str := data.(string)
			// check for raw data (32 bytes)
			length := len(str)
			if length == 64 {
				// Hex format
				key._data = HexDecode(str)
			}
			// TODO: PEM format
		}
	}
	return key._data
}

//-------- IPrivateKey(ISignKey)

func (key *ECCPrivateKey) Sign(data []byte) []byte {
	sig := secp256k1.Sign(key.Data(), SHA256(data))
	return secp256k1.SignatureToDER(sig)
}

//-------- IPrivateKey

func (key *ECCPrivateKey) PublicKey() PublicKey {
	if key._publicKey == nil {
		pri := key.Data()
		pub := secp256k1.GetPublicKey(pri)
		data := "04" + HexEncode(pub)
		key._publicKey = PublicKeyParse(map[string]interface{}{
			"algorithm": ECC,
			"data": data,
			"curve": "secp256k1",
			"digest": "SHA256",
		})
	}
	return key._publicKey
}
