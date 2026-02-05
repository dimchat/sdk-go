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
	. "github.com/dimchat/sdk-go/core"
	. "github.com/dimchat/sdk-go/msg"
)

type MessagePacker struct {
	//Packer
	TwinsHelper

	// protected
	InstantPacker  InstantMessagePacker
	SecurePacker   SecureMessagePacker
	ReliablePacker ReliableMessagePacker
}

func (packer *MessagePacker) Init(facebook IFacebook, messenger IMessenger) Packer {
	if packer.TwinsHelper.Init(facebook, messenger) != nil {
		packer.InstantPacker = CreateInstantMessagePacker(messenger)
		packer.SecurePacker = CreateSecureMessagePacker(messenger)
		packer.ReliablePacker = CreateReliableMessagePacker(messenger)
	}
	return packer
}

//
//  InstantMessage -> SecureMessage -> ReliableMessage -> Data
//

// Override
func (packer *MessagePacker) EncryptMessage(iMsg InstantMessage) SecureMessage {
	// TODO: check receiver before calling this, make sure the visa.key exists;
	//       otherwise, suspend this message for waiting receiver's visa/meta;
	//       if receiver is a group, query all members' visa too!
	facebook := packer.Facebook
	messenger := packer.Messenger

	var sMsg SecureMessage
	// NOTICE: before sending group message, you can decide whether expose the group ID
	//      (A) if you don't want to expose the group ID,
	//          you can split it to multi-messages before encrypting,
	//          replace the 'receiver' to each member and keep the group hidden in the content;
	//          in this situation, the packer will use the personal message key (user to user);
	//      (B) if the group ID is overt, no need to worry about the exposing,
	//          you can keep the 'receiver' being the group ID, or set the group ID as 'group'
	//          when splitting to multi-messages to let the remote packer knows it;
	//          in these situations, the local packer will use the group msg key (user to group)
	//          to encrypt the message, and the remote packer can get the overt group ID before
	//          decrypting to take the right message key.
	receiver := iMsg.Receiver()

	//
	//  1. get message key with direction (sender -> receiver) or (sender -> group)
	//
	password := messenger.GetEncryptKey(iMsg)
	if password == nil {
		//panic("failed to get msg key")
		return nil
	}

	//
	//  2. encrypt 'content' to 'data' for receiver/group members
	//
	if receiver.IsGroup() {
		// group message
		members := facebook.GetMembers(receiver)
		if members == nil || len(members) == 0 {
			//panic("group not ready")
			return nil
		}
		// a station will never send group message, so here must be a client;
		// the client messenger should check the group's meta & members before encrypting,
		// so we can trust that the group members MUST exist here.
		sMsg = packer.InstantPacker.EncryptMessage(iMsg, password, members)
	} else {
		// personal message (or split group message)
		sMsg = packer.InstantPacker.EncryptMessage(iMsg, password, nil)
	}
	if sMsg == nil {
		// public key for encryption not found
		// TODO: suspend this message for waiting receiver's meta
		return nil
	}

	// NOTICE: copy content type to envelope
	//         this help the intermediate nodes to recognize message type
	envelope := sMsg.Envelope()
	content := iMsg.Content()
	envelope.SetType(content.Type())

	// OK
	return sMsg
}

// Override
func (packer *MessagePacker) SignMessage(sMsg SecureMessage) ReliableMessage {
	// sign 'data' by sender
	delegate := packer.SecurePacker
	return delegate.SignMessage(sMsg)
}

/*/
// Override
func (packer *MessagePacker) SerializeMessage(rMsg ReliableMessage) []byte {
	dict := rMsg.Map()
	json := JSONEncodeMap(dict)
	return UTF8Encode(json)
}
/*/

//
//  Data -> ReliableMessage -> SecureMessage -> InstantMessage
//

/*/
// Override
func (packer *MessagePacker) DeserializeMessage(data []byte) ReliableMessage {
	json := UTF8Decode(data)
	dict := JSONDecodeMap(json)
	return ParseReliableMessage(dict)
}
/*/

// Override
func (packer *MessagePacker) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	// verify 'data' with 'signature'
	delegate := packer.ReliablePacker
	return delegate.VerifyMessage(rMsg)
}

// Override
func (packer *MessagePacker) DecryptMessage(sMsg SecureMessage) InstantMessage {
	// TODO: check receiver before calling this, make sure you are the receiver,
	//       or you are a member of the group when this is a group message,
	//       so that you will have a private key (decrypt key) to decrypt it.
	receiver := sMsg.Receiver()
	user := packer.SelectLocalUser(receiver)
	if user == nil {
		// not for you?
		sender := sMsg.Sender()
		panic("receiver error: " + receiver.String() + ", from " + sender.String())
		return nil
	}
	// decrypt 'data' to 'content'
	delegate := packer.SecurePacker
	return delegate.DecryptMessage(sMsg, user.ID())
	// TODO: check top-secret message
	//       (do it by application)
}
