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

import . "github.com/dimchat/dkd-go/protocol"

/**
 *  CPU: Content Processing Unit
 */
type ContentProcessor interface {

	/**
	 *  Process message content
	 *
	 * @param content - content received
	 * @param rMsg    - reliable message
	 * @return {Content} response to sender
	 */
	ProcessContent(content Content, rMsg ReliableMessage) []Content
}

/**
 *  CPU Creator
 */
type ContentProcessorCreator interface {

	/**
	 *  Create content processor with type
	 *
	 * @param msgType - content type
	 * @return ContentProcessor
	 */
	CreateContentProcessor(msgType MessageType) ContentProcessor

	/**
	 *  Create command processor with name
	 *
	 * @param msgType - content type
	 * @param cmdName - command name
	 * @return CommandProcessor
	 */
	CreateCommandProcessor(msgType MessageType, cmdName string) ContentProcessor
}

/**
 *  CPU Factory
 */
type ContentProcessorFactory interface {

	/**
	 *  Get content/command processor
	 *
	 * @param content - Content/Command
	 * @return ContentProcessor
	 */
	GetContentProcessor(content Content) ContentProcessor

	GetContentProcessorForType(msgType MessageType) ContentProcessor
}
