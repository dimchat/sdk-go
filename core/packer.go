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

import . "github.com/dimchat/dkd-go/protocol"

// Packer defines the interface for end-to-end message packing/unpacking workflows
//
// Orchestrates the complete message lifecycle between plaintext and network formats:
//
//	Outbound Flow (Sending Message):
//	    InstantMessage → [Encrypt] → SecureMessage → [Sign] → ReliableMessage → [Serialize] → Binary Data
//
//	Inbound Flow (Receiving Message):
//	    Binary Data → [Deserialize] → ReliableMessage → [Verify] → SecureMessage → [Decrypt] → InstantMessage
type Packer interface {

	// -------------------------------------------------------------------------
	//  Outbound Message Flow (Plain → Encrypted → Signed → Stream)
	// -------------------------------------------------------------------------

	// EncryptMessage converts a plaintext InstantMessage to an encrypted SecureMessage
	//
	// Encrypts message content with symmetric key and wraps key with receiver's public key
	//
	// Parameters:
	//   - iMsg - Plaintext message to encrypt
	// Returns: Encrypted SecureMessage (nil if encryption fails)
	EncryptMessage(iMsg InstantMessage) SecureMessage

	// SignMessage adds a digital signature to a SecureMessage to create a ReliableMessage
	//
	// Signs encrypted content with sender's private key to ensure authenticity
	//
	// Parameters:
	//   - sMsg - Encrypted SecureMessage to sign
	// Returns: Signed ReliableMessage (nil if signing fails)
	SignMessage(sMsg SecureMessage) ReliableMessage

	// SerializeMessage converts ReliableMessage to binary data (see Transformer)
	//
	// Parameters:
	//   - rMsg - Signed ReliableMessage to serialize
	// Returns: Binary data package for network transmission
	//SerializeMessage(rMsg ReliableMessage) []byte

	// -------------------------------------------------------------------------
	//  Inbound Message Flow (Stream → Signed → Encrypted → Plain)
	// -------------------------------------------------------------------------

	// DeserializeMessage parses binary data to ReliableMessage (see Transformer)
	//
	// Parameters:
	//   - data - Binary data package from network
	// Returns: Structured ReliableMessage
	//DeserializeMessage(data []byte) ReliableMessage

	// VerifyMessage validates the signature on a ReliableMessage to get a SecureMessage
	//
	// Verifies signature with sender's public key; returns nil if verification fails (untrusted message)
	//
	// Parameters:
	//   - rMsg - Received signed ReliableMessage to verify
	// Returns: Trusted SecureMessage (nil if signature is invalid)
	VerifyMessage(rMsg ReliableMessage) SecureMessage

	// DecryptMessage converts an encrypted SecureMessage back to plaintext InstantMessage
	//
	// Decrypts symmetric key with receiver's private key and content with symmetric key
	//
	// Parameters:
	//   - sMsg - Encrypted SecureMessage to decrypt
	// Returns: Plaintext InstantMessage (nil if decryption fails)
	DecryptMessage(sMsg SecureMessage) InstantMessage
}
