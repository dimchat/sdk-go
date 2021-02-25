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
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	"strings"
)

type FileContentProcessor interface {
	ContentProcessor
	IFileContentProcessor
}
type IFileContentProcessor interface {

	UploadFileContent(content FileContent, pwd SymmetricKey, iMsg InstantMessage) bool
	DownloadFileContent(content FileContent, pwd SymmetricKey, sMsg SecureMessage) bool
}

//
//  File content processor implementation
//
type BaseFileContentProcessor struct {
	BaseContentProcessor
}

func (cpu *BaseFileContentProcessor) Init() *BaseFileContentProcessor {
	if cpu.BaseContentProcessor.Init() != nil {
	}
	return cpu
}

func (cpu *BaseFileContentProcessor) UploadFileContent(content FileContent, pwd SymmetricKey, iMsg InstantMessage) bool {
	data := content.Data()
	if data == nil || len(data) == 0 {
		panic("failed to get file data")
		return false
	}
	// encrypt and upload file data onto CDN and save the URL in message content
	encrypted := pwd.Encrypt(data)
	if encrypted == nil || len(encrypted) == 0 {
		panic("failed to encrypt file data with key")
		return false
	}
	url := cpu.Messenger().UploadData(encrypted, iMsg)
	if url == "" {
		return false
	} else {
		// replace 'data' with 'URL'
		content.SetURL(url)
		content.SetData(nil)
		return true
	}
}

func (cpu *BaseFileContentProcessor) DownloadFileContent(content FileContent, pwd SymmetricKey, sMsg SecureMessage) bool {
	url := content.URL()
	if !strings.Contains(url, "://") {
		// download URL not found
		return false
	}
	iMsg := InstantMessageCreate(sMsg.Envelope(), content)
	// download from CDN
	encrypted := cpu.Messenger().DownloadData(url, iMsg)
	if encrypted == nil || len(encrypted) == 0 {
		// save symmetric key for decrypting file data after download from CDN
		content.SetPassword(pwd)
		return false
	} else {
		// decrypt file data
		plain := pwd.Decrypt(encrypted)
		if plain == nil || len(plain) == 0 {
			panic("failed to decrypt file data with key")
			return false
		}
		content.SetData(plain)
		content.SetURL("")
		return true
	}
}

func (cpu *BaseFileContentProcessor) Process(_ Content, _ ReliableMessage) Content {
	// TODO: process file content

	return nil
}
