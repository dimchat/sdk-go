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
	. "github.com/dimchat/mkm-go/types"
	. "github.com/dimchat/sdk-go/sdk"
)

/**
 *  CPU - Content Processing Unit
 */

// BaseContentProcessor is the base implementation of ContentProcessor
//
// Provides default content processing behavior and utility methods for response generation
// Serves as the foundation for all specialized content processors (text, file, image, etc.)
type BaseContentProcessor struct {
	//ContentProcessor
	*TwinsHelper
}

// Override
func (cpu *BaseContentProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	return cpu.RespondReceipt("Content not support.", rMsg.Envelope(), content, StringKeyMap{
		"template": "Content (type: ${type}) not support yet!",
		"replacements": StringKeyMap{
			"type": content.Type(),
		},
	})
}

// -------------------------------------------------------------------------
//  Response Utility Methods (Protected - for internal/extension use)
// -------------------------------------------------------------------------

// RespondReceipt is a protected utility method to generate standardized receipt responses
//
// Creates a ReceiptCommand with the specified text, context, and template parameters
// Used by base and specialized processors to generate consistent response formatting
//
// Parameters:
//   - text  - Human-readable response text (fallback if template is used)
//   - head  - Original message envelope (for response routing/context)
//   - body  - Original content being processed (for correlation)
//   - extra - Template parameters and additional metadata (supports ${var} replacement)
//
// Returns: Slice containing a single ReceiptCommand response
func (cpu *BaseContentProcessor) RespondReceipt(text string, head Envelope, body Content, extra StringKeyMap) []Content {
	// create base receipt command with text & original envelope
	res := createReceipt(text, head, body, extra)
	return []Content{res}
}

// createReceipt is a helper function to build a complete ReceiptCommand with context
//
// Populates receipt with original message context (group ID, envelope) and extra metadata
// Automatically adds group ID from content if present and merges extra key-value pairs
//
// Parameters:
//   - text  - Base text message for the receipt
//   - head  - Original message envelope (for sender/receiver/timestamp context)
//   - body  - Original content (used to extract group ID if present)
//   - extra - Additional key-value metadata (template params, custom fields)
//
// Returns: Fully populated ReceiptCommand with context and metadata
func createReceipt(text string, head Envelope, body Content, extra StringKeyMap) ReceiptCommand {
	// create base receipt command with text, original envelope, serial number & group ID
	res := NewReceiptCommand(text, head, body)
	if body != nil {
		// check group
		group := body.Group()
		if group != nil {
			res.SetGroup(group)
		}
	}
	// add extra key-value
	if extra != nil {
		for key, value := range extra {
			res.Set(key, value)
		}
	}
	return res
}

/**
 *  CPU - Command Processing Unit
 */

// BaseCommandProcessor is the base implementation for command-specific ContentProcessors
//
// Extends BaseContentProcessor with command-specific default behavior
// Handles type checking for Command content and provides "command not supported" fallback
type BaseCommandProcessor struct {
	*BaseContentProcessor
}

// Override
func (cpu *BaseCommandProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	command, ok := content.(Command)
	if !ok {
		//panic("command error")
		return nil
	}
	return cpu.RespondReceipt("Command not support.", rMsg.Envelope(), content, StringKeyMap{
		"template": "Command (name: ${command}) not support yet!",
		"replacements": StringKeyMap{
			"command": command.CMD(),
		},
	})
}
