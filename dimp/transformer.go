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
	. "github.com/dimchat/core-go/core"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

type MessengerTransformer struct {
	MessageTransformer
}

func (transformer *MessengerTransformer) Init(messenger IMessenger) *MessengerTransformer {
	if transformer.MessageTransformer.Init(messenger) != nil {
	}
	return transformer
}

func (transformer *MessengerTransformer) Messenger() IMessenger {
	return transformer.Transceiver().(IMessenger)
}

func (transformer *MessengerTransformer) Facebook() IFacebook {
	return transformer.Messenger().Facebook()
}

func (transformer *MessengerTransformer) getFileContentProcessor() FileContentProcessor {
	processor := ContentProcessorGetByType(FILE)
	processor.SetMessenger(transformer.Messenger())
	return processor.(FileContentProcessor)
}

//-------- InstantMessageDelegate

func (transformer *MessengerTransformer) SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte {
	// check attachment for File/Image/Audio/Video message content
	file, ok := content.(FileContent)
	if ok {
		fpu := transformer.getFileContentProcessor()
		fpu.UploadFileContent(file, password, iMsg)
	}
	return transformer.MessageTransformer.SerializeContent(content, password, iMsg)
}

func (transformer *MessengerTransformer) EncryptKey(data []byte, receiver ID, iMsg InstantMessage) []byte {
	key := transformer.Facebook().GetPublicKeyForEncryption(receiver)
	if key == nil {
		// save this message in a queue waiting receiver's meta/document response
		transformer.Messenger().SuspendInstantMessage(iMsg)
		return nil
	}
	return transformer.MessageTransformer.EncryptKey(data, receiver, iMsg)
}

//-------- SecureMessageDelegate

func (transformer *MessengerTransformer) DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content {
	content := transformer.MessageTransformer.DeserializeContent(data, password, sMsg)
	// check attachment for File/Image/Audio/Video message content
	file, ok := content.(FileContent)
	if ok {
		fpu := transformer.getFileContentProcessor()
		fpu.DownloadFileContent(file, password, sMsg)
	}
	return content
}
