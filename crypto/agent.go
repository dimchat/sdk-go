/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/ext"
	. "github.com/dimchat/mkm-go/protocol"
)

type VisaAgent interface {

	/**
	 *  Encrypt plaintext to ciphertexts with all visa keys
	 *
	 * @param plaintext - key data
	 * @param meta      - meta for public key
	 * @param documents - visa documents for public keys
	 * @return encrypted data with terminals
	 */
	EncryptBundle(plaintext []byte, meta Meta, documents []Document) EncryptedBundle

	/**
	 *  Get all verify keys from documents and meta
	 *
	 * @param meta      - meta for public key
	 * @param documents - visa documents for public keys
	 * @return verify keys
	 */
	GetVerifyKeys(meta Meta, documents []Document) []VerifyKey

	/**
	 *  Get all terminals from documents
	 *
	 * @param documents - visa documents
	 * @return terminals
	 */
	GetTerminals(documents []Document) []string
}

/**
 *  Base VisaAgent
 */
type DefaultVisaAgent struct {
	//VisaAgent
}

// protected
func (agent DefaultVisaAgent) GetVerifyKey(doc Document) VerifyKey {
	if visa, ok := doc.(Visa); ok {
		pKey := visa.PublicKey()
		if vKey, ok := pKey.(VerifyKey); ok {
			return vKey
		}
		//panic("invalid key type")
		return nil
	}
	// public key in user profile?
	key := doc.GetProperty("key")
	return ParsePublicKey(key)
}

// protected
func (agent DefaultVisaAgent) GetEncryptKey(doc Document) EncryptKey {
	if visa, ok := doc.(Visa); ok {
		pKey := visa.PublicKey()
		if pKey != nil {
			return pKey
		}
		//panic("public key not visa")
		return nil
	}
	key := doc.GetProperty("key")
	pubKey := ParsePublicKey(key)
	if pubKey == nil {
		// profile document?
		return nil
	} else if encKey, ok := pubKey.(EncryptKey); ok {
		return encKey
	}
	//panic("public key is not encrypt key")
	return nil
}

// protected
func (agent DefaultVisaAgent) GetTerminal(doc Document) string {
	terminal := doc.GetString("terminal", "")
	if terminal == "" {
		// get from document ID
		helper := GetGeneralAccountHelper()
		did := helper.GetDocumentID(doc.Map())
		if did != nil {
			terminal = did.Terminal()
		} else {
			//panic("terminal not found")
			// TODO: get from property?
		}
	}
	return terminal
}

// Override
func (agent DefaultVisaAgent) EncryptBundle(plaintext []byte, meta Meta, documents []Document) EncryptedBundle {
	// NOTICE: meta.key will never be changed, so use Visa.key to encrypt message
	//         is a better way
	bundle := NewEncryptedBundle()
	var terminal string
	var pubKey EncryptKey
	var ciphertext []byte
	//
	//  1. encrypt with visa keys
	//
	for _, doc := range documents {
		// encrypt by public key
		pubKey = agent.GetEncryptKey(doc)
		if pubKey == nil {
			continue
		}
		// get visa.terminal
		terminal = agent.GetTerminal(doc)
		if terminal == "" {
			terminal = "*"
		}
		if bundle.Contains(terminal) {
			//panic("duplicate terminal detected")
			continue
		}
		ciphertext = pubKey.Encrypt(plaintext, nil)
		bundle.Set(terminal, ciphertext)
	}
	if bundle.IsEmpty() {
		//
		//  2. encrypt with meta key
		//
		metaKey := meta.PublicKey()
		if encKey, ok := metaKey.(EncryptKey); ok {
			//terminal = "*"
			ciphertext = encKey.Encrypt(plaintext, nil)
			bundle.Set("*", ciphertext)
		}
	}
	// OK
	return bundle
}

// Override
func (agent DefaultVisaAgent) GetVerifyKeys(meta Meta, documents []Document) []VerifyKey {
	keys := make([]VerifyKey, 0)
	var pubKey VerifyKey
	// the sender may use communication key to sign message.data,
	// try to verify it with visa.key first;
	for _, doc := range documents {
		pubKey = agent.GetVerifyKey(doc)
		if pubKey == nil {
			//panic("verify key not found")
			continue
		}
		keys = append(keys, pubKey)
	}
	// the sender may use identity key to sign message.data,
	// try to verify it with meta.key too.
	pubKey = meta.PublicKey()
	if pubKey != nil {
		keys = append(keys, pubKey)
	} else {
		//panic("failed to get meta key")
	}
	// OK
	return keys
}

// Override
func (agent DefaultVisaAgent) GetTerminals(documents []Document) []string {
	devices := make([]string, 0)
	var terminal string
	for _, doc := range documents {
		terminal = agent.GetTerminal(doc)
		if terminal == "" {
			terminal = "*"
		}
		devices = append(devices, terminal)
	}
	return devices
}

var sharedVisaAgent VisaAgent = &DefaultVisaAgent{}

func SetVisaAgent(agent VisaAgent) {
	sharedVisaAgent = agent
}

func GetVisaAgent() VisaAgent {
	return sharedVisaAgent
}
