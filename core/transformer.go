/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
 *
 *                                Written in 2020 by Moky <albert.moky@gmail.com>
 *
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
package sdk

import (
	. "github.com/dimchat/core-go/msg"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/crypto"
	. "github.com/dimchat/sdk-go/mkm"
	. "github.com/dimchat/sdk-go/msg"
)

type Transformer interface {
	InstantMessageDelegate
	SecureMessageDelegate
	ReliableMessageDelegate

	/**
	 *  Serialize network message
	 *
	 * @param rMsg - network message
	 * @return data package
	 */
	SerializeMessage(rMsg ReliableMessage) []byte

	/**
	 *  Deserialize network message
	 *
	 * @param data - data package
	 * @return network message
	 */
	DeserializeMessage(data []byte) ReliableMessage
}

/**
 *  Message Transformer
 *  <p>
 *      Converting message format between PlainMessage and NetworkMessage
 *  </p>
 */
type MessageTransformer struct {
	//Transformer

	// protected
	EntityDelegate EntityDelegate // facebook
	Compressor     Compressor
}

func NewMessageTransformer(facebook EntityDelegate) *MessageTransformer {
	return &MessageTransformer{
		EntityDelegate: facebook,
		Compressor:     CreateCompressor(),
	}
}

func (transformer *MessageTransformer) SerializeMessage(rMsg ReliableMessage) []byte {
	compressor := transformer.Compressor
	info := rMsg.Map()
	return compressor.CompressReliableMessage(info)
}

func (transformer *MessageTransformer) DeserializeMessage(data []byte) ReliableMessage {
	compressor := transformer.Compressor
	info := compressor.ExtractReliableMessage(data)
	return ParseReliableMessage(info)
}

//-------- InstantMessageDelegate

// Override
func (transformer *MessageTransformer) SerializeContent(content Content, password SymmetricKey, _ InstantMessage) []byte {
	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         before serialize content, this job should be done in subclass
	compressor := transformer.Compressor
	info := content.Map()
	dict := password.Map()
	return compressor.CompressContent(info, dict)
}

// Override
func (transformer *MessageTransformer) EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte {
	// store 'IV' in iMsg for AES decryption
	return password.Encrypt(data, iMsg.Map())
}

/*/
// Override
func (transformer *MessageTransformer) EncodeData(data []byte, iMsg InstantMessage) interface{} {
	if IsBroadcastMessage(iMsg) {
		// broadcast message content will not be encrypted (just encoded to JSON),
		// so no need to encode to Base64 here
		return UTF8Decode(data)
	}
	// message content had been encrypted by a symmetric key,
	// so the data should be encoded here (with algorithm 'base64' as default).
	ted := NewBase64DataWithBytes(data)
	return ted.Serialize()
}
/*/

// Override
func (transformer *MessageTransformer) SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte {
	if IsBroadcastMessage(iMsg) {
		// broadcast message has no key
		return nil
	}
	compressor := transformer.Compressor
	dict := password.Map()
	return compressor.CompressSymmetricKey(dict)
}

// Override
func (transformer *MessageTransformer) EncryptKey(data []byte, receiver ID, _ InstantMessage) EncryptedBundle {
	// TODO: make sure the receiver's public key exists
	facebook := transformer.EntityDelegate
	contact := facebook.GetUser(receiver)
	if contact == nil {
		//panic("failed to encrypt message key for contact")
		return nil
	}
	// encrypt with public key of the receiver (or group member)
	return contact.EncryptBundle(data)
}

// Override
func (transformer *MessageTransformer) EncodeKey(bundle EncryptedBundle, receiver ID, _ InstantMessage) StringKeyMap {
	// message key had been encrypted by a public key,
	// so the data should be encoded here (with algorithm 'base64' as default).
	return bundle.Encode(receiver)
	// TODO: check for wildcard
}

//-------- ISecureMessageDelegate

