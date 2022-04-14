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
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/dimp"
)

var (
	StrDocCmdError = "Document command error."
	FmtDocNotFound = "Sorry, document not found for ID: %s"
	FmtDocNotAccepted = "Document not accepted: %s"
	FmtDocAccepted = "Document received: %s"
)

type DocumentCommandProcessor struct {
	MetaCommandProcessor
}

func NewDocumentCommandProcessor(facebook IFacebook, messenger IMessenger) ContentProcessor {
	cpu := new(DocumentCommandProcessor)
	cpu.Init(facebook, messenger)
	return cpu
}

//func (cpu *DocumentCommandProcessor) Init(facebook IFacebook, messenger IMessenger) ContentProcessor {
//	if cpu.MetaCommandProcessor.Init(facebook, messenger) != nil {
//	}
//	return cpu
//}

func (cpu *DocumentCommandProcessor) getDocument(identifier ID, docType string) []Content {
	facebook := cpu.Facebook()
	// query entity document for ID
	doc := facebook.GetDocument(identifier, docType)
	if doc == nil {
		text := fmt.Sprintf(FmtDocNotFound, identifier.String())
		return cpu.RespondText(text, nil)
	} else {
		meta := facebook.GetMeta(identifier)
		res := DocumentCommandRespond(identifier, meta, doc)
		return cpu.RespondContent(res)
	}
}

func (cpu *DocumentCommandProcessor) putDocument(identifier ID, meta Meta, doc Document) []Content {
	if doc.ID().Equal(identifier) == false {
		return cpu.RespondText(StrDocCmdError, nil)
	}
	facebook := cpu.Facebook()
	if !ValueIsNil(meta) {
		// received a meta for ID
		if facebook.SaveMeta(meta, identifier) == false {
			// meta not match
			text := fmt.Sprintf(FmtMetaNotAccepted, identifier.String())
			return cpu.RespondText(text, nil)
		}
	}
	// received a document for ID
	if facebook.SaveDocument(doc) {
		// document saved
		text := fmt.Sprintf(FmtDocAccepted, identifier.String())
		return cpu.RespondReceipt(text)
	} else {
		// save document failed
		text := fmt.Sprintf(FmtDocNotAccepted, identifier.String())
		return cpu.RespondText(text, nil)
	}
}

//-------- IContentProcessor

func (cpu *DocumentCommandProcessor) Process(content Content, _ ReliableMessage) []Content {
	cmd, _ := content.(DocumentCommand)
	identifier := cmd.ID()
	doc := cmd.Document()
	if identifier == nil {
		// error
		return cpu.RespondText(StrDocCmdError, nil)
	} else if doc == nil {
		// query document for ID
		docType, ok := cmd.Get("doc_type").(string)
		if !ok || docType == "" {
			docType = "*"  // ANY
		}
		return cpu.getDocument(identifier, docType)
	} else {
		// received a new document for ID
		return cpu.putDocument(identifier, cmd.Meta(), doc)
	}
}
