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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
)

type DocumentCommandProcessor struct {
	MetaCommandProcessor
}

func (cpu *DocumentCommandProcessor) Init() *DocumentCommandProcessor {
	if cpu.MetaCommandProcessor.Init() != nil {
	}
	return cpu
}

func (cpu *DocumentCommandProcessor) getDocument(identifier ID, docType string) Content {
	facebook := cpu.Facebook()
	// query entity document for ID
	doc := facebook.GetDocument(identifier, docType)
	if doc == nil {
		// document not found
		text := "Sorry, document not found for ID: " + identifier.String()
		return NewTextContent(text)
	}
	// response
	meta := facebook.GetMeta(identifier)
	return DocumentCommandRespond(identifier, meta, doc)
}

func (cpu *DocumentCommandProcessor) putDocument(identifier ID, meta Meta, doc Document) Content {
	facebook := cpu.Facebook()
	if meta != nil {
		// received a meta for ID
		if facebook.SaveMeta(meta, identifier) == false {
			// meta not match
			text := "Meta not accepted: " + identifier.String()
			return NewTextContent(text)
		}
	}
	// received a document for ID
	if facebook.SaveDocument(doc) == false {
		// save document failed
		text := "Document not accepted: " + identifier.String()
		return NewTextContent(text)
	}
	// response
	//text := "Document received: " + identifier.String()
	//return new(ReceiptCommand).InitWithText(text)
	return nil
}

func (cpu *DocumentCommandProcessor) Execute(cmd Command, _ ReliableMessage) Content {
	mCmd, _ := cmd.(*DocumentCommand)
	identifier := mCmd.ID()
	doc := mCmd.Document()
	if doc == nil {
		docType := cmd.Get("doc_type")
		if docType == nil {
			return cpu.getDocument(identifier, "*")
		} else {
			return cpu.getDocument(identifier, docType.(string))
		}
	} else {
		meta := mCmd.Meta()
		return cpu.putDocument(identifier, meta, doc)
	}
}