// Override
func (transformer *MessageTransformer) DecodeKey(msgKeys StringKeyMap, receiver ID, _ SecureMessage) EncryptedBundle {
	facebook := transformer.EntityDelegate
	user := facebook.GetUser(receiver)
	if user == nil {
		//panic("failed to decode key")
		return nil
	}
	// decode key bundle for all terminals
	terminals := user.Terminals()
	bundle := DecodeEncryptedBundle(msgKeys, receiver, terminals)
	if bundle == nil || bundle.IsEmpty() {
		// check for wildcard
		for _, target := range terminals {
			if target == "*" {
				//panic("failed to decode key")
				return nil
			}
		}
		// decode key bundle for "*"
		terminals = []string{
			"*",
		}
		bundle = DecodeEncryptedBundle(msgKeys, receiver, terminals)
		if bundle == nil || bundle.IsEmpty() {
			//panic("failed to decode key")
			return nil
		}
	}
	return bundle
}

// Override
func (transformer *MessageTransformer) DecryptKey(bundle EncryptedBundle, receiver ID, _ SecureMessage) []byte {
	// NOTICE: the receiver must be a member ID
	//         if it's a group message
	facebook := transformer.EntityDelegate
	user := facebook.GetUser(receiver)
	if user == nil {
		//panic("failed to decrypt key")
		return nil
	}
	// decrypt with private key of the receiver (or group member)
	return user.DecryptBundle(bundle)
}

// Override
func (transformer *MessageTransformer) DeserializeKey(key []byte, _ SecureMessage) SymmetricKey {
	if len(key) == 0 {
		//panic("reused key? get it from local cache")
		return nil
	}
	compressor := transformer.Compressor
	info := compressor.ExtractSymmetricKey(key)
	return ParseSymmetricKey(info)
}

/*/
// Override
func (transformer *MessageTransformer) DecodeData(data interface{}, sMsg SecureMessage) []byte {
	if IsBroadcastMessage(sMsg) {
		// broadcast message content will not be encrypted (just encoded to JSON),
		// so return the string data directly
		if str, ok := data.(string); ok {
			return UTF8Encode(str)
		}
		//panic("content data error")
		return nil
	}
	// message content had been encrypted by a symmetric key,
	// so the data should be encoded here (with algorithm 'base64' as default).
	ted := ParseTransportableData(data)
	return ted.Bytes()
}
/*/

// Override
func (transformer *MessageTransformer) DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte {
	// check 'IV' in sMsg for AES decryption
	return password.Decrypt(data, sMsg.Map())
}

// Override
func (transformer *MessageTransformer) DeserializeContent(data []byte, password SymmetricKey, _ SecureMessage) Content {
	compressor := transformer.Compressor
	dict := password.Map()
	info := compressor.ExtractContent(data, dict)
	return ParseContent(info)
	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         after deserialize content, this job should be done in subclass
}

// Override
func (transformer *MessageTransformer) SignData(data []byte, sMsg SecureMessage) []byte {
	facebook := transformer.EntityDelegate
	sender := sMsg.Sender()
	user := facebook.GetUser(sender)
	if user == nil {
		//panic("failed to sign message data for user")
		return nil
	}
	return user.Sign(data)
}

/*/
// Override
func (transformer *MessageTransformer) EncodeSignature(signature []byte, _ SecureMessage) interface{} {
	ted := NewBase64DataWithBytes(signature)
	return ted.Serialize()
}
/*/

//-------- IReliableMessageDelegate

/*/
// Override
func (transformer *MessageTransformer) DecodeSignature(signature interface{}, _ ReliableMessage) []byte {
	ted := ParseTransportableData(signature)
	return ted.Bytes()
}
/*/

// Override
func (transformer *MessageTransformer) VerifyDataSignature(data []byte, signature []byte, rMsg ReliableMessage) bool {
	facebook := transformer.EntityDelegate
	sender := rMsg.Sender()
	contact := facebook.GetUser(sender)
	if contact == nil {
		//panic("failed to verify message signature for contact")
		return false
	}
	return contact.Verify(data, signature)
}
