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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/plugins/types"
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
	BasePublicKey
	IEncryptKey

	_rsaPublicKey *rsa.PublicKey

	_data []byte
}

func NewRSAPublicKey(dict map[string]interface{}) *RSAPublicKey {
	key := new(RSAPublicKey).Init(dict)
	ObjectRetain(key)
	return key
}

func (key *RSAPublicKey) Init(dict map[string]interface{}) *RSAPublicKey {
	if key.BasePublicKey.Init(dict) != nil {
		// lazy load
		key._rsaPublicKey = nil
		key._data = nil
	}
	return key
}

func (key *RSAPublicKey) getPublicKey() *rsa.PublicKey {
	if key._rsaPublicKey == nil {
		data := key.Get("data")
		block, _ := pem.Decode(UTF8Encode(data.(string)))
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err !=  nil {
			panic(err)
		}
		key._rsaPublicKey = pub.(*rsa.PublicKey)
	}
	return key._rsaPublicKey
}

func (key *RSAPublicKey) getHash() crypto.Hash {
	return crypto.SHA256
}

//-------- ICryptographyKey

func (key *RSAPublicKey) Data() []byte {
	if key._data == nil {
		// TODO: encode public key data to PKCS1
		pub := key.getPublicKey()
		key._data = pub.N.Bytes()
	}
	return key._data
}

//-------- IPublicKey/IVerifyKey

func (key *RSAPublicKey) Verify(data []byte, signature []byte) bool {
	pub := key.getPublicKey()
	h := key.getHash().New()
	h.Write(data)
	sum := h.Sum(nil)
	err := rsa.VerifyPKCS1v15(pub, key.getHash(), sum, signature)
	return err == nil
}

//-------- IEncryptKey

func (key *RSAPublicKey) Encrypt(plaintext []byte) []byte {
	pub := key.getPublicKey()
	part := pub.N.BitLen() / 8 - 11
	chunks := BytesSplit(plaintext, part)
	buffer := bytes.NewBufferString("")
	for _, line := range chunks {
		data, err := rsa.EncryptPKCS1v15(rand.Reader, pub, line)
		if err != nil {
			panic(err)
		}
		buffer.Write(data)
	}
	return buffer.Bytes()
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
	BasePrivateKey
	IDecryptKey

	_rsaPrivateKey *rsa.PrivateKey

	_data []byte

	_publicKey PublicKey
}

func NewRSAPrivateKey(dict map[string]interface{}) *RSAPrivateKey {
	key := new(RSAPrivateKey).Init(dict)
	ObjectRetain(key)
	return key
}

func (key *RSAPrivateKey) Init(dict map[string]interface{}) *RSAPrivateKey {
	if key.BasePrivateKey.Init(dict) != nil {
		// lazy load
		key._rsaPrivateKey = nil
		key._data = nil
		key.setPublicKey(nil)
	}
	return key
}

func (key *RSAPrivateKey) Release() int {
	cnt := key.BasePrivateKey.Release()
	if cnt == 0 {
		// this object is going to be destroyed,
		// release children
		key.setPublicKey(nil)
	}
	return cnt
}

func (key *RSAPrivateKey) setPublicKey(pKey PublicKey) {
	if pKey != key._publicKey {
		ObjectRetain(pKey)
		ObjectRelease(key._publicKey)
		key._publicKey = pKey
	}
}

func (key *RSAPrivateKey) getPrivateKey() *rsa.PrivateKey {
	if key._rsaPrivateKey == nil {
		data := key.Get("data")
		if data == nil {
			// generate new key with size
			key._rsaPrivateKey, _ = key.generate(1024)
		} else {
			block, _ := pem.Decode(UTF8Encode(data.(string)))
			pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err !=  nil {
				panic(err)
			}
			key._rsaPrivateKey = pri.(*rsa.PrivateKey)
		}
	}
	return key._rsaPrivateKey
}

func (key *RSAPrivateKey) generate(bits int) (*rsa.PrivateKey, []byte) {
	pri, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	der := MarshalPKCS8PrivateKey(pri)
	block := &pem.Block{
		Type: "PRIVATE KEY",
		Bytes: der,
	}
	data := pem.EncodeToMemory(block)
	key.Set("data", UTF8Decode(data))
	key.Set("mode", "ECB")
	key.Set("padding", "PKCS1")
	key.Set("digest", "SHA256")
	return pri, der
}

func (key *RSAPrivateKey) getHash() crypto.Hash {
	return crypto.SHA256
}

//-------- ICryptographyKey

func (key *RSAPrivateKey) Data() []byte {
	if key._data == nil {
		// TODO: encode private key data to PKCS1
		pri := key.getPrivateKey()
		key._data = pri.D.Bytes()
	}
	return key._data
}

//-------- IPrivateKey(ISignKey)

func (key *RSAPrivateKey) Sign(data []byte) []byte {
	pri := key.getPrivateKey()
	h := key.getHash().New()
	h.Write(data)
	sum := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, pri, key.getHash(), sum)
	if err != nil {
		panic(err)
	}
	return sig
}

//-------- IDecryptKey

func (key *RSAPrivateKey) Decrypt(ciphertext []byte) []byte {
	pri := key.getPrivateKey()
	part := pri.N.BitLen() / 8
	chunks := BytesSplit(ciphertext, part)
	buffer := bytes.NewBufferString("")
	for _, line := range chunks {
		data, err := rsa.DecryptPKCS1v15(rand.Reader, pri, line)
		if err != nil {
			panic(err)
		}
		buffer.Write(data)
	}
	return buffer.Bytes()
}

//-------- IPrivateKey

func (key *RSAPrivateKey) PublicKey() PublicKey {
	if key._publicKey == nil {
		sKey := key.getPrivateKey()
		pKey := &sKey.PublicKey
		der, err := x509.MarshalPKIXPublicKey(pKey)
		if err != nil {
			panic(err)
		}
		block := &pem.Block{
			Type: "PUBLIC KEY",
			Bytes: der,
		}
		data := pem.EncodeToMemory(block)
		key.setPublicKey(PublicKeyParse(map[string]interface{}{
			"algorithm": RSA,
			"data": UTF8Decode(data),
			"mode": "ECB",
			"padding": "PKCS1",
			"digest": "SHA256",
		}))
	}
	return key._publicKey
}

//
//  PKCS8
//

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k
}
