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
	EncryptKey

	_publicKey *rsa.PublicKey
}

func NewRSAPublicKey(dict map[string]interface{}) *RSAPublicKey {
	return new(RSAPublicKey).Init(dict)
}

func (key *RSAPublicKey) Init(dict map[string]interface{}) *RSAPublicKey {
	if key.BasePublicKey.Init(dict) != nil {
		// lazy load
		key._publicKey = nil
	}
	return key
}

func (key *RSAPublicKey) getPublicKey() *rsa.PublicKey {
	if key._publicKey == nil {
		data := key.Get("data")
		if data == nil {
			//panic(key)
			return nil
		}
		block, _ := pem.Decode(UTF8Encode(data.(string)))
		if block == nil {
			panic(data)
			return nil
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err !=  nil {
			panic(err)
			return nil
		}
		key._publicKey, _ = pub.(*rsa.PublicKey)
	}
	return key._publicKey
}

func (key *RSAPublicKey) getHash() crypto.Hash {
	return crypto.SHA256
}

func (key *RSAPublicKey) Data() []byte {
	// TODO:
	return nil
}

func (key *RSAPublicKey) Verify(data []byte, signature []byte) bool {
	pub := key.getPublicKey()
	h := key.getHash().New()
	h.Write(data)
	sum := h.Sum(nil)
	err := rsa.VerifyPKCS1v15(pub, key.getHash(), sum, signature)
	return err == nil
}

func (key *RSAPublicKey) Encrypt(plaintext []byte) []byte {
	pub := key.getPublicKey()
	if pub == nil {
		//panic(key)
		return nil
	}
	part := pub.N.BitLen() / 8 - 11
	chunks := BytesSplit(plaintext, part)
	buffer := bytes.NewBufferString("")
	for _, line := range chunks {
		data, err := rsa.EncryptPKCS1v15(rand.Reader, pub, line)
		if err != nil {
			panic(err)
			return nil
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
	DecryptKey

	_privateKey *rsa.PrivateKey

	_publicKey PublicKey
}

func NewRSAPrivateKey(dict map[string]interface{}) *RSAPrivateKey {
	return new(RSAPrivateKey).Init(dict)
}

func (key *RSAPrivateKey) Init(dict map[string]interface{}) *RSAPrivateKey {
	if key.BasePrivateKey.Init(dict) != nil {
		// lazy load
		key._privateKey = nil
		key._publicKey = nil
	}
	return key
}

func (key *RSAPrivateKey) getPrivateKey() *rsa.PrivateKey {
	if key._privateKey == nil {
		data := key.Get("data")
		if data == nil {
			// generate new key with size
			key._privateKey = key.generate(1024)
		} else {
			block, _ := pem.Decode(UTF8Encode(data.(string)))
			if block == nil {
				panic(data)
				return nil
			}
			pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err !=  nil {
				panic(err)
				return nil
			}
			key._privateKey, _ = pri.(*rsa.PrivateKey)
		}
	}
	return key._privateKey
}

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

func (key *RSAPrivateKey) generate(bits int) *rsa.PrivateKey {
	pri, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
		return nil
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
	return pri
}

func (key *RSAPrivateKey) getHash() crypto.Hash {
	return crypto.SHA256
}

func (key *RSAPrivateKey) Data() []byte {
	// TODO:
	return nil
}

func (key *RSAPrivateKey) Sign(data []byte) []byte {
	pri := key.getPrivateKey()
	h := key.getHash().New()
	h.Write(data)
	sum := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, pri, key.getHash(), sum)
	if err != nil {
		panic(err)
		return nil
	}
	return sig
}

func (key *RSAPrivateKey) Decrypt(ciphertext []byte) []byte {
	pri := key.getPrivateKey()
	part := pri.N.BitLen() / 8
	chunks := BytesSplit(ciphertext, part)
	buffer := bytes.NewBufferString("")
	for _, line := range chunks {
		data, err := rsa.DecryptPKCS1v15(rand.Reader, pri, line)
		if err != nil {
			panic(err)
			return nil
		}
		buffer.Write(data)
	}
	return buffer.Bytes()
}

func (key *RSAPrivateKey) PublicKey() PublicKey {
	if key._publicKey == nil {
		sKey := key.getPrivateKey()
		pKey := &sKey.PublicKey
		der, err := x509.MarshalPKIXPublicKey(pKey)
		if err != nil {
			return nil
		}
		block := &pem.Block{
			Type: "PUBLIC KEY",
			Bytes: der,
		}
		data := pem.EncodeToMemory(block)
		key._publicKey = PublicKeyParse(map[string]interface{}{
			"algorithm": RSA,
			"data": UTF8Decode(data),
			"mode": "ECB",
			"padding": "PKCS1",
			"digest": "SHA256",
		})
	}
	return key._publicKey
}
