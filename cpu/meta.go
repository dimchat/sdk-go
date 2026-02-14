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
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/core"
)

/**
 *  CPU for MetaCommand
 *  ~~~~~~~~~~~~~~~~~~~
 */
type MetaCommandProcessor struct {
	BaseCommandProcessor
}

func (cpu *MetaCommandProcessor) GetArchivist() Archivist {
	facebook := cpu.Facebook
	return facebook.GetArchivist()
}

// Override
func (cpu *MetaCommandProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	command, ok := content.(MetaCommand)
	if !ok {
		//panic("meta command error")
		return nil
	}
	did := command.ID()
	meta := command.Meta()
	if did == nil {
		// error
		return cpu.RespondReceipt("Meta command error.", rMsg.Envelope(), content, nil)
	} else if meta == nil {
		// query meta for ID
		return cpu.getMeta(did, rMsg.Envelope(), command)
	}
	// received a meta for ID
	return cpu.putMeta(meta, did, rMsg.Envelope(), command)
}

func (cpu *MetaCommandProcessor) getMeta(did ID, envelope Envelope, content MetaCommand) []Content {
	facebook := cpu.Facebook
	meta := facebook.GetMeta(did)
	if meta == nil {
		return cpu.RespondReceipt("Meta not found.", envelope, content, StringKeyMap{
			"template": "Meta not found: ${did}.",
			"replacements": StringKeyMap{
				"did": did.String(),
			},
		})
	}
	// meta got
	return cpu.respondMeta(did, meta, envelope.Sender())
}

// protected
func (cpu *MetaCommandProcessor) respondMeta(did ID, meta Meta, receiver ID) []Content {
	if receiver.Equal(did) {
		panic("cycled response: " + receiver.String())
	}
	// TODO: check response expired
	res := NewCommandForRespondMeta(did, meta)
	return []Content{res}
}

func (cpu *MetaCommandProcessor) putMeta(meta Meta, did ID, envelope Envelope, content MetaCommand) []Content {
	var errors []Content
	// 1. try to save meta
	errors = cpu.saveMeta(meta, did, envelope, content)
	if errors != nil {
		return errors
	}
	// 2. success
	return cpu.RespondReceipt("Meta received.", envelope, content, StringKeyMap{
		"template": "Meta received: ${did}.",
		"replacements": StringKeyMap{
			"did": did.String(),
		},
	})
}

// protected
func (cpu *MetaCommandProcessor) saveMeta(meta Meta, did ID, envelope Envelope, content MetaCommand) []Content {
	archivist := cpu.GetArchivist()
	// check meta
	if !cpu.checkMeta(meta, did) {
		// meta invalid
		return cpu.RespondReceipt("Meta not valid.", envelope, content, StringKeyMap{
			"template": "Meta not valid: ${did}.",
			"replacements": StringKeyMap{
				"did": did.String(),
			},
		})
	} else if !archivist.SaveMeta(meta, did) {
		// DB error?
		return cpu.RespondReceipt("Meta not accepted.", envelope, content, StringKeyMap{
			"template": "Meta not accepted: ${did}.",
			"replacements": StringKeyMap{
				"did": did.String(),
			},
		})
	}
	// meta saved, return no error
	return nil
}

// protected
func (cpu *MetaCommandProcessor) checkMeta(meta Meta, did ID) bool {
	if !meta.IsValid() {
		return false
	}
	old := did.Address()
	gen := GenerateAddress(meta, old.Network())
	return old.Equal(gen)
}
