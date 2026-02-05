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
	. "github.com/dimchat/sdk-go/dimp/crypto"
	. "github.com/dimchat/sdk-go/dimp/mkm"
	. "github.com/dimchat/sdk-go/dimp/msg"
)

type ITransceiver interface {
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
 *  Message Transceiver
 *  <p>
 *      Converting message format between PlainMessage and NetworkMessage
 *  </p>
 */
type Transceiver struct {
	//ITransceiver

	// protected
	_delegate   EntityDelegate // facebook
	_compressor Compressor
}

func (transceiver *Transceiver) Init(facebook EntityDelegate) ITransceiver {
	transceiver._delegate = facebook
	transceiver._compressor = CreateCompressor()
	return transceiver
}

func (transceiver *Transceiver) EntityDelegate() EntityDelegate {
	return transceiver._delegate
}

func (transceiver *Transceiver) Compressor() Compressor {
	return transceiver._compressor
}

// Override
func (transceiver *Transceiver) SerializeMessage(rMsg ReliableMessage) []byte {
	compressor := transceiver.Compressor()
	info := rMsg.Map()
	return compressor.CompressReliableMessage(info)
}

// Override
func (transceiver *Transceiver) DeserializeMessage(data []byte) ReliableMessage {
	compressor := transceiver.Compressor()
	info := compressor.ExtractReliableMessage(data)
	return ParseReliableMessage(info)
}

//-------- InstantMessageDelegate

// Override
func (transceiver *Transceiver) SerializeContent(content Content, password SymmetricKey, _ InstantMessage) []byte {
	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         before serialize content, this job should be done in subclass
	compressor := transceiver.Compressor()
	info := content.Map()
	dict := password.Map()
	return compressor.CompressContent(info, dict)
}

// Override
func (transceiver *Transceiver) EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte {
	// store 'IV' in iMsg for AES decryption
	return password.Encrypt(data, iMsg.Map())
}

/*/
// Override
func (transceiver *Transceiver) EncodeData(data []byte, iMsg InstantMessage) interface{} {
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
func (transceiver *Transceiver) SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte {
	if IsBroadcastMessage(iMsg) {
		// broadcast message has no key
		return nil
	}
	compressor := transceiver.Compressor()
	dict := password.Map()
	return compressor.CompressSymmetricKey(dict)
}

// Override
func (transceiver *Transceiver) EncryptKey(data []byte, receiver ID, _ InstantMessage) EncryptedBundle {
	// TODO: make sure the receiver's public key exists
	facebook := transceiver.EntityDelegate()
	contact := facebook.GetUser(receiver)
	if contact == nil {
		//panic("failed to encrypt message key for contact")
		return nil
	}
	// encrypt with public key of the receiver (or group member)
	return contact.EncryptBundle(data)
}

// Override
func (transceiver *Transceiver) EncodeKey(bundle EncryptedBundle, receiver ID, _ InstantMessage) StringKeyMap {
	// message key had been encrypted by a public key,
	// so the data should be encoded here (with algorithm 'base64' as default).
	return bundle.Encode(receiver)
	// TODO: check for wildcard
}

//-------- ISecureMessageDelegate

// Override
func (transceiver *Transceiver) DecodeKey(msgKeys StringKeyMap, receiver ID, _ SecureMessage) EncryptedBundle {
	facebook := transceiver.EntityDelegate()
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
func (transceiver *Transceiver) DecryptKey(bundle EncryptedBundle, receiver ID, _ SecureMessage) []byte {
	// NOTICE: the receiver must be a member ID
	//         if it's a group message
	facebook := transceiver.EntityDelegate()
	user := facebook.GetUser(receiver)
	if user == nil {
		//panic("failed to decrypt key")
		return nil
	}
	// decrypt with private key of the receiver (or group member)
	return user.DecryptBundle(bundle)
}

// Override
func (transceiver *Transceiver) DeserializeKey(key []byte, _ SecureMessage) SymmetricKey {
	if key == nil || len(key) == 0 {
		//panic("reused key? get it from local cache")
		return nil
	}
	compressor := transceiver.Compressor()
	info := compressor.ExtractSymmetricKey(key)
	return ParseSymmetricKey(info)
}

/*/
// Override
func (transceiver *Transceiver) DecodeData(data interface{}, sMsg SecureMessage) []byte {
	if IsBroadcastMessage(sMsg) {
		// broadcast message content will not be encrypted (just encoded to JSON),
		// so return the string data directly
		if text, ok := data.(string); ok {
			return UTF8Encode(text)
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
func (transceiver *Transceiver) DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte {
	// check 'IV' in sMsg for AES decryption
	return password.Decrypt(data, sMsg.Map())
}

// Override
func (transceiver *Transceiver) DeserializeContent(data []byte, password SymmetricKey, _ SecureMessage) Content {
	compressor := transceiver.Compressor()
	dict := password.Map()
	info := compressor.ExtractContent(data, dict)
	return ParseContent(info)
	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         after deserialize content, this job should be done in subclass
}

// Override
func (transceiver *Transceiver) SignData(data []byte, sMsg SecureMessage) []byte {
	facebook := transceiver.EntityDelegate()
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
func (transceiver *Transceiver) EncodeSignature(signature []byte, _ SecureMessage) interface{} {
	ted := NewBase64DataWithBytes(signature)
	return ted.Serialize()
}
/*/

//-------- IReliableMessageDelegate

/*/
// Override
func (transceiver *Transceiver) DecodeSignature(signature interface{}, _ ReliableMessage) []byte {
	ted := ParseTransportableData(signature)
	return ted.Bytes()
}
/*/

// Override
func (transceiver *Transceiver) VerifyDataSignature(data []byte, signature []byte, rMsg ReliableMessage) bool {
	facebook := transceiver.EntityDelegate()
	sender := rMsg.Sender()
	contact := facebook.GetUser(sender)
	if contact == nil {
		//panic("failed to verify message signature for contact")
		return false
	}
	return contact.Verify(data, signature)
}
