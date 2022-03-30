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
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

type IMessenger interface {
	ITransceiver
	CipherKeyDelegate
	Packer
	Processor

	CipherKeyDelegate() CipherKeyDelegate
	SetCipherKeyDelegate(delegate CipherKeyDelegate)

	Packer() Packer
	SetPacker(packer Packer)

	Processor() Processor
	SetProcessor(processor Processor)
}

type Messenger struct {
	Transceiver

	_keyCache CipherKeyDelegate
	_packer Packer
	_processor Processor
}

func (messenger *Messenger) Init() *Messenger {
	if messenger.Transceiver.Init() != nil {
		messenger._keyCache = nil
		messenger._packer = nil
		messenger._processor = nil
	}
	return messenger
}

func (messenger *Messenger) CipherKeyDelegate() CipherKeyDelegate {
	return messenger._keyCache
}
func (messenger *Messenger) SetCipherKeyDelegate(delegate CipherKeyDelegate) {
	messenger._keyCache = delegate
}

func (messenger *Messenger) Packer() Packer {
	return messenger._packer
}
func (messenger *Messenger) SetPacker(packer Packer) {
	messenger._packer = packer
}

func (messenger *Messenger) Processor() Processor {
	return messenger._processor
}
func (messenger *Messenger) SetProcessor(processor Processor) {
	messenger._processor = processor
}

//-------- ICipherKeyDelegate

func (messenger *Messenger) GetCipherKey(sender, receiver ID, generate bool) SymmetricKey {
	return messenger.CipherKeyDelegate().GetCipherKey(sender, receiver, generate)
}

func (messenger *Messenger) CacheCipherKey(sender, receiver ID, key SymmetricKey) {
	messenger.CipherKeyDelegate().CacheCipherKey(sender, receiver, key)
}

//-------- IPacker

func (messenger *Messenger) GetOvertGroup(content Content) ID {
	return messenger.Packer().GetOvertGroup(content)
}

func (messenger *Messenger) EncryptMessage(iMsg InstantMessage) SecureMessage {
	return messenger.Packer().EncryptMessage(iMsg)
}

func (messenger *Messenger) SignMessage(sMsg SecureMessage) ReliableMessage {
	return messenger.Packer().SignMessage(sMsg)
}

func (messenger *Messenger) SerializeMessage(rMsg ReliableMessage) []byte {
	return messenger.Packer().SerializeMessage(rMsg)
}

func (messenger *Messenger) DeserializeMessage(data []byte) ReliableMessage {
	return messenger.Packer().DeserializeMessage(data)
}

func (messenger *Messenger) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	return messenger.Packer().VerifyMessage(rMsg)
}

func (messenger *Messenger) DecryptMessage(sMsg SecureMessage) InstantMessage {
	return messenger.Packer().DecryptMessage(sMsg)
}

//-------- IProcessor

func (messenger *Messenger) ProcessPackage(data []byte) [][]byte {
	return messenger.Processor().ProcessPackage(data)
}

func (messenger *Messenger) ProcessReliableMessage(rMsg ReliableMessage) []ReliableMessage {
	return messenger.Processor().ProcessReliableMessage(rMsg)
}

func (messenger *Messenger) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) []SecureMessage {
	return messenger.Processor().ProcessSecureMessage(sMsg, rMsg)
}

func (messenger *Messenger) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) []InstantMessage {
	return messenger.Processor().ProcessInstantMessage(iMsg, rMsg)
}

func (messenger *Messenger) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	return messenger.Processor().ProcessContent(content, rMsg)
}

//-------- ISecureMessageDelegate (ITransceiver)

func (messenger *Messenger) DeserializeKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) SymmetricKey {
	if key == nil {
		// get key from cache
		return messenger.GetCipherKey(sender, receiver, false)
	} else {
		// call super()
		return messenger.Transceiver.DeserializeKey(key, sender, receiver, sMsg)
	}
}

func (messenger *Messenger) DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content {
	// call super()
	content := messenger.Transceiver.DeserializeContent(data, password, sMsg)
	if !messenger.IsBroadcast(sMsg) {
		// check and cache key for reuse
		group := messenger.GetOvertGroup(content)
		if group == nil {
			// personal message or (group) command
			// cache key with direction (sender -> receiver)
			messenger.CacheCipherKey(sMsg.Sender(), sMsg.Receiver(), password)
		} else {
			// group message (excludes group command)
			// cache the key with direction (sender -> group)
			messenger.CacheCipherKey(sMsg.Sender(), group, password)
		}
	}
	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         after deserialize content, this job should be do in subclass
	return content
}
