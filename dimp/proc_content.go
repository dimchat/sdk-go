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
	"fmt"
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp/dkd"
	. "github.com/dimchat/sdk-go/dimp/protocol"
)

var (
	FmtContentNotSupport = "Content (type: %d) not support yet!"
)

/**
 *  CPU: Content Processing Unit
 */
type ContentProcessor interface {

	/**
	 *  Process message content
	 *
	 * @param content - content received
	 * @param message - reliable message
	 * @return contents responding to msg.sender
	 */
	Process(content Content, rMsg ReliableMessage) []Content
}

type BaseContentProcessor struct {
	TwinsHelper
}

func NewContentProcessor(facebook IFacebook, messenger IMessenger) ContentProcessor {
	cpu := new(BaseContentProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

//-------- IContentProcessor

func (cpu *BaseContentProcessor) Process(content Content, _ ReliableMessage) []Content {
	text := fmt.Sprintf(FmtContentNotSupport, content.Type())
	return cpu.RespondText(text, content.Group())
}

//
//  Convenient responding
//

func (cpu *BaseContentProcessor) RespondText(text string, group ID) []Content {
	res := NewTextContent(text)
	if group != nil {
		res.SetGroup(group)
	}
	return []Content{res}
}

func (cpu *BaseContentProcessor) RespondReceipt(text string) []Content {
	res := NewReceiptCommand(text, nil, 0, nil)
	return []Content{res}
}

func (cpu *BaseContentProcessor) RespondContent(content Content) []Content {
	if content == nil {
		return nil
	} else {
		return []Content{content}
	}
}
