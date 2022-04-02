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
type BaseProcessorCreator struct {
	TwinsHelper
}

//-------- IProcessorCreator

func (factory *BaseProcessorCreator) CreateContentProcessor(msgType ContentType) ContentProcessor {
	// core contents
	if msgType == FORWARD {
		return NewForwardContentProcessor(factory.Facebook(), factory.Messenger())
	}
	// unknown
	return nil
}

func (factory *BaseProcessorCreator) CreateCommandProcessor(msgType ContentType, cmdName string) CommandProcessor {
	// meta
	if cmdName == META {
		return NewMetaCommandProcessor(factory.Facebook(), factory.Messenger())
	}
	// document
	if cmdName == DOCUMENT {
		return NewDocumentCommandProcessor(factory.Facebook(), factory.Messenger())
	} else if cmdName == "profile" || cmdName == "visa" || cmdName == "bulletin" {
		return NewDocumentCommandProcessor(factory.Facebook(), factory.Messenger())
	}
	// group
	if cmdName == "group" {
		return NewGroupCommandProcessor(factory.Facebook(), factory.Messenger())
	} else if cmdName == INVITE {
		return NewInviteCommandProcessor(factory.Facebook(), factory.Messenger())
	} else if cmdName == EXPEL {
		return NewExpelCommandProcessor(factory.Facebook(), factory.Messenger())
	} else if cmdName == QUIT {
		return NewQuitCommandProcessor(factory.Facebook(), factory.Messenger())
	} else if cmdName == RESET {
		return NewResetCommandProcessor(factory.Facebook(), factory.Messenger())
	} else if cmdName == QUERY {
		return NewQueryCommandProcessor(factory.Facebook(), factory.Messenger())
	}
	// others
	if msgType == COMMAND {
		return NewCommandProcessor(factory.Facebook(), factory.Messenger())
	} else if msgType == HISTORY {
		return NewHistoryCommandProcessor(factory.Facebook(), factory.Messenger())
	}
	// unknown
	return nil
}
