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
package dkd

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

// InstantMessagePacker defines the interface for packing (encrypting) InstantMessages into SecureMessages
//
//	Core responsibility: Convert plaintext instant messages to encrypted secure messages by:
//	    1. Encrypting message content with a symmetric key (PW) → "data" field
//	    2. Encrypting the symmetric key with receiver's public key (PK) → "key/keys" field
//	Supports both personal (1:1) and group (1:N) message encryption
type InstantMessagePacker interface {

	/*
	 *  Encrypt the Instant Message to Secure Message
	 *
	 *    +----------+      +----------+
	 *    | sender   |      | sender   |
	 *    | receiver |      | receiver |
	 *    | time     |  ->  | time     |
	 *    |          |      |          |
	 *    | content  |      | data     |  1. data = encrypt(content, PW)
	 *    +----------+      | key/keys |  2. key  = encrypt(PW, receiver.PK)
	 *                      +----------+
	 */

	// EncryptMessage converts a plaintext InstantMessage to an encrypted SecureMessage
	//
	// Replaces the plaintext "content" field with encrypted "data" and adds encrypted key material
	// Handles both personal messages (members = nil) and group messages (members = group member IDs)
	//
	// Parameters:
	//   - iMsg     - Plaintext instant message to encrypt
	//   - password - Symmetric key used to encrypt the message content
	//   - members  - List of group member IDs (nil for personal messages)
	// Returns: Encrypted SecureMessage (nil if receiver's visa/encryption key is not found)
	EncryptMessage(iMsg InstantMessage, password SymmetricKey, members []ID) SecureMessage
}

// SecureMessagePacker defines the interface for two core packing operations:
//  1. Unpacking (decrypting) SecureMessages back to InstantMessages
//  2. Packing (signing) SecureMessages into ReliableMessages
//
// Decryption Workflow:
//  1. Decrypt symmetric key (PW) from "key/keys" using receiver's private key (SK)
//  2. Decrypt message content from "data" using the symmetric key → restore "content" field
//
// Signing Workflow:
//  1. Generate digital signature for "data" field using sender's private key (SK)
//  2. Add signature to "signature" field → convert to ReliableMessage
type SecureMessagePacker interface {

	/*
	 *  Decrypt the Secure Message to Instant Message
	 *
	 *    +----------+      +----------+
	 *    | sender   |      | sender   |
	 *    | receiver |      | receiver |
	 *    | time     |  ->  | time     |
	 *    |          |      |          |  1. PW      = decrypt(key, receiver.SK)
	 *    | data     |      | content  |  2. content = decrypt(data, PW)
	 *    | key/keys |      +----------+
	 *    +----------+
	 */

	// DecryptMessage converts an encrypted SecureMessage back to a plaintext InstantMessage
	//
	// Replaces the encrypted "data" field with the original plaintext "content" field
	// Uses the receiver's private key to decrypt the symmetric key and message content
	//
	// Parameters:
	//   - sMsg     - Encrypted secure message to decrypt
	//   - receiver - Actual recipient ID (local user with access to private decryption keys)
	// Returns: Plaintext InstantMessage (nil if decryption fails)
	DecryptMessage(sMsg SecureMessage, receiver ID) InstantMessage

	/*
	 *  Sign the Secure Message to Reliable Message
	 *
	 *    +----------+      +----------+
	 *    | sender   |      | sender   |
	 *    | receiver |      | receiver |
	 *    | time     |  ->  | time     |
	 *    |          |      |          |
	 *    | data     |      | data     |
	 *    | key/keys |      | key/keys |
	 *    +----------+      | signature|  1. signature = sign(data, sender.SK)
	 *                      +----------+
	 */

	// SignMessage adds a digital signature to a SecureMessage to create a ReliableMessage
	//
	// Signs the encrypted "data" field with the sender's private key and adds to "signature" field
	// The signature proves message authenticity (only the sender could generate a valid signature)
	//
	// Parameters:
	//   - sMsg - Encrypted secure message to sign
	// Returns: Signed ReliableMessage (nil if signing fails)
	SignMessage(sMsg SecureMessage) ReliableMessage
}

// ReliableMessagePacker defines the interface for unpacking (verifying) ReliableMessages to SecureMessages
//
// Core responsibility: Validate the digital signature on a ReliableMessage to ensure authenticity:
//  1. Verify the "data" field against the "signature" field using the sender's public key (PK)
//  2. If verification succeeds, return the underlying SecureMessage (without signature validation risk)
type ReliableMessagePacker interface {

	/*
	 *  Verify the Reliable Message to Secure Message
	 *
	 *    +----------+      +----------+
	 *    | sender   |      | sender   |
	 *    | receiver |      | receiver |
	 *    | time     |  ->  | time     |
	 *    |          |      |          |
	 *    | data     |      | data     |  1. verify(data, signature, sender.PK)
	 *    | key/keys |      | key/keys |
	 *    | signature|      +----------+
	 *    +----------+
	 */

	// VerifyMessage validates a signed ReliableMessage and returns the underlying SecureMessage
	//
	// Verifies the "data" and "signature" fields using the sender's public key to ensure authenticity
	// If verification fails (invalid signature), returns nil to indicate untrusted message
	//
	// Parameters:
	//   - rMsg - Received signed reliable message to verify
	// Returns: SecureMessage (nil if signature verification fails)
	VerifyMessage(rMsg ReliableMessage) SecureMessage
}

//
//  Factories
//

func CreateInstantMessagePacker(messenger InstantMessageDelegate) InstantMessagePacker {
	factory := GetMessagePackerFactory()
	return factory.CreateInstantMessagePacker(messenger)
}

func CreateSecureMessagePacker(messenger SecureMessageDelegate) SecureMessagePacker {
	factory := GetMessagePackerFactory()
	return factory.CreateSecureMessagePacker(messenger)
}

func CreateReliableMessagePacker(messenger ReliableMessageDelegate) ReliableMessagePacker {
	factory := GetMessagePackerFactory()
	return factory.CreateReliableMessagePacker(messenger)
}

type MessagePackerFactory interface {
	CreateInstantMessagePacker(messenger InstantMessageDelegate) InstantMessagePacker
	CreateSecureMessagePacker(messenger SecureMessageDelegate) SecureMessagePacker
	CreateReliableMessagePacker(messenger ReliableMessageDelegate) ReliableMessagePacker
}

var sharedMessagePackerFactory MessagePackerFactory = &msgPackerFactory{}

func SetMessagePackerFactory(factory MessagePackerFactory) {
	sharedMessagePackerFactory = factory
}

func GetMessagePackerFactory() MessagePackerFactory {
	return sharedMessagePackerFactory
}

type msgPackerFactory struct {
	//MessagePackerFactory
}

// Override
func (msgPackerFactory) CreateInstantMessagePacker(messenger InstantMessageDelegate) InstantMessagePacker {
	return &PlainMessagePacker{
		Transformer: messenger,
	}
}

// Override
func (msgPackerFactory) CreateSecureMessagePacker(messenger SecureMessageDelegate) SecureMessagePacker {
	return &EncryptedMessagePacker{
		Transformer: messenger,
	}
}

// Override
func (msgPackerFactory) CreateReliableMessagePacker(messenger ReliableMessageDelegate) ReliableMessagePacker {
	return &NetworkMessagePacker{
		Transformer: messenger,
	}
}
