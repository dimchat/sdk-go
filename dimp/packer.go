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
package dimp

import (
	. "github.com/dimchat/core-go/core"
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
)

type MessengerPacker struct {
	MessagePacker
}

func (packer *MessengerPacker) Init(messenger IMessenger) *MessengerPacker {
	if packer.MessagePacker.Init(messenger) != nil {
	}
	return packer
}

func (packer *MessengerPacker) Messenger() IMessenger {
	return packer.Transceiver().(IMessenger)
}

func (packer *MessengerPacker) Facebook() IFacebook {
	return packer.Messenger().Facebook()
}

func (packer *MessengerPacker) isWaiting(identifier ID) bool {
	if identifier.IsBroadcast() {
		// broadcast ID doesn't contain meta or visa
		return false
	}
	if identifier.IsGroup() {
		// if group is not broadcast ID, its meta should be exists
		return packer.Facebook().GetMeta(identifier) == nil
	}
	// if user is not broadcast ID, its visa key should be exists
	return packer.Facebook().GetPublicKeyForEncryption(identifier) == nil
}

func (packer *MessengerPacker) EncryptMessage(iMsg InstantMessage) SecureMessage {
	receiver := iMsg.Receiver()
	group := iMsg.Group()
	if packer.isWaiting(receiver) || (group != nil && packer.isWaiting(group)) {
		// NOTICE: the application will query visa automatically,
		//         save this message in a queue waiting sender's visa response
		packer.Messenger().SuspendInstantMessage(iMsg)
		return nil
	}
	// make sure visa.key exists before encrypting message
	return packer.MessagePacker.EncryptMessage(iMsg)
}

func (packer *MessengerPacker) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	facebook := packer.Facebook()
	sender := rMsg.Sender()
	// [Meta Protocol]
	meta := rMsg.Meta()
	if meta == nil {
		// get from local storage
		meta = facebook.GetMeta(sender)
	} else if !facebook.SaveMeta(meta, sender) {
		// failed to save meta attached to message
		meta = nil
	}
	if meta == nil {
		// NOTICE: the application will query meta automatically,
		//         save this message in a queue waiting sender's meta response
		packer.Messenger().SuspendReliableMessage(rMsg)
		return nil
	}
	// [Visa Protocol]
	visa := rMsg.Visa()
	if visa != nil {
		// check & save visa attached to message
		doc, ok := visa.(Document)
		if ok {
			facebook.SaveDocument(doc)
		}
	}
	// make sure meta exists before verifying message
	return packer.MessagePacker.VerifyMessage(rMsg)
}

func (packer *MessengerPacker) DecryptMessage(sMsg SecureMessage) InstantMessage {
	// check message delegate
	if sMsg.Delegate() == nil {
		sMsg.SetDelegate(packer.Messenger())
	}
	var trimmed SecureMessage
	receiver := sMsg.Receiver()
	user := packer.Messenger().SelectLocalUser(receiver)
	if user == nil {
		// current users not match
		trimmed = nil
	} else if receiver.IsGroup() {
		// trim group message
		trimmed = sMsg.Trim(user.ID())
	} else {
		trimmed = sMsg
	}
	if trimmed == nil {
		// not for you?
		panic("receiver error: " + receiver.String())
		return nil
	}
	// make sure private key (decrypt key) exists before decrypting message
	return packer.MessagePacker.DecryptMessage(sMsg)
}
