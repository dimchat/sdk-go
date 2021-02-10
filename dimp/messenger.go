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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
)

type Messenger struct {
	MessageTransceiver

	_delegate MessengerDelegate
	_dataSource MessengerDataSource

	_transmitter Transmitter
}

func (messenger *Messenger) Init() *Messenger {
	if messenger.MessageTransceiver.Init() != nil {
		messenger._delegate = nil
		messenger._dataSource = nil

		messenger._transmitter = nil
	}
	return messenger
}

func (messenger *Messenger) Facebook() *Facebook {
	barrack := messenger.EntityDelegate()
	facebook, ok := barrack.(*Facebook)
	if ok {
		return facebook
	} else {
		panic(facebook)
	}
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
 *  Delegate for transmitting message
 *
 * @param transmitter - message transmitter
 */
func (messenger *Messenger) SetTransmitter(transmitter Transmitter) {
	messenger._transmitter = transmitter
}
func (messenger *Messenger) Transmitter() Transmitter {
	return messenger._transmitter
}

//
//  Interfaces for transmitting Message
//
func (messenger *Messenger) SendContent(sender ID, receiver ID, content Content, callback MessengerCallback, priority int) bool {
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
