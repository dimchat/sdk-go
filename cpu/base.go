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
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/sdk"
)

/**
 *  CPU - Content Processing Unit
 */
type BaseContentProcessor struct {
	//ContentProcessor
	TwinsHelper
}

// Override
func (cpu *BaseContentProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	return cpu.RespondReceipt("Content not support.", rMsg.Envelope(), content, StringKeyMap{
		"template": "Content (type: ${type}) not support yet!",
		"replacements": StringKeyMap{
			"type": content.Type(),
		},
	})
}

//
//  Convenient responding
//

// protected
func (cpu *BaseContentProcessor) RespondReceipt(text string, head Envelope, body Content, extra StringKeyMap) []Content {
	// create base receipt command with text & original envelope
	res := createReceipt(text, head, body, extra)
	return []Content{res}
}

/**
 *  Create receipt command with text, original envelope, serial number &amp; group
 *
 * @param text     - text message
 * @param head     - original envelope
 * @param body     - original content
 * @param extra    - extra info
 * @return receipt command
 */
func createReceipt(text string, head Envelope, body Content, extra StringKeyMap) ReceiptCommand {
	// create base receipt command with text, original envelope, serial number & group ID
	res := NewReceiptCommand(text, head, body)
	if body != nil {
		// check group
		group := body.Group()
		if group != nil {
			res.SetGroup(group)
		}
	}
	// add extra key-value
	if extra != nil {
		for key, value := range extra {
			res.Set(key, value)
		}
	}
	return res
}

/**
 *  CPU - Command Processing Unit
 */
type BaseCommandProcessor struct {
	BaseContentProcessor
}

// Override
func (cpu *BaseCommandProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	command, ok := content.(Command)
	if !ok {
		//panic("command error")
		return nil
	}
	return cpu.RespondReceipt("Command not support.", rMsg.Envelope(), content, StringKeyMap{
		"template": "Command (name: ${command}) not support yet!",
		"replacements": StringKeyMap{
			"command": command.CMD(),
		},
	})
}
