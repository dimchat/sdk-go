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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp/cpu"
)

type Messenger struct {
	Transceiver

	_delegate MessengerDelegate
	_dataSource MessengerDataSource

	_transmitter Transmitter

	facebook *Facebook
	packer *MessengerPacker
	processor *MessengerProcessor
	transmitter *MessengerTransmitter
}

func (messenger *Messenger) Init() *Messenger {
	if messenger.Transceiver.Init() != nil {
		messenger._delegate = nil
		messenger._dataSource = nil

		messenger._transmitter = nil

		messenger.facebook = nil
		messenger.packer = nil
		messenger.processor = nil
		messenger.transmitter = nil
	}
	return messenger
}

/**
 *  Delegate for Station
 *
 * @param delegate - message delegate
 */
func (messenger *Messenger) SetDelegate(delegate MessengerDelegate) {
	messenger._delegate = delegate
}
func (messenger *Messenger) Delegate() MessengerDelegate {
	return messenger._delegate
}

/**
 *  Delegate for Message Storage
 *
 * @param delegate - message data source
 */
func (messenger *Messenger) SetDataSource(delegate MessengerDataSource) {
	messenger._dataSource = delegate
}
func (messenger *Messenger) DataSource() MessengerDataSource {
	return messenger._dataSource
}

/**
 *  Delegate for getting entity
 *
 * @param barrack - facebook
 */
func (messenger *Messenger) SetEntityDelegate(barrack EntityDelegate) {
	messenger.Transceiver.SetEntityDelegate(barrack)
	facebook, ok := barrack.(*Facebook)
	if ok {
		messenger.facebook = facebook
	}
}
func (messenger *Messenger) EntityDelegate() EntityDelegate {
	delegate := messenger.Transceiver.EntityDelegate()
	if delegate == nil {
		delegate = messenger.Facebook()
		messenger.Transceiver.SetEntityDelegate(delegate)
	}
	return delegate
}
func (messenger *Messenger) Facebook() *Facebook {
	if messenger.facebook == nil {
		messenger.facebook = messenger.NewFacebook()
	}
	return messenger.facebook
}
func (messenger *Messenger) NewFacebook() *Facebook {
	return new(Facebook).Init()
}

/**
 *  Delegate for packing message
 *
 * @param packer - message packer
 */
func (messenger *Messenger) SetPacker(packer Packer) {
	messenger.Transceiver.SetPacker(packer)
	mp, ok := packer.(*MessengerPacker)
	if ok {
		messenger.packer = mp
	}
}
func (messenger *Messenger) Packer() Packer {
	packer := messenger.Transceiver.Packer()
	if packer == nil {
		packer = messenger.MessengerPacker()
		messenger.Transceiver.SetPacker(packer)
	}
	return packer
}
func (messenger *Messenger) MessengerPacker() *MessengerPacker {
	if messenger.packer == nil {
		messenger.packer = messenger.NewPacker()
	}
	return messenger.packer
}
func (messenger *Messenger) NewPacker() *MessengerPacker {
	return new(MessengerPacker).Init(messenger)
}

/**
 *  Delegate for processing message
 *
 * @param processor - message processor
 */
func (messenger *Messenger) SetProcessor(processor Processor) {
	messenger.Transceiver.SetProcessor(processor)
	mp, ok := processor.(*MessengerProcessor)
	if ok {
		messenger.processor = mp
	}
}
func (messenger *Messenger) Processor() Processor {
	processor := messenger.Transceiver.Processor()
	if processor == nil {
		processor = messenger.MessengerProcessor()
		messenger.Transceiver.SetProcessor(processor)
	}
	return processor
}
func (messenger *Messenger) MessengerProcessor() *MessengerProcessor {
	if messenger.processor == nil {
		messenger.processor = messenger.NewProcessor()
	}
	return messenger.processor
}
func (messenger *Messenger) NewProcessor() *MessengerProcessor {
	return new(MessengerProcessor).Init(messenger)
}

/**
 *  Delegate for transmitting message
 *
 * @param transmitter - message transmitter
 */
