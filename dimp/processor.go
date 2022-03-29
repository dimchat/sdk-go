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
	"time"
)

type MessageProcessor struct {
	MessengerHelper
	IProcessor

	_factory IProcessorFactory
}

func (processor *MessageProcessor) Init(facebook IFacebook, messenger IMessenger, factory IProcessorFactory) *MessageProcessor {
	if processor.MessengerHelper.Init(facebook, messenger) != nil {
		processor._factory = factory
	}
	return processor
}

func (processor *MessageProcessor) Factory() IProcessorFactory {
	return processor._factory
}

//-------- Processor Factory

func (processor *MessageProcessor) GetProcessor(content Content) IContentProcessor {
	return processor.Factory().GetProcessor(content)
}
func (processor *MessageProcessor) GetContentProcessor(msgType uint8) IContentProcessor {
	return processor.Factory().GetContentProcessor(msgType)
}
func (processor *MessageProcessor) GetCommandProcessor(msgType uint8, cmdName string) ICommandProcessor {
	return processor.Factory().GetCommandProcessor(msgType, cmdName)
}

//-------- IProcessor

func (processor *MessageProcessor) ProcessPackage(data []byte) [][]byte {
	messenger := processor.Messenger()
	// 1. deserialize message
	rMsg := messenger.DeserializeMessage(data)
	if rMsg == nil {
		// no valid message received
		return nil
	}
	// 2. process message
	responses := messenger.ProcessReliableMessage(rMsg)
	if responses == nil || len(responses) == 0 {
		// nothing to respond
		return nil
	}
	// 3. serialize message
	packages := make([][]byte, 0, len(responses))
	for _, item := range responses {
		pack := messenger.SerializeMessage(item)
		if pack != nil {
			packages = append(packages, pack)
		}
	}
	return packages
}

func (processor *MessageProcessor) ProcessReliableMessage(rMsg ReliableMessage) []ReliableMessage {
	// NOTICE: override to check broadcast message before calling it
	messenger := processor.Messenger()

	// 1. verify message
	sMsg := messenger.VerifyMessage(rMsg)
	if sMsg == nil {
		// waiting for sender's meta if not exists
		return nil
	}
	// 2. process message
	responses := messenger.ProcessSecureMessage(sMsg, rMsg)
	if responses == nil || len(responses) == 0 {
		// nothing to respond
		return nil
	}
	// 3. sign message
	messages := make([]ReliableMessage, 0, len(responses))
	for _, item := range responses {
		msg := messenger.SignMessage(item)
		if msg != nil {
			messages = append(messages, msg)
		}
	}
	return messages
	// NOTICE: override to deliver to the receiver when catch exception "receiver error ..."
}

func (processor *MessageProcessor) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) []SecureMessage {
	messenger := processor.Messenger()
	// 1. decrypt message
	iMsg := messenger.DecryptMessage(sMsg)
	if iMsg == nil {
		// cannot decrypt this message, not for you?
		// delivering message to other receiver?
		return nil
	}
	// 2. process message
	responses := messenger.ProcessInstantMessage(iMsg, rMsg)
	if responses == nil || len(responses) == 0 {
		// nothing to respond
		return nil
	}
	// 3. encrypt message
	messages := make([]SecureMessage, 0, len(responses))
	for _, item := range responses {
		msg := messenger.EncryptMessage(item)
		if msg != nil {
			messages = append(messages, msg)
		}
	}
	return messages
}

func (processor *MessageProcessor) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) []InstantMessage {
	messenger := processor.Messenger()
	// 1. process content
	responses := messenger.ProcessContent(iMsg.Content(), rMsg)
	if responses == nil || len(responses) == 0 {
		// nothing to respond
		return nil
	}

	// 2. select a local user to build message
	facebook := processor.Facebook()
	sender := iMsg.Sender()
	receiver := iMsg.Receiver()
	user := facebook.SelectLocalUser(receiver)

	// 3. pack messages
	messages := make([]InstantMessage, 0, len(responses))
	for _, item := range responses {
		env := EnvelopeCreate(user.ID(), sender, time.Time{})
		msg := InstantMessageCreate(env, item)
		messages = append(messages, msg)
	}
	return messages
}

func (processor *MessageProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	// TODO: override to check group
	cpu := processor.GetProcessor(content)
	return cpu.Process(content, rMsg)
	// TODO: override to filter the response
}
