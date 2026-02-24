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

// ContentProcessor defines the core interface for processing message content
//
// Core responsibility: Handle business logic for specific types of message content (text, command, file, etc.)
// Processes incoming content and generates appropriate response content for the sender
type ContentProcessor interface {

	// ProcessContent executes business logic on incoming message content and returns response content
	//
	// Uses context from the ReliableMessage (sender/receiver/...) to inform processing logic
	//
	// Parameters:
	//   - content - Raw incoming message content to process (text/command/file/etc.)
	//   - rMsg    - Parent ReliableMessage providing context (sender ID, receiver ID, timestamp, etc.)
	// Returns: Slice of response Content objects (empty slice if no response needed)
	ProcessContent(content Content, rMsg ReliableMessage) []Content
}

/**
 *  CPU Creator
 */

// ContentProcessorCreator defines the factory interface for creating ContentProcessor instances
//
// Acts as a "CPU Creator" - instantiates specialized processors for different content/command types
// Supports both generic content type creation and specific command name creation
type ContentProcessorCreator interface {

	// CreateContentProcessor instantiates a ContentProcessor for a specific message type
	//
	// Creates generic processors for broad content categories (text, file, image, etc.)
	//
	// Parameters:
	//   - msgType - Target message/content type to create a processor for
	// Returns: Specialized ContentProcessor (nil if no processor exists for the type)
	CreateContentProcessor(msgType MessageType) ContentProcessor

	// CreateCommandProcessor instantiates a ContentProcessor for a specific command
	//
	// Creates specialized processors for named commands (e.g., "meta", "documents", "receipt", ...)
	//
	// Parameters:
	//   - msgType - Base message type (typically "command" type)
	//   - cmd     - Specific command name to create a processor for
	// Returns: Specialized CommandProcessor (implements ContentProcessor, nil if command not found)
	CreateCommandProcessor(msgType MessageType, cmd string) ContentProcessor
}

/**
 *  CPU Factory
 */

// ContentProcessorFactory defines the factory interface for retrieving pre-created ContentProcessor instances
//
// Acts as a "CPU Factory" - manages a cache of processors and provides quick access
// Eliminates redundant processor creation by reusing existing instances
type ContentProcessorFactory interface {

	// GetContentProcessor retrieves the appropriate ContentProcessor for a given Content instance
	//
	// Determines the content type/command from the Content object and returns the matching processor
	//
	// Parameters:
	//   - content - Incoming Content to find a processor for
	// Returns: Matching ContentProcessor (nil if no processor exists for the content type/command)
	GetContentProcessor(content Content) ContentProcessor

	// GetContentProcessorForType retrieves the ContentProcessor for a specific MessageType
	//
	// Provides direct access to generic processors by content type (without parsing Content object)
	//
	// Parameters:
	//   - msgType - Target message/content type to get a processor for
	// Returns: Generic ContentProcessor for the type (nil if no processor exists)
	GetContentProcessorForType(msgType MessageType) ContentProcessor
}
