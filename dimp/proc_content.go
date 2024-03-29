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
	. "github.com/dimchat/dkd-go/protocol"
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

/**
 *  CPU Factory
 *  ~~~~~~~~~~~
 *
 *  Delegate for Message Processor
 */
type ContentProcessorFactory interface {

	/**
	 *  Get processor for content
	 */
	GetProcessor(content Content) ContentProcessor

	/**
	 *  Get processor for content type
	 */
	GetContentProcessor(msgType ContentType) ContentProcessor

	/**
	 *  Get processor for command name
	 */
	GetCommandProcessor(msgType ContentType, cmdName string) ContentProcessor
}
