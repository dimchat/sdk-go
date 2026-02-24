/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
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
package sdk

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimchat/sdk-go/mkm"
)

// Messenger defines the core interface for all message-related operations in the system
//
// Acts as the unified entry point for message processing, combining:
//   - Transformer: Message serialization/deserialization and format conversion
//   - Packer: Message encryption/signing/verification/decryption workflows
//   - Processor: Message processing at all format levels (binary to content)
//   - Cipher key management: Retrieval and caching of symmetric encryption keys
type Messenger interface {
	Transformer
	Packer
	Processor

	// -------------------------------------------------------------------------
	//  Cipher Key Management (Symmetric Key Handling)
	// -------------------------------------------------------------------------

	// GetEncryptKey retrieves a symmetric key for encrypting an InstantMessage
	//
	// Uses CipherKeyDelegate to get/cache direction-specific (sender→receiver) keys
	//
	// Parameters:
	//   - iMsg - Plaintext InstantMessage to be encrypted
	// Returns: SymmetricKey for message encryption (nil if key generation/retrieval fails)
	GetEncryptKey(iMsg InstantMessage) SymmetricKey

	// GetDecryptKey retrieves a symmetric key for decrypting a SecureMessage
	//
	// Extracts and decrypts the key from the message's encrypted key bundle
	//
	// Parameters:
	//   - sMsg - Encrypted SecureMessage to be decrypted
	// Returns: SymmetricKey for message decryption (nil if key extraction fails)
	GetDecryptKey(sMsg SecureMessage) SymmetricKey

	// CacheDecryptKey stores a decrypted symmetric key for future reuse
	//
	// Caches the key with direction-specific context from the SecureMessage
	//
	// Parameters:
	//   - key  - Decrypted symmetric key to cache (must not be nil)
	//   - sMsg - Encrypted SecureMessage associated with the key
	CacheDecryptKey(key SymmetricKey, sMsg SecureMessage)
}

// BaseMessenger is the base implementation of the Messenger interface
//
// Provides core dependencies and infrastructure for message processing operations
type BaseMessenger struct {
	//Messenger
	*MessageTransformer

	// CipherKeyDelegate manages caching/retrieval of symmetric encryption keys
	//
	// Backend implementation for GetEncryptKey/GetDecryptKey/CacheDecryptKey
	CipherKeyDelegate CipherKeyDelegate

	// Packer handles message encryption/signing/verification/decryption workflows
	//
	// Implements the Packer interface methods
	Packer Packer

	// Processor handles message processing at all format levels
	//
	// Implements the Processor interface methods
	Processor Processor
}

func NewBaseMessenger(facebook EntityDelegate, delegate CipherKeyDelegate) *BaseMessenger {
	return &BaseMessenger{
		MessageTransformer: NewMessageTransformer(facebook),
		CipherKeyDelegate:  delegate,
		Packer:             nil,
		Processor:          nil,
	}
}

//-------- SecureMessageDelegate

// Override
func (messenger *BaseMessenger) DeserializeKey(key []byte, sMsg SecureMessage) SymmetricKey {
	if len(key) == 0 {
		// get key from cache with direction: sender -> receiver(group)
		return messenger.GetDecryptKey(sMsg)
	}
	password := messenger.MessageTransformer.DeserializeKey(key, sMsg)
	// cache decrypt key when success
	if password != nil {
		// cache the key with direction: sender -> receiver(group)
		messenger.CacheDecryptKey(password, sMsg)
	}
	return password
}

//
//  Interfaces for Cipher Key
//

func (messenger *BaseMessenger) GetEncryptKey(iMsg InstantMessage) SymmetricKey {
	sender := iMsg.Sender()
	target := CipherKeyDestinationForMessage(iMsg)
	db := messenger.CipherKeyDelegate
	return db.GetCipherKey(sender, target, true)
}

func (messenger *BaseMessenger) GetDecryptKey(sMsg SecureMessage) SymmetricKey {
	sender := sMsg.Sender()
	target := CipherKeyDestinationForMessage(sMsg)
	db := messenger.CipherKeyDelegate
	return db.GetCipherKey(sender, target, false)
}

func (messenger *BaseMessenger) CacheDecryptKey(key SymmetricKey, sMsg SecureMessage) {
	sender := sMsg.Sender()
	target := CipherKeyDestinationForMessage(sMsg)
	db := messenger.CipherKeyDelegate
	db.CacheCipherKey(sender, target, key)
}

//-------- IPacker

// Override
func (messenger *BaseMessenger) EncryptMessage(iMsg InstantMessage) SecureMessage {
	packer := messenger.Packer
	return packer.EncryptMessage(iMsg)
}

// Override
func (messenger *BaseMessenger) SignMessage(sMsg SecureMessage) ReliableMessage {
	packer := messenger.Packer
	return packer.SignMessage(sMsg)
}

/*/
// Override
func (messenger *BaseMessenger) SerializeMessage(rMsg ReliableMessage) []byte {
	packer := messenger.Packer
	return packer.SerializeMessage(rMsg)
}

// Override
func (messenger *BaseMessenger) DeserializeMessage(data []byte) ReliableMessage {
	packer := messenger.Packer
	return packer.DeserializeMessage(data)
}
/*/

// Override
func (messenger *BaseMessenger) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	packer := messenger.Packer
	return packer.VerifyMessage(rMsg)
}

// Override
func (messenger *BaseMessenger) DecryptMessage(sMsg SecureMessage) InstantMessage {
	packer := messenger.Packer
	return packer.DecryptMessage(sMsg)
}

//-------- IProcessor

// Override
func (messenger *BaseMessenger) ProcessPackage(data []byte) [][]byte {
	processor := messenger.Processor
	return processor.ProcessPackage(data)
}

// Override
func (messenger *BaseMessenger) ProcessReliableMessage(rMsg ReliableMessage) []ReliableMessage {
	processor := messenger.Processor
	return processor.ProcessReliableMessage(rMsg)
}

// Override
func (messenger *BaseMessenger) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) []SecureMessage {
	processor := messenger.Processor
	return processor.ProcessSecureMessage(sMsg, rMsg)
}

// Override
func (messenger *BaseMessenger) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) []InstantMessage {
	processor := messenger.Processor
	return processor.ProcessInstantMessage(iMsg, rMsg)
}

// Override
func (messenger *BaseMessenger) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	processor := messenger.Processor
	return processor.ProcessContent(content, rMsg)
}
