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
	"fmt"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

//
//  CPU
//
type ContentProcessor interface {

	SetMessenger(messenger *Messenger)

	Process(content Content, rMsg ReliableMessage) Content
}

//
//  CPU Factories
//
var contentProcessors = make(map[uint8]ContentProcessor)

func ContentProcessorRegister(msgType uint8, cpu ContentProcessor) {
	contentProcessors[msgType] = cpu
}
func ContentProcessorGet(content Content) ContentProcessor {
	return ContentProcessorGetByType(content.Type())
}
func ContentProcessorGetByType(msgType uint8) ContentProcessor {
	return contentProcessors[msgType]
}

/**
 *  Base Content Processor
 */
type BaseContentProcessor struct {
	ContentProcessor

	_messenger *Messenger
}

func (cpu *BaseContentProcessor) Init() *BaseContentProcessor {
	return cpu
}

func (cpu *BaseContentProcessor) SetMessenger(messenger *Messenger) {
	cpu._messenger = messenger
}

func (cpu *BaseContentProcessor) Messenger() *Messenger {
	return cpu._messenger
}

func (cpu *BaseContentProcessor) Facebook() *Facebook {
	return cpu.Messenger().Facebook()
}

func (cpu *BaseContentProcessor) Process(content Content, _ ReliableMessage) Content {
	text := fmt.Sprintf("Content (type: %d) not support yet!", content.Type())
	res := NewTextContent(text)
	// check group message
	group := content.Group()
	if group != nil {
		res.SetGroup(group)
	}
	return res
}

//
//  Register content processors
//
func BuildContentProcessors() {
	
}
