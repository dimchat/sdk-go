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
)

type MessengerProcessor struct {
	MessageProcessor
}

func (processor *MessengerProcessor) Init(messenger IMessenger) *MessengerProcessor {
	if processor.MessageProcessor.Init(messenger) != nil {
	}
	return processor
}

func (processor *MessengerProcessor) Messenger() IMessenger {
	return processor.Transceiver().(IMessenger)
}

func (processor *MessengerProcessor) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) InstantMessage {
	res := processor.MessageProcessor.ProcessInstantMessage(iMsg, rMsg)
	if processor.Messenger().SaveMessage(iMsg) {
		return res
	}
	// error
	return nil
}

func (processor *MessengerProcessor) ProcessContent(content Content, rMsg ReliableMessage) Content {
	// TODO: override to check group
	cpu := ContentProcessorGet(content)
	if cpu == nil {
		cpu = ContentProcessorGetByType(0)  // unknown
	}
	cpu.SetMessenger(processor.Messenger())
	return cpu.Process(content, rMsg)
	// TODO: override to filter the response
}
