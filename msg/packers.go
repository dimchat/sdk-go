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

	/**
	 *  Encrypt personal / group message
	 *  <p>
	 *      Replace 'content' field with encrypted 'data'
	 *  </p>
	 *
	 * @param iMsg     - plain message
	 * @param password - symmetric key
	 * @param members  - group members for group message; null for personal message
	 * @return SecureMessage object, null on visa not found
	 */
	EncryptMessage(iMsg InstantMessage, password SymmetricKey, members []ID) SecureMessage
}

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

	/**
	 *  Decrypt message, replace encrypted 'data' with 'content' field
	 *
	 * @param sMsg     - encrypted message
	 * @param receiver - actual receiver (local user)
	 * @return InstantMessage object
	 */
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

	/**
	 *  Sign message.data, add 'signature' field
	 *
	 * @param sMsg - encrypted message
	 * @return ReliableMessage object
	 */
	SignMessage(sMsg SecureMessage) ReliableMessage
}

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

	/**
	 *  Verify 'data' and 'signature' field with sender's public key
	 *
	 * @param rMsg - received message
	 * @return SecureMessage object
	 */
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
		Transceiver: messenger,
	}
}

// Override
func (msgPackerFactory) CreateSecureMessagePacker(messenger SecureMessageDelegate) SecureMessagePacker {
	return &EncryptedMessagePacker{
		Transceiver: messenger,
	}
}

// Override
func (msgPackerFactory) CreateReliableMessagePacker(messenger ReliableMessageDelegate) ReliableMessagePacker {
	return &NetworkMessagePacker{
		Transceiver: messenger,
	}
}
