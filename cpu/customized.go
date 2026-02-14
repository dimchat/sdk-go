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

/**
 *  Handler for CustomizedContent
 */
type CustomizedContentHandler interface {

	/**
	 *  Do your job
	 *
	 * @param content   - customized content
	 * @param rMsg      - network message
	 * @param messenger - message transceiver
	 * @return responses
	 */
	HandleContent(content CustomizedContent, rMsg ReliableMessage, messenger IMessenger) []Content
}

/**
 *  Default Handler
 */
type BaseCustomizedHandler struct {
	//CustomizedContentHandler
}

// Override
func (handler BaseCustomizedHandler) HandleContent(content CustomizedContent, rMsg ReliableMessage, _ IMessenger) []Content {
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

// protected
func (handler BaseCustomizedHandler) RespondReceipt(text string, head Envelope, body Content, extra StringKeyMap) []Content {
	// create base receipt command with text & original envelope
	res := createReceipt(text, head, body, extra)
	return []Content{res}
}

/**
 *  Filter for CustomizedContent Handler
 */
type CustomizedContentFilter interface {

	/**
	 *  Fetch a handler for 'app' and 'mod'
	 *
	 * @param content - customized content
	 * @param rMsg    - message with envelope
	 * @return customized handler
	 */
	FilterContent(content CustomizedContent, rMsg ReliableMessage) CustomizedContentHandler
}

type defaultCustomizedFilter struct {
	//CustomizedContentFilter

	// protected
	DefaultHandler CustomizedContentHandler
}

// Override
func (filter defaultCustomizedFilter) FilterContent(_ CustomizedContent, _ ReliableMessage) CustomizedContentHandler {
	// if the application has too many modules, I suggest you to
	// use different handler to do the jobs for each module.
	return filter.DefaultHandler
}

var sharedCustomizedContentFilter CustomizedContentFilter = &defaultCustomizedFilter{
	DefaultHandler: &BaseCustomizedHandler{},
}

func SetCustomizedContentFilter(filter CustomizedContentFilter) {
	sharedCustomizedContentFilter = filter
}

func GetCustomizedContentFilter() CustomizedContentFilter {
	return sharedCustomizedContentFilter
}

/**
 *  Customized Content Processing Unit
 *  <p>
 *      Handle content for application customized
 *  </p>
 */
type CustomizedContentProcessor struct {
	BaseContentProcessor
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
