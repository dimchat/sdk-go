/* license: https://mit-license.org
 *
 *  Dao-Ke-Dao: Universal Message Module
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
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/crypto"
)

// InstantMessageDelegate defines the interface for encrypting InstantMessages to SecureMessages
//
//	Implements the 6-step encryption workflow for message content and symmetric keys:
//	    1. Serialize content (JSON/Protobuf)
//	    2. Encrypt content with symmetric key
//	    3. Encode encrypted content (Base64)
//	    4. Serialize symmetric key (JSON/Protobuf)
//	    5. Encrypt key with receiver's public key (per terminal)
//	    6. Encode encrypted key data (Base64 and ID+terminal mapping)
//
//	Encryption Workflow Overview:
//	    InstantMessage → [Content Serialize → Encrypt → Key Serialize → Encrypt → Encode] → SecureMessage
//
//	Data Flow:
//	    content → encrypt(content, PW) → data;
//	    PW → encrypt(PW, receiver.PK) → key/keys
type InstantMessageDelegate interface {

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

	// -------------------------------------------------------------------------
	//  Content Encryption (Steps 1-3)
	// -------------------------------------------------------------------------

	// SerializeContent converts message content to raw byte data for encryption
	//
	// Step 1. Supports multiple serialization formats (JSON, Protobuf, etc.)
	//
	// Parameters:
	//   - content  - Message content to serialize (Text/File/Command/...)
	//   - password - Symmetric key (includes compression algorithm metadata)
	//   - iMsg     - Source InstantMessage context
	// Returns: Serialized byte array of the content
	SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte

	// EncryptContent encrypts serialized content data with a symmetric key
	//
	// Step 2. Generates the "data" field for the resulting SecureMessage
	//
	// Parameters:
	//   - data     - Serialized content bytes (from SerializeContent)
	//   - password - Symmetric encryption key (AES/DES/...)
	//   - iMsg     - Source InstantMessage context
	// Returns: Encrypted content bytes
	EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte

	// EncodeData base64-encodes encrypted content data to string
	//
	// Step 3. Converts binary encrypted data to string format for network transmission
	//
	// Parameters:
	//   - data - Encrypted content bytes (from EncryptContent)
	//   - iMsg - Source InstantMessage context
	// Returns: Base64-encoded string (interface{} for format flexibility)
	//EncodeData(data []byte, iMsg InstantMessage) interface{}

	// -------------------------------------------------------------------------
	//  Key Encryption (Steps 4-6)
	// -------------------------------------------------------------------------

	// SerializeKey converts symmetric key to raw byte data for encryption
	//
	// Step 4. Returns nil for reused keys (broadcast messages or cached keys)
	//
	// Parameters:
	//   - password - Symmetric key to serialize
	//   - iMsg     - Source InstantMessage context
	// Returns: Serialized key bytes (nil for reused/broadcast keys)
	SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte

	// EncryptKey encrypts serialized key data with the receiver's public key
	//
	// Step 5. Generates terminal-specific encrypted keys (EncryptedBundle) for the receiver
	//
	// Parameters:
	//   - data     - Serialized key bytes (from SerializeKey)
	//   - receiver - Actual message recipient (user/group member ID)
	//   - iMsg     - Source InstantMessage context
	// Returns: EncryptedBundle mapping ID.terminals to encrypted key data
	EncryptKey(data []byte, receiver ID, iMsg InstantMessage) EncryptedBundle

	// EncodeKey converts EncryptedBundle to StringKeyMap for network transmission
	//
	// Step 6. Structures key data with receiver ID + terminal mappings (Base64 encoded)
	//
	// Parameters:
	//   - bundle   - Terminal-specific encrypted key data (from EncryptKey)
	//   - receiver - Actual message recipient (user/group member ID)
	//   - iMsg     - Source InstantMessage context
	// Returns: Encoded StringKeyMap (compatible with "message.key/keys" field)
	EncodeKey(bundle EncryptedBundle, receiver ID, iMsg InstantMessage) StringKeyMap
}

// SecureMessageDelegate defines the interface for decrypting SecureMessages to InstantMessages
// and signing SecureMessages to ReliableMessages. Implements two core workflows:
//
//	Decryption Workflow (SecureMessage → InstantMessage):
//	    1. Decode encrypted key data (Base64)
//	    2. Decrypt key with receiver's private key
//	    3. Deserialize symmetric key (JSON/Protobuf)
//	    4. Decode encrypted content (Base64)
//	    5. Decrypt content with symmetric key
//	    6. Deserialize message content (JSON/Protobuf)
//
//	Signing Workflow Overview:
//	    SecureMessage → [Data Sign → Encode] → ReliableMessage
//
//	Decryption Workflow Overview:
//	    SecureMessage → [Key Decode → Decrypt → Deserialize → Content Decode → Decrypt → Deserialize] → InstantMessage
//
//	Data Flow:
//	    key/keys → decrypt(key, receiver.SK) → PW;
//	    data → decrypt(data, PW) → content
type SecureMessageDelegate interface {

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

	// -------------------------------------------------------------------------
	//  Key Decryption (Steps 1-3)
	// -------------------------------------------------------------------------

	// DecodeKey parses encrypted key data from "message.key/keys"
	//
	// Step 1. Converts encoded key data back to terminal-specific EncryptedBundle
	//
	// Parameters:
	//   - msgKeys  - Encoded key data (from "message.key/keys" field)
	//   - receiver - Actual message recipient (user/group member ID)
	//   - sMsg     - Source SecureMessage context
	// Returns: EncryptedBundle with terminal-specific encrypted key data
	DecodeKey(msgKeys StringKeyMap, receiver ID, sMsg SecureMessage) EncryptedBundle

	// DecryptKey decrypts terminal-specific key data with receiver's private key
	//
	// Step 2. Retrieves the serialized symmetric key from the EncryptedBundle
	//
	// Parameters:
	//   - bundle   - Terminal-specific encrypted key data (from DecodeKey)
	//   - receiver - Actual message recipient (user/group member ID)
	//   - sMsg     - Source SecureMessage context
	// Returns: Serialized symmetric key bytes (nil if decryption fails)
	DecryptKey(bundle EncryptedBundle, receiver ID, sMsg SecureMessage) []byte

	// DeserializeKey reconstructs symmetric key from serialized byte data
	//
	// Step 3. If key data is empty, retrieves reused key from key cache
	//
	// Parameters:
	//   - key  - Serialized key bytes (from DecryptKey, nil for reused keys)
	//   - sMsg - Source SecureMessage context
	// Returns: Reconstructed SymmetricKey (nil if deserialization fails)
	DeserializeKey(key []byte, sMsg SecureMessage) SymmetricKey

	// -------------------------------------------------------------------------
	//  Content Decryption (Steps 4-6)
	// -------------------------------------------------------------------------

	// DecodeData decodes base64-encoded encrypted content to bytes
	//
	// Step 4. Converts string-encoded content back to binary encrypted data
	//
	// Parameters:
	//   - data - Base64-encoded encrypted content string
	//   - sMsg - Source SecureMessage context
	// Returns: Binary encrypted content bytes (nil if decoding fails)
	//DecodeData(data interface{}, sMsg SecureMessage) []byte

	// DecryptContent decrypts encrypted content data with symmetric key
	//
	// Step 5. Retrieves the serialized message content from encrypted bytes
	//
	// Parameters:
	//   - data     - Encrypted content bytes (from DecodeData or direct from message)
	//   - password - Symmetric decryption key (from DeserializeKey)
	//   - sMsg     - Source SecureMessage context
	// Returns: Serialized content bytes (nil if decryption fails)
	DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte

	// DeserializeContent reconstructs message content from serialized bytes
	//
	// Step 6. Uses compression algorithm metadata from symmetric key for decompression
	//
	// Parameters:
	//   - data     - Serialized content bytes (from DecryptContent)
	//   - password - Symmetric key (includes compression algorithm info)
	//   - sMsg     - Source SecureMessage context
	// Returns: Reconstructed Content object (nil if deserialization fails)
	DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content

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

	// -------------------------------------------------------------------------
	//  Signature Generation (SecureMessage → ReliableMessage)
	// -------------------------------------------------------------------------

	// SignData generates a digital signature for encrypted content data
	//
	// Uses sender's private key to sign the "data" field (proves authenticity)
	//
	// Parameters:
	//   - data - Encrypted content bytes (from SecureMessage.data)
	//   - sMsg - Source SecureMessage context
	// Returns: Digital signature bytes (nil if signing fails)
	SignData(data []byte, sMsg SecureMessage) []byte

	// EncodeSignature base64-encodes signature to string
	//
	// Converts binary signature to string format for network transmission
	//
	// Parameters:
	//   - signature - Digital signature bytes (from SignData)
	//   - sMsg      - Source SecureMessage context
	// Returns: Base64-encoded signature string (interface{} for format flexibility)
	//EncodeSignature(signature []byte, sMsg SecureMessage) interface{}
}

// ReliableMessageDelegate defines the interface for verifying ReliableMessages to SecureMessages
//
//	Implements the signature verification workflow:
//	    1. Decode signature (Base64)
//	    2. Verify signature with sender's public key
//
//	Verification Workflow Overview:
//	    ReliableMessage → [Signature Decode → Verify] → SecureMessage
//
//	Validation Rule:
//	    verify(data, signature, sender.PK) returns true for authentic messages
type ReliableMessageDelegate interface {
	//SecureMessageDelegate

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

	// -------------------------------------------------------------------------
	//  Signature Verification (Steps 1-2)
	// -------------------------------------------------------------------------

	// DecodeSignature decodes base64-encoded signature to bytes
	//
	// Step 1. Converts string-encoded signature back to binary format
	//
	// Parameters:
	//   - signature - Base64-encoded signature string
	//   - rMsg      - Source ReliableMessage context
	// Returns: Binary signature bytes (nil if decoding fails)
	//DecodeSignature(signature interface{}, rMsg ReliableMessage) []byte

	// VerifyDataSignature validates the authenticity of encrypted content data
	//
	// Step 2. Uses sender's public key to verify the signature matches the encrypted data
	//
	// Parameters:
	//   - data      - Encrypted content bytes (from ReliableMessage.data)
	//   - signature - Digital signature bytes (from DecodeSignature)
	//   - rMsg      - Source ReliableMessage context
	// Returns: true if signature is valid (authentic), false otherwise
	VerifyDataSignature(data []byte, signature []byte, rMsg ReliableMessage) bool
}
