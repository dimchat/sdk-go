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
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/sdk-go/dimp"
)

var (
	StrMetaCmdError = "Meta command error."
	FmtMetaNotFound = "Sorry, meta not found for ID: %s"
	FmtMetaNotAccepted = "Meta not accepted: %s"
	FmtMetaAccepted = "Meta received: %s"
)

type MetaCommandProcessor struct {
	CommandProcessor
}

func NewMetaCommandProcessor(facebook IFacebook, messenger IMessenger) *MetaCommandProcessor {
	cpu := new(MetaCommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

func (cpu *MetaCommandProcessor) getMeta(identifier ID) []Content {
	// query meta for ID
	meta := cpu.Facebook().GetMeta(identifier)
	if meta == nil {
		text := fmt.Sprintf(FmtMetaNotFound, identifier.String())
		return cpu.RespondText(text, nil)
	} else {
		res := MetaCommandRespond(identifier, meta)
		return cpu.RespondContent(res)
	}
}

func (cpu *MetaCommandProcessor) putMeta(identifier ID, meta Meta) []Content {
	// received a meta for ID
	if cpu.Facebook().SaveMeta(meta, identifier) {
		// meta saved
		text := fmt.Sprintf(FmtMetaAccepted, identifier.String())
		return cpu.RespondReceipt(text)
	} else {
		// save meta failed
		text := fmt.Sprintf(FmtMetaNotAccepted, identifier.String())
		return cpu.RespondText(text, nil)
	}
}

func (cpu *MetaCommandProcessor) Execute(cmd Command, _ ReliableMessage) []Content {
	mCmd, _ := cmd.(MetaCommand)
	identifier := mCmd.ID()
	meta := mCmd.Meta()
	if identifier == nil {
		// error
		return cpu.RespondText(StrMetaCmdError, cmd.Group())
	} else if meta == nil {
		// query meta for ID
		return cpu.getMeta(identifier)
	} else {
		// received a meta for ID
		return cpu.putMeta(identifier, meta)
	}
}
