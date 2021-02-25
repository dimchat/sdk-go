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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

type MetaCommandProcessor struct {
	BaseCommandProcessor
}

func (cpu *MetaCommandProcessor) Init() *MetaCommandProcessor {
	if cpu.BaseCommandProcessor.Init() != nil {
	}
	return cpu
}

func (cpu *MetaCommandProcessor) getMeta(identifier ID) Content {
	// query meta for ID
	meta := cpu.Facebook().GetMeta(identifier)
	if meta == nil {
		// meta not found
		text := "Sorry, meta not found for ID: " + identifier.String()
		return NewTextContent(text)
	}
	// response
	return MetaCommandRespond(identifier, meta)
}

func (cpu *MetaCommandProcessor) putMeta(identifier ID, meta Meta) Content {
	// received a meta for ID
	if cpu.Facebook().SaveMeta(meta, identifier) == false {
		// save meta failed
		text := "Meta not accepted: " + identifier.String()
		return NewTextContent(text)
	}
	// response
	text := "Meta received: " + identifier.String()
	//return new(ReceiptCommand).InitWithText(text)
	return receipt(text)
}

func (cpu *MetaCommandProcessor) Execute(cmd Command, _ ReliableMessage) Content {
	mCmd, _ := cmd.(MetaCommand)
	identifier := mCmd.ID()
	meta := mCmd.Meta()
	if meta == nil {
		return cpu.getMeta(identifier)
	} else {
		return cpu.putMeta(identifier, meta)
	}
}
