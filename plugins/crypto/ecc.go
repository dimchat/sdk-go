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
 *          data         : "..." // base64_encode()
 *      }
 */
type ECCPublicKey struct {
	BasePublicKey

	_data []byte
}

func NewECCPublicKey(dict map[string]interface{}) PublicKey {
	key := new(ECCPublicKey)
	key.Init(dict)
	return key
}

func (key *ECCPublicKey) Init(dict map[string]interface{}) PublicKey {
	if key.BasePublicKey.Init(dict) != nil {
		// lazy load
		key._data = nil
	}
	return key
}

//-------- ICryptographyKey

func (key *ECCPublicKey) Data() []byte {
	if key._data == nil {
		data, _ := key.Get("data").(string)
		// check for raw data (33/65 bytes)
		dataLen := len(data)
		if dataLen == 66 || dataLen == 130 {
			// Hex format
			key._data = HexDecode(data)
		}
		// TODO: PEM format?
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
 *          data         : "..." // base64_encode()
 *      }
 */
type ECCPrivateKey struct {
	BasePrivateKey

	_data []byte

	_publicKey PublicKey
}

func NewECCPrivateKey(dict map[string]interface{}) PrivateKey {
	key := new(ECCPrivateKey)
	key.Init(dict)
	return key
}

func (key *ECCPrivateKey) Init(dict map[string]interface{}) PrivateKey {
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
		data, _ := key.Get("data").(string)
		dataLen := len(data)
		if dataLen == 0 {
			// generate key
			_, pri := secp256k1.Generate()
			key._data = pri
			key.Set("data", HexEncode(pri))
			key.Set("curve", "SECP256k1")
			key.Set("digest", "SHA256")
		} else if dataLen == 64 {
			// check for raw data (32 bytes)
			// Hex format
			key._data = HexDecode(data)
		}
		// TODO: PEM format?
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
			"curve": "SECP256k1",
			"digest": "SHA256",
		})
	}
	return key._publicKey
}
