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
	. "github.com/dimchat/mkm-go/ext"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  CPU for DocumentCommand
 *  ~~~~~~~~~~~~~~~~~~~~~~~
 */
type DocumentCommandProcessor struct {
	*MetaCommandProcessor
}

// Override
func (cpu *DocumentCommandProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	command, ok := content.(DocumentCommand)
	if !ok {
		//panic("document command error")
		return nil
	}
	did := command.ID()
	docs := command.Documents()
	if did == nil {
		// error
		return cpu.RespondReceipt("document command error.", rMsg.Envelope(), content, nil)
	} else if docs == nil {
		// query entity documents for ID
		return cpu.getDocuments(did, rMsg.Envelope(), command)
	}
	// received new documents
	return cpu.putDocuments(docs, did, rMsg.Envelope(), command)
}

func (cpu *DocumentCommandProcessor) getDocuments(did ID, envelope Envelope, content DocumentCommand) []Content {
	facebook := cpu.Facebook
	documents := facebook.GetDocuments(did)
	if len(documents) == 0 {
		return cpu.RespondReceipt("Document not found.", envelope, content, StringKeyMap{
			"template": "Document not found: ${did}",
			"replacements": StringKeyMap{
				"did": did.String(),
			},
		})
	}
	// documents got
	queryTime := content.LastTime()
	if !TimeIsNil(queryTime) {
		// check last document time
		lastDoc := cpu.lastDocument(documents)
		lastTime := lastDoc.Time()
		if TimeIsNil(lastTime) {
			//panic("document error")
		} else if TimeToInt64(lastTime) <= TimeToInt64(queryTime) {
			// document not updated
			return cpu.RespondReceipt("Document not updated.", envelope, content, StringKeyMap{
				"template": "Document not updated: ${did}",
				"replacements": StringKeyMap{
					"did": did.String(),
				},
			})
		}
	}
	// document got
	return cpu.respondDocuments(did, documents, envelope.Sender())
}

// protected
func (cpu *DocumentCommandProcessor) respondDocuments(did ID, documents []Document, receiver ID) []Content {
	if receiver.Equal(did) {
		panic("cycled response: " + receiver.String())
	}
	// TODO: check response expired
	facebook := cpu.Facebook
	meta := facebook.GetMeta(did)
	res := NewCommandForRespondDocuments(did, meta, documents)
	return []Content{res}
}

// protected
func (cpu *DocumentCommandProcessor) lastDocument(documents []Document) Document {
	var lastDoc Document
	var lastTime Time
	var docTime Time
	for _, doc := range documents {
		docTime = doc.Time()
		if lastDoc == nil {
			// first document
			lastDoc = doc
			lastTime = docTime
		} else if lastTime == nil {
			// the first document has no time (old version),
			// if this document has time, use the new one
			if !TimeIsNil(docTime) {
				// first document with time
				lastDoc = doc
				lastTime = docTime
			}
		} else if docTime != nil && TimeToInt64(docTime) > TimeToInt64(lastTime) {
			// new document
			lastDoc = doc
			lastTime = docTime
		}
	}
	return lastDoc
}

func (cpu *DocumentCommandProcessor) putDocuments(documents []Document, did ID, envelope Envelope, content DocumentCommand) []Content {
	var errors []Content
	meta := content.Meta()
	// 0. check meta
	if meta == nil {
		facebook := cpu.Facebook
		meta = facebook.GetMeta(did)
		if meta == nil {
			return cpu.RespondReceipt("Meta not found.", envelope, content, StringKeyMap{
				"template": "Meta not found: ${did}.",
				"replacements": StringKeyMap{
					"did": did.String(),
				},
			})
		}
	} else {
		// 1. try to save meta
		errors = cpu.saveMeta(meta, did, envelope, content)
		if errors != nil {
			// failed
			return errors
		}
	}
	// 2. try to save documents
	errors = make([]Content, 0)
	var responses []Content
	for _, doc := range documents {
		responses = cpu.saveDocument(doc, meta, did, envelope, content)
		if responses != nil {
			// failed
			errors = append(errors, responses...)
		}
	}
	if len(errors) > 0 {
		// failed
		return errors
	}
	// 3. success
	return cpu.RespondReceipt("Documents received.", envelope, content, StringKeyMap{
		"template": "Documents received: ${did}",
		"replacements": StringKeyMap{
			"did": did.String(),
		},
	})
}

// protected
func (cpu *DocumentCommandProcessor) saveDocument(doc Document, meta Meta, did ID, envelope Envelope, content DocumentCommand) []Content {
	facebook := cpu.Facebook
	// check document
	if !cpu.checkDocument(doc, meta, did) {
		// document invalid
		return cpu.RespondReceipt("Document not accepted.", envelope, content, StringKeyMap{
			"template": "Document not accepted: ${did}",
			"replacements": StringKeyMap{
				"did": did.String(),
			},
		})
	} else if !facebook.SaveDocument(doc, did) {
		// document expired
		return cpu.RespondReceipt("Document not changed.", envelope, content, StringKeyMap{
			"template": "Document not changed: ${did}",
			"replacements": StringKeyMap{
				"did": did.String(),
			},
		})
	}
	// document saved, return no error
	return nil
}

// protected
func (cpu *DocumentCommandProcessor) checkDocument(doc Document, meta Meta, did ID) bool {
	// check meta with ID
	if !cpu.checkMeta(meta, did) {
		// meta error
		return false
	}
	// check document ID
	helper := GetGeneralAccountHelper()
	docID := helper.GetDocumentID(doc.Map())
	if docID == nil {
		//panic("document ID not found")
	} else if !docID.Address().Equal(did.Address()) {
		//panic("document ID not matched")
		return false
	}
	// NOTICE: if this is a bulletin document for group,
	//             verify it with the group owner's meta.key
	//         else (this is a visa document for user)
	//             verify it with the user's meta.key
	pubKey := meta.PublicKey()
	return doc.Verify(pubKey)
	// TODO: check for group document
}
