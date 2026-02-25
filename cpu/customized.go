/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/sdk"
)

// CustomizedContentHandler defines the interface for processing application-specific CustomizedContent
//
// # Core extension point for applications to implement custom business logic for their proprietary message types
//
// Handles content with app/module/action identifiers (app-specific structured messages)
type CustomizedContentHandler interface {

	// HandleContent executes application-specific logic for CustomizedContent
	//
	// Implements business logic for the given app/module/action combination
	//
	// Parameters:
	//   - content   - Application-specific CustomizedContent to process (contains app/mod/act identifiers)
	//   - rMsg      - Parent ReliableMessage providing envelope/context (sender/receiver/timestamp)
	//   - messenger - Messenger instance for message transformation/processing utilities
	// Returns: Slice of response Content objects (empty slice if no response needed)
	HandleContent(content CustomizedContent, rMsg ReliableMessage, messenger Messenger) []Content
}

// BaseCustomizedHandler provides a default implementation of CustomizedContentHandler
//
// Serves as fallback for unimplemented app/module/action combinations
//
// Returns standardized "content not supported" receipt responses with app/mod/act context
type BaseCustomizedHandler struct {
	//CustomizedContentHandler
}

// Override
func (handler BaseCustomizedHandler) HandleContent(content CustomizedContent, rMsg ReliableMessage, _ Messenger) []Content {
	app := content.Application()
	mod := content.Module()
	act := content.Action()
	return handler.RespondReceipt("Content not support.", rMsg.Envelope(), content, StringKeyMap{
		"template": "Customized content (app: ${app}, mod: ${mod}, act: ${act}) not support yet!",
		"replacements": StringKeyMap{
			"app": app,
			"mod": mod,
			"act": act,
		},
	})
}

// RespondReceipt is a protected utility method to generate standardized receipt responses
//
// Creates a ReceiptCommand with app-specific context (used by base and custom handlers)
//
// Parameters:
//   - text  - Human-readable response text (fallback for template)
//   - head  - Original message envelope (for response routing)
//   - body  - Original CustomizedContent (for correlation)
//   - extra - Template parameters (supports ${app}/${mod}/${act} replacement)
//
// Returns: Slice containing a single ReceiptCommand response
func (handler BaseCustomizedHandler) RespondReceipt(text string, head Envelope, body Content, extra StringKeyMap) []Content {
	// create base receipt command with text & original envelope
	res := createReceipt(text, head, body, extra)
	return []Content{res}
}

// CustomizedContentFilter defines the interface for routing CustomizedContent to appropriate handlers
//
// Acts as a handler resolver - maps app/module combinations to specific CustomizedContentHandler instances
// Core extension point for implementing handler routing logic (e.g., by app ID, module ID, or action)
type CustomizedContentFilter interface {

	// FilterContent resolves the appropriate CustomizedContentHandler for a given CustomizedContent
	//
	// Implements routing logic to select handler based on app/mod/act or message context
	//
	// Parameters:
	//   - content - CustomizedContent to find a handler for (contains app/mod/act identifiers)
	//   - rMsg    - Parent ReliableMessage providing additional context for routing
	// Returns: CustomizedContentHandler for the content (nil if no handler found)
	FilterContent(content CustomizedContent, rMsg ReliableMessage) CustomizedContentHandler
}

// defaultCustomizedFilter is the default implementation of CustomizedContentFilter
//
// Provides simple routing logic that returns the same default handler for all content
// Recommended to replace with application-specific filter for multi-module applications
type defaultCustomizedFilter struct {
	//CustomizedContentFilter

	// defaultHandler is the fallback handler returned for all content types
	// Used when no specific handler is found for an app/module combination
	defaultHandler CustomizedContentHandler
}

// Override
func (filter defaultCustomizedFilter) FilterContent(_ CustomizedContent, _ ReliableMessage) CustomizedContentHandler {
	// if the application has too many modules, I suggest you to
	// use different handler to do the jobs for each module.
	return filter.defaultHandler
}

var sharedCustomizedContentFilter CustomizedContentFilter = &defaultCustomizedFilter{
	defaultHandler: &BaseCustomizedHandler{},
}

func SetCustomizedContentFilter(filter CustomizedContentFilter) {
	sharedCustomizedContentFilter = filter
}

func GetCustomizedContentFilter() CustomizedContentFilter {
	return sharedCustomizedContentFilter
}

// CustomizedContentProcessor is the ContentProcessor implementation for CustomizedContent
//
// # Integrates with CustomizedContentFilter and CustomizedContentHandler to process app-specific content
//
// Serves as the bridge between the core message processing system and application-specific logic
type CustomizedContentProcessor struct {
	*BaseContentProcessor
}

// Override
func (cpu *CustomizedContentProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	customized, ok := content.(CustomizedContent)
	if !ok {
		//panic("customized content error")
		return nil
	}
	// get handler for 'app' & 'mod'
	filter := GetCustomizedContentFilter()
	handler := filter.FilterContent(customized, rMsg)
	if handler == nil {
		//panic("should not happen")
		return nil
	}
	// handle the action
	messenger := cpu.Messenger
	return handler.HandleContent(customized, rMsg, messenger)
}
