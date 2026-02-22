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
	"fmt"

	. "github.com/dimchat/core-go/format"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/crypto"
)

type EncryptedMessagePacker struct {
	//SecureMessagePacker

	// protected
	Transformer SecureMessageDelegate
}

// protected
func (packer *EncryptedMessagePacker) DecodeKey(sMsg SecureMessage, receiver ID) EncryptedBundle {
	msgKeys := sMsg.EncryptedKeys()
	if msgKeys == nil {
		// get from 'key'
		base64 := sMsg.Get("key")
		if base64 == nil {
			// broadcast message?
			// reused key?
			return nil
		}
		msgKeys = NewMap()
		msgKeys[receiver.String()] = base64
	}
	transformer := packer.Transformer
	if transformer == nil {
		//panic("secure message delegate not found")
		return nil
	}
	return transformer.DecodeKey(msgKeys, receiver, sMsg)
}

// Override
func (packer *EncryptedMessagePacker) DecryptMessage(sMsg SecureMessage, receiver ID) InstantMessage {
	transformer := packer.Transformer
	if transformer == nil {
		//panic("secure message delegate not found")
		return nil
	}

	var pwd []byte // serialized symmetric key data

	//
	//  1. Decode 'message.key' to encrypted symmetric key data
	//
	bundle := packer.DecodeKey(sMsg, receiver)
	if bundle == nil || bundle.IsEmpty() {
		// broadcast message?
		// reused key?
		pwd = nil
	} else {
		//
		//  2. Decrypt 'message.key' with receiver's private key
		//
		pwd = transformer.DecryptKey(bundle, receiver, sMsg)
		if len(pwd) == 0 {
			// A: my visa updated but the sender doesn't got the new one;
			// B: key data error.
			panic(fmt.Sprintf("failed to decrypt message key: %s, %s => %s",
				bundle.String(), sMsg.Sender().String(), receiver.String()))
			// TODO: check whether my visa key is changed, push new visa to this contact
		}
	}

	//
	//  3. Deserialize message key from data (JsON / ProtoBuf / ...)
	//     (if key is empty, means it should be reused, get it from key cache)
	//
	password := transformer.DeserializeKey(pwd, sMsg)
	if password == nil {
		// A: key data is empty, and cipher key not found from local storage;
		// B: key data error.
		panic(fmt.Sprintf("failed to get message key: %d byte(s), %s => %s",
			len(pwd), sMsg.Sender().String(), receiver.String()))
		// TODO: ask the sender to send again (with new message key)
	}

	//
	//  4. Decode 'message.data' to encrypted content data
	//
	ciphertext := sMsg.Data()
	if ciphertext == nil || ciphertext.IsEmpty() {
		//panic("failed to decode message data")
		return nil
	}

	//
	//  5. Decrypt 'message.data' with symmetric key
	//
	body := transformer.DecryptContent(ciphertext.Bytes(), password, sMsg)
	if len(body) == 0 {
		// A: password is a reused key loaded from local storage, but it's expired;
		// B: key error.
		panic(fmt.Sprintf("failed to decrypt message data with key: %v, data length: %d, %s => %s",
			password.Map(), ciphertext.Size(), sMsg.Sender().String(), receiver.String()))
		// TODO: ask the sender to send again
	}

	//
	//  6. Deserialize message content from data (JsON / ProtoBuf / ...)
	//
	content := transformer.DeserializeContent(body, password, sMsg)
	if content == nil {
		//panic("failed to deserialize content")
		return nil
	}

	// TODO: check attachment for File/Image/Audio/Video message content
	//      if URL exists, means file data was uploaded to a CDN,
	//          1. save password as 'content.key';
	//          2. try to download file data from CDN;
	//          3. decrypt downloaded data with 'content.key'.
	//      (do it by application)

	// OK, pack message
	info := sMsg.CopyMap(false)
	delete(info, "key")
	delete(info, "keys")
	delete(info, "data")
	info["content"] = content.Map()
	return ParseInstantMessage(info)
}

// Override
func (packer *EncryptedMessagePacker) SignMessage(sMsg SecureMessage) ReliableMessage {
	transformer := packer.Transformer
	if transformer == nil {
		//panic("secure message delegate not found")
		return nil
	}

	//
	//  0. decode message data
	//
	ciphertext := sMsg.Data()
	if ciphertext == nil || ciphertext.IsEmpty() {
		//panic("failed to decrypt message data")
		return nil
	}

	//
	//  1. Sign 'message.data' with sender's private key
	//
	signature := transformer.SignData(ciphertext.Bytes(), sMsg)
	if len(signature) == 0 {
		//panic("failed to sign message")
		return nil
	}

	//
	//  2. Encode 'message.signature' to String (Base64)
	//
	base64 := NewBase64DataWithBytes(signature)
	if base64 == nil || base64.IsEmpty() {
		panic("failed to encode signature")
		return nil
	}

	// OK, pack message
	info := sMsg.CopyMap(false)
	info["signature"] = base64.Serialize()
	return ParseReliableMessage(info)
}
