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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp/core"
	. "github.com/dimchat/sdk-go/dimp/dkd"
)

type MessageProcessor struct {
	//Processor
	TwinsHelper

	// protected
	Factory ContentProcessorFactory
}

func (processor *MessageProcessor) Init(facebook IFacebook, messenger IMessenger) Processor {
	if processor.TwinsHelper.Init(facebook, messenger) != nil {
		processor.Factory = CreateContentProcessorFactory(facebook, messenger)
	}
	return processor
}

// Override
func (processor *MessageProcessor) ProcessPackage(data []byte) [][]byte {
	messenger := processor.Messenger
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
	for _, res := range responses {
		pack := messenger.SerializeMessage(res)
		if pack == nil || len(pack) == 0 {
			// should not happen
			continue
		}
		packages = append(packages, pack)
	}
	return packages
}

// Override
func (processor *MessageProcessor) ProcessReliableMessage(rMsg ReliableMessage) []ReliableMessage {
	// TODO: override to check broadcast message before calling it
	messenger := processor.Messenger
	// 1. verify message
	sMsg := messenger.VerifyMessage(rMsg)
	if sMsg == nil {
		// TODO: suspend and waiting for sender's meta if not exists
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
	for _, res := range responses {
		msg := messenger.SignMessage(res)
		if msg == nil {
			// should not happen
			continue
		}
		messages = append(messages, msg)
	}
	return messages
	// TODO: override to deliver to the receiver when catch exception "receiver error ..."
}

// Override
func (processor *MessageProcessor) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) []SecureMessage {
	messenger := processor.Messenger
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
	for _, res := range responses {
		msg := messenger.EncryptMessage(res)
		if msg == nil {
			// receiver not ready?
			continue
		}
		messages = append(messages, msg)
	}
	return messages
}

// Override
func (processor *MessageProcessor) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) []InstantMessage {
	messenger := processor.Messenger
	// 1. process content
	responses := messenger.ProcessContent(iMsg.Content(), rMsg)
	if responses == nil || len(responses) == 0 {
		// nothing to respond
		return nil
	}
	// 2. select a local user to build message
	sender := iMsg.Sender()
	receiver := iMsg.Receiver()
	user := processor.SelectLocalUser(receiver)
	if user == nil {
		//panic("receiver error")
		return nil
	}
	// 3. pack messages
	messages := make([]InstantMessage, 0, len(responses))
	for _, res := range responses {
		env := CreateEnvelope(user.ID(), sender, nil)
		msg := CreateInstantMessage(env, res)
		messages = append(messages, msg)
	}
	return messages
}

// Override
func (processor *MessageProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	// TODO: override to check group
	factory := processor.Factory
	cpu := factory.GetContentProcessor(content)
	if cpu == nil {
		// default content processor
		cpu = factory.GetContentProcessorForType(ContentType.ANY)
		if cpu == nil {
			panic("failed to get default CPU")
			return nil
		}
	}
	return cpu.ProcessContent(content, rMsg)
	// TODO: override to filter the response
}

//
//  CPU Factory Helper
//

func CreateContentProcessorFactory(facebook IFacebook, messenger IMessenger) ContentProcessorFactory {
	helper := GetContentProcessorHelper()
	return helper.CreateContentProcessorFactory(facebook, messenger)
}

type ContentProcessorHelper interface {
	CreateContentProcessorFactory(facebook IFacebook, messenger IMessenger) ContentProcessorFactory
}

var sharedContentProcessorHelper ContentProcessorHelper = nil

func SetContentProcessorHelper(helper ContentProcessorHelper) {
	sharedContentProcessorHelper = helper
}

func GetContentProcessorHelper() ContentProcessorHelper {
	return sharedContentProcessorHelper
}
