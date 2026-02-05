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
package cpu

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dkd"
	. "github.com/dimchat/sdk-go/sdk"
)

type ForwardContentProcessor struct {
	BaseContentProcessor
}

// Override
func (cpu *ForwardContentProcessor) ProcessContent(content Content, _ ReliableMessage) []Content {
	forwardContent, ok := content.(ForwardContent)
	if !ok {
		//panic("forward content error")
		return nil
	}
	secrets := forwardContent.SecretMessages()
	size := len(secrets)
	if size < 1 {
		size = 1
	}
	responses := make([]Content, size)
	// call messenger to process it
	messenger := cpu.Messenger
	var res Content
	var results []ReliableMessage
	for _, item := range secrets {
		results = messenger.ProcessReliableMessage(item)
		if results == nil {
			res = NewForwardMessages([]ReliableMessage{})
		} else if len(results) == 1 {
			res = NewForwardMessage(results[0])
		} else {
			res = NewForwardMessages(results)
		}
		responses = append(responses, res)
	}
	return responses
}

type ArrayContentProcessor struct {
	BaseContentProcessor
}

// Override
func (cpu *ArrayContentProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	arrayContent, ok := content.(ArrayContent)
	if !ok {
		//panic("array content error")
		return nil
	}
	array := arrayContent.Contents()
	size := len(array)
	if size < 1 {
		size = 1
	}
	responses := make([]Content, size)
	// call messenger to process it
	messenger := cpu.Messenger
	var res Content
	var results []Content
	for _, item := range array {
		results = messenger.ProcessContent(item, rMsg)
		if results == nil {
			res = NewArrayContent([]Content{})
		} else if len(results) == 1 {
			res = results[0]
		} else {
			res = NewArrayContent(results)
		}
		responses = append(responses, res)
	}
	return responses
}

//
//  Factories
//

func NewForwardContentProcessor(facebook IFacebook, messenger IMessenger) ContentProcessor {
	return &ForwardContentProcessor{
		BaseContentProcessor{
			TwinsHelper{
				Facebook:  facebook,
				Messenger: messenger,
			},
		},
	}
}

func NewArrayContentProcessor(facebook IFacebook, messenger IMessenger) ContentProcessor {
	return &ArrayContentProcessor{
		BaseContentProcessor{
			TwinsHelper{
				Facebook:  facebook,
				Messenger: messenger,
			},
		},
	}
}
