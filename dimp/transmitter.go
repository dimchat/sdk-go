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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Message Transmitter
 *  ~~~~~~~~~~~~~~~~~~~
 */
type Transmitter interface {

	/**
	 *  Send message content to receiver
	 *
	 * @param sender - sender ID
	 * @param receiver - receiver ID
	 * @param content - message content
	 * @param callback - if needs callback, set it here
	 * @return true on success
	 */
	SendContent(sender ID, receiver ID, content Content, callback MessengerCallback, priority int) bool

	/**
	 *  Send instant message (encrypt and sign) onto DIM network
	 *
	 * @param iMsg - instant message
	 * @param callback - if needs callback, set it here
	 * @return true on success
	 */
	SendInstantMessage(iMsg InstantMessage, callback MessengerCallback, priority int) bool

	SendReliableMessage(rMsg ReliableMessage, callback MessengerCallback, priority int) bool
}

type MessengerTransmitter struct {
	Transmitter

	_messenger *Messenger
}

func (transmitter *MessengerTransmitter) Init(messenger *Messenger) *MessengerTransmitter {
	transmitter._messenger = messenger
	return transmitter
}

func (transmitter *MessengerTransmitter) Messenger() *Messenger {
	return transmitter._messenger
}
