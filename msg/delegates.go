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

/**
 *  Instant Message Delegate
 */
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

	//
	//  Encrypt Content
	//

	/**
	 *  1. Serialize 'message.content' to data (JsON / ProtoBuf / ...)
	 *
	 * @param content  - message.content
	 * @param password - symmetric key (includes data compression algorithm)
	 * @param iMsg     - instant message object
	 * @return serialized content data
	 */
	SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte

	/**
	 *  2. Encrypt content data to 'message.data' with symmetric key
	 *
	 * @param data     - serialized data of message.content
	 * @param password - symmetric key
	 * @param iMsg     - instant message object
	 * @return encrypted message content data
	 */
	EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte

	/*
	 *  3. Encode 'message.data' to String (Base64)
	 *
	 * @param data - encrypted content data
	 * @param iMsg - instant message object
	 * @return String object
	 */
	//EncodeData(data []byte, iMsg InstantMessage) interface{}

	//
	//  Encrypt Key
	//

	/**
	 *  4. Serialize message key to data (JsON / ProtoBuf / ...)
	 *
	 * @param password - symmetric key
	 * @param iMsg     - instant message object
	 * @return serialized key data, null for reused (or broadcast message)
	 */
	SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte

	/**
	 *  5. Encrypt key data to 'message.key/keys' with receiver's public key
	 *
	 * @param data     - serialized data of symmetric key
	 * @param receiver - actual receiver (user, or group member)
	 * @param iMsg     - instant message object
	 * @return encrypted symmetric key data and targets (ID terminals)
	 */
	EncryptKey(data []byte, receiver ID, iMsg InstantMessage) EncryptedBundle

	/**
	 *  6. Encode 'message.key' to String (Base64)
	 *
	 * @param bundle   - encrypted symmetric key data and targets (ID terminals)
	 * @param receiver - actual receiver (user, or group member)
	 * @param iMsg     - instant message object
	 * @return encoded key data and targets (ID + terminals)
	 */
	EncodeKey(bundle EncryptedBundle, receiver ID, iMsg InstantMessage) StringKeyMap
}

/**
 *  Secure Message Delegate
 */
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

	//
	//  Decrypt Key
	//

	/**
	 *  1. Decode 'message.key' to encrypted symmetric key data
	 *
	 * @param msgKeys  - encoded key data and targets (ID + terminals)
	 * @param receiver - actual receiver (user, or group member)
	 * @param sMsg     - secure message object
	 * @return encrypted symmetric key data and targets (ID terminals)
	 */
	DecodeKey(msgKeys StringKeyMap, receiver ID, sMsg SecureMessage) EncryptedBundle

	/**
	 *  2. Decrypt 'message.key' with receiver's private key
	 *
	 * @param bundle   - encrypted symmetric key data and targets (ID terminals)
	 * @param receiver - actual receiver (user, or group member)
	 * @param sMsg     - secure message object
	 * @return serialized data of symmetric key
	 */
	DecryptKey(bundle EncryptedBundle, receiver ID, sMsg SecureMessage) []byte

	/**
	 *  3. Deserialize message key from data (JsON / ProtoBuf / ...)
	 *  <p>
	 *      (if key data is empty, means it should be reused, get it from key cache)
	 *  </p>
	 *
	 * @param key      - serialized key data, null for reused key
	 * @param sMsg     - secure message object
	 * @return symmetric key
	 */
	DeserializeKey(key []byte, sMsg SecureMessage) SymmetricKey

	//
	//  Decrypt Content
	//

	/*
	 *  4. Decode 'message.data' to encrypted content data
	 *
	 * @param data - base64 string object
	 * @param sMsg - secure message object
	 * @return encrypted content data
	 */
	//DecodeData(data interface{}, sMsg SecureMessage) []byte

	/**
	 *  5. Decrypt 'message.data' with symmetric key
	 *
	 * @param data     - encrypt content data
	 * @param password - symmetric key
	 * @param sMsg     - secure message object
	 * @return serialized message content
	 */
	DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte

	/**
	 *  6. Deserialize message content from data (JsON / ProtoBuf / ...)
	 *
	 * @param data     - serialized content data
	 * @param password - symmetric key (includes data compression algorithm)
	 * @param sMsg     - secure message object
	 * @return message content
	 */
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

	//
	//  Signature
	//

	/**
	 *  1. Sign 'message.data' with sender's private key
	 *
	 * @param data - encrypted message data
	 * @param sMsg - secure message object
	 * @return signature of encrypted message data
	 */
	SignData(data []byte, sMsg SecureMessage) []byte

	/*
	 *  2. Encode 'message.signature' to String (Base64)
	 *
	 * @param signature - signature of message.data
	 * @param sMsg      - secure message object
	 * @return String object
	 */
	//EncodeSignature(signature []byte, sMsg SecureMessage) interface{}
}

/**
 *  Reliable Message Delegate
 */
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

	/*
	 *  1. Decode 'message.signature' from String (Base64)
	 *
	 * @param signature - base64 string object
	 * @param rMsg      - reliable message
	 * @return signature data
	 */
	//DecodeSignature(signature interface{}, rMsg ReliableMessage) []byte

	/**
	 *  2. Verify the message data and signature with sender's public key
	 *
	 * @param data      - message content(encrypted) data
	 * @param signature - signature for message content(encrypted) data
	 * @param rMsg      - reliable message object
	 * @return YES on signature matched
	 */
	VerifyDataSignature(data []byte, signature []byte, rMsg ReliableMessage) bool
}
