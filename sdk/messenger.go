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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/sdk-go/dimp/core"
)

type IMessenger interface {
	ITransceiver
	Packer
	Processor

	//
	//  Interfaces for Cipher Key
	//
	GetEncryptKey(iMsg InstantMessage) SymmetricKey
	GetDecryptKey(sMsg SecureMessage) SymmetricKey
	CacheDecryptKey(key SymmetricKey, sMsg SecureMessage)
}

type Messenger struct {
	//IMessenger
	Transceiver

	Tee MessageTee
}

//func (messenger *Messenger) Self() IMessenger {
//	return messenger.Tee.Messenger()
//}

func (messenger *Messenger) Packer() Packer {
	return messenger.Tee.Packer()
}

func (messenger *Messenger) Processor() Processor {
	return messenger.Tee.Processor()
}

func (messenger *Messenger) CipherKeyDelegate() CipherKeyDelegate {
	return messenger.Tee.CipherKeyDelegate()
}

// Override
func (messenger *Messenger) GetEncryptKey(iMsg InstantMessage) SymmetricKey {
	sender := iMsg.Sender()
	target := CipherKeyDestinationForMessage(iMsg)
	db := messenger.CipherKeyDelegate()
	return db.GetCipherKey(sender, target, true)
}

// Override
func (messenger *Messenger) GetDecryptKey(sMsg SecureMessage) SymmetricKey {
	sender := sMsg.Sender()
	target := CipherKeyDestinationForMessage(sMsg)
	db := messenger.CipherKeyDelegate()
	return db.GetCipherKey(sender, target, false)
}

// Override
func (messenger *Messenger) CacheDecryptKey(key SymmetricKey, sMsg SecureMessage) {
	sender := sMsg.Sender()
	target := CipherKeyDestinationForMessage(sMsg)
	db := messenger.CipherKeyDelegate()
	db.CacheCipherKey(sender, target, key)
}

//-------- IPacker

// Override
func (messenger *Messenger) EncryptMessage(iMsg InstantMessage) SecureMessage {
	packer := messenger.Packer()
	return packer.EncryptMessage(iMsg)
}

// Override
func (messenger *Messenger) SignMessage(sMsg SecureMessage) ReliableMessage {
	packer := messenger.Packer()
	return packer.SignMessage(sMsg)
}

/*/
// Override
func (messenger *Messenger) SerializeMessage(rMsg ReliableMessage) []byte {
	packer := messenger.Packer()
	return packer.SerializeMessage(rMsg)
}

// Override
func (messenger *Messenger) DeserializeMessage(data []byte) ReliableMessage {
	packer := messenger.Packer()
	return packer.DeserializeMessage(data)
}
/*/

// Override
func (messenger *Messenger) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	packer := messenger.Packer()
	return packer.VerifyMessage(rMsg)
}

// Override
func (messenger *Messenger) DecryptMessage(sMsg SecureMessage) InstantMessage {
	packer := messenger.Packer()
	return packer.DecryptMessage(sMsg)
}

//-------- IProcessor

// Override
func (messenger *Messenger) ProcessPackage(data []byte) [][]byte {
	processor := messenger.Processor()
	return processor.ProcessPackage(data)
}

// Override
func (messenger *Messenger) ProcessReliableMessage(rMsg ReliableMessage) []ReliableMessage {
	processor := messenger.Processor()
	return processor.ProcessReliableMessage(rMsg)
}

// Override
func (messenger *Messenger) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) []SecureMessage {
	processor := messenger.Processor()
	return processor.ProcessSecureMessage(sMsg, rMsg)
}

// Override
func (messenger *Messenger) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) []InstantMessage {
	processor := messenger.Processor()
	return processor.ProcessInstantMessage(iMsg, rMsg)
}

// Override
func (messenger *Messenger) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	processor := messenger.Processor()
	return processor.ProcessContent(content, rMsg)
}
