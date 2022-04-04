/* license: https://mit-license.org
 *
 *  DIM-SDK : Decentralized Instant Messaging Software Development Kit
 *
 *                                Written in 2022 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Albert Moky
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
	. "github.com/dimchat/sdk-go/dimp"
)

/**
 *  CPU Creator
 *  ~~~~~~~~~~~
 *
 *  Delegate for CPU factory
 */
type CPCreator struct {
	TwinsHelper
}

//-------- IContentProcessorCreator

func (factory *CPCreator) CreateContentProcessor(msgType ContentType) ContentProcessor {
	switch msgType {
	// forward content
	case FORWARD:
		return NewForwardContentProcessor(factory.Facebook(), factory.Messenger())
	// default commands
	case COMMAND:
		return NewBaseCommandProcessor(factory.Facebook(), factory.Messenger())
	case HISTORY:
		return NewHistoryCommandProcessor(factory.Facebook(), factory.Messenger())
	// default contents
	case 0:
		return NewBaseContentProcessor(factory.Facebook(), factory.Messenger())
	// unknown
	default:
		return nil
	}
}

func (factory *CPCreator) CreateCommandProcessor(_ ContentType, cmdName string) ContentProcessor {
	switch cmdName {
	// meta command
	case META:
		return NewMetaCommandProcessor(factory.Facebook(), factory.Messenger())
	// document command
	case DOCUMENT:
		return NewDocumentCommandProcessor(factory.Facebook(), factory.Messenger())
	// group commands
	case "group":
		return NewGroupCommandProcessor(factory.Facebook(), factory.Messenger())
	case INVITE:
		return NewInviteCommandProcessor(factory.Facebook(), factory.Messenger())
	case EXPEL:
		return NewExpelCommandProcessor(factory.Facebook(), factory.Messenger())
	case QUIT:
		return NewQuitCommandProcessor(factory.Facebook(), factory.Messenger())
	case RESET:
		return NewResetCommandProcessor(factory.Facebook(), factory.Messenger())
	case QUERY:
		return NewQueryCommandProcessor(factory.Facebook(), factory.Messenger())
	// unknown
	default:
		return nil
	}
}