func (messenger *Messenger) SetTransmitter(transmitter Transmitter) {
	messenger._transmitter = transmitter
	mt, ok := transmitter.(*MessengerTransmitter)
	if ok {
		messenger._transmitter = mt
	}
}
func (messenger *Messenger) Transmitter() Transmitter {
	if messenger._transmitter == nil {
		messenger._transmitter = messenger.MessengerTransmitter()
	}
	return messenger._transmitter
}
func (messenger *Messenger) MessengerTransmitter() *MessengerTransmitter {
	if messenger.transmitter == nil {
		messenger.transmitter = messenger.NewTransmitter()
	}
	return messenger.transmitter
}
func (messenger *Messenger) NewTransmitter() *MessengerTransmitter {
	return new(MessengerTransmitter).Init(messenger)
}

func (messenger *Messenger) getFileContentProcessor() *FileContentProcessor {
	processor := ContentProcessorGetByType(FILE)
	processor.SetMessenger(messenger)
	return processor.(*FileContentProcessor)
}

//-------- InstantMessageDelegate

func (messenger *Messenger) SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte {
	// check attachment for File/Image/Audio/Video message content
	file, ok := content.(*FileContent)
	if ok {
		fpu := messenger.getFileContentProcessor()
		fpu.UploadFileContent(file, password, iMsg)
	}
	return messenger.Transceiver.SerializeContent(content, password, iMsg)
}

func (messenger *Messenger) EncryptKey(data []byte, receiver ID, iMsg InstantMessage) []byte {
	key := messenger.Facebook().GetPublicKeyForEncryption(receiver)
	if key == nil {
		// save this message in a queue waiting receiver's meta/document response
		messenger.SuspendInstantMessage(iMsg)
		return nil
	}
	return messenger.Transceiver.EncryptKey(data, receiver, iMsg)
}

//-------- SecureMessageDelegate

func (messenger *Messenger) DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content {
	content := messenger.Transceiver.DeserializeContent(data, password, sMsg)
	// check attachment for File/Image/Audio/Video message content
	file, ok := content.(*FileContent)
	if ok {
		fpu := messenger.getFileContentProcessor()
		fpu.DownloadFileContent(file, password, sMsg)
	}
	return content
}

//
//  Interfaces for transmitting Message
//
func (messenger *Messenger) SendContent(sender ID, receiver ID, content Content, callback MessengerCallback, priority int) bool {
	if sender == nil {
		// Application Layer should make sure user is already login before it send message to server.
		// Application layer should put message into queue so that it will send automatically after user login
		user := messenger.Facebook().GetCurrentUser()
		sender = user.ID()
	}
	return messenger.Transmitter().SendContent(sender, receiver, content, callback, priority)
}

func (messenger *Messenger) SendInstantMessage(iMsg InstantMessage, callback MessengerCallback, priority int) bool {
	return messenger.Transmitter().SendInstantMessage(iMsg, callback, priority)
}

func (messenger *Messenger) SendReliableMessage(rMsg ReliableMessage, callback MessengerCallback, priority int) bool {
	return messenger.SendReliableMessage(rMsg, callback, priority)
}

//
//  Interfaces for Station
//
func (messenger *Messenger) UploadData(data []byte, iMsg InstantMessage) string {
	return messenger.Delegate().UploadData(data, iMsg)
}

func (messenger *Messenger) DownloadData(url string, iMsg InstantMessage) []byte {
	return messenger.Delegate().DownloadData(url, iMsg)
}

func (messenger *Messenger) SendPackage(data []byte, handler MessengerCompletionHandler, priority int) bool {
	return messenger.Delegate().SendPackage(data, handler, priority)
}

//
//  Interfaces for Message Storage
//
func (messenger *Messenger) SaveMessage(iMsg InstantMessage) bool {
	return messenger.DataSource().SaveMessage(iMsg)
}

func (messenger *Messenger) SuspendInstantMessage(iMsg InstantMessage) {
	messenger.DataSource().SuspendInstantMessage(iMsg)
}

func (messenger *Messenger) SuspendReliableMessage(rMsg ReliableMessage) {
	messenger.DataSource().SuspendReliableMessage(rMsg)
}
