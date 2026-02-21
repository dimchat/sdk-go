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
	. "github.com/dimchat/core-go/format"
	. "github.com/dimchat/core-go/msg"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/crypto"
)

type PlainMessagePacker struct {
	//InstantMessagePacker

	// protected
	Transceiver InstantMessageDelegate
}

// Override
func (packer PlainMessagePacker) EncryptMessage(iMsg InstantMessage, password SymmetricKey, members []ID) SecureMessage {
	// TODO: check attachment for File/Image/Audio/Video message content
	//      (do it by application)
	transceiver := packer.Transceiver
	if transceiver == nil {
		//panic("instant message delegate not found")
		return nil
	}

	//
	//  1. Serialize 'message.content' to data (JsON / ProtoBuf / ...)
	//
	body := transceiver.SerializeContent(iMsg.Content(), password, iMsg)
	if len(body) == 0 {
		panic("fail to serialize content")
		return nil
	}

	//
	//  2. Encrypt content data to 'message.data' with symmetric key
	//
	ciphertext := transceiver.EncryptContent(body, password, iMsg)
	if len(ciphertext) == 0 {
		//panic("fail to encrypt content with key")
		return nil
	}

	//
	//  3. Encode 'message.data' to String (Base64)
	//
	var encodedData TransportableData
	if IsBroadcastMessage(iMsg) {
		// broadcast message content will not be encrypted (just encoded to JSON),
		// so no need to encode to Base64 here
		encodedData = NewPlainDataWithBytes(ciphertext) // UTF8Decode(ciphertext);
	} else {
		// message content had been encrypted by a symmetric key,
		// so the data should be encoded here (with algorithm 'base64' as default).
		encodedData = NewBase64DataWithBytes(ciphertext)
	}
	if encodedData == nil || encodedData.IsEmpty() {
		//panic("fail to encode content data")
		return nil
	}

	//
	//  4. Serialize message key to data (JsON / ProtoBuf / ...)
	//
	pwd := transceiver.SerializeKey(password, iMsg)
	// NOTICE:
	//    if the key is reused, iMsg must be updated with key digest.
	info := iMsg.CopyMap(false)

	// replace 'content' with encrypted and encoded 'data'
	delete(info, "content")
	info["data"] = encodedData.Serialize()

	// check serialized key data,
	// if key data is null here, build the secure message directly.
	if len(pwd) == 0 {
		// A) broadcast message has no key
		// B) reused key
		return ParseSecureMessage(info)
	}
	// encrypt and encode key

	if members == nil {
		// personal message
		receiver := iMsg.Receiver()
		members = []ID{
			receiver,
		}
	}

	bundleMap := make(map[ID]EncryptedBundle, len(members))
	var bundle EncryptedBundle
	for _, receiver := range members {
		//
		//  5. Encrypt key data to 'message.keys' with member's public keys
		//
		bundle = transceiver.EncryptKey(pwd, receiver, iMsg)
		if bundle == nil || bundle.IsEmpty() {
			// public key for member not found
			// TODO: suspend this message for waiting member's visa
			continue
		}
		bundleMap[receiver] = bundle
	}

	//
	//  6. Encode message key to String (Base64)
	//
	msgKeys := packer.EncodeKeys(bundleMap, iMsg)
	if len(msgKeys) == 0 {
		// public key for member(s) not found
		// TODO: suspend this message for waiting member's visa
		return nil
	}

	// insert as 'keys'
	info["keys"] = msgKeys

	// OK, pack message
	return ParseSecureMessage(info)
}

// protected
func (packer PlainMessagePacker) EncodeKeys(bundleMap map[ID]EncryptedBundle, iMsg InstantMessage) StringKeyMap {
	transceiver := packer.Transceiver
	if transceiver == nil {
		//panic("instant message delegate not found")
		return nil
	}
	msgKeys := NewMap()
	var encodedKeys StringKeyMap
	for receiver, bundle := range bundleMap {
		encodedKeys = transceiver.EncodeKey(bundle, receiver, iMsg)
		if len(encodedKeys) == 0 {
			//panic("fail to encode key data")
			continue
		}
		// insert to 'message.keys' with ID + terminal
		for k, v := range encodedKeys {
			msgKeys[k] = v
		}
	}
	// TODO: put key digest
	return msgKeys
}
