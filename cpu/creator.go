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
type ContentProcessorCreator interface {

	/**
	 *  Create content processor with type
	 *
	 * @param msgType - content type
	 * @return ContentProcessor
	 */
	CreateContentProcessor(msgType ContentType) ContentProcessor

	/**
	 *  Create command processor with name
	 *
	 * @param msgType - content type
	 * @param cmdName - command name
	 * @return CommandProcessor
	 */
	CreateCommandProcessor(msgType ContentType, cmdName string) ContentProcessor
}

type BaseCreator struct {
	TwinsHelper
}

//func (creator *BaseCreator) Init(facebook IFacebook, messenger IMessenger) ContentProcessorCreator {
//	if creator.TwinsHelper.Init(facebook, messenger) != nil {
//	}
//	return creator
//}

//-------- IContentProcessorCreator

func (creator *BaseCreator) CreateContentProcessor(msgType ContentType) ContentProcessor {
	switch msgType {
	// forward content
	case FORWARD:
		return NewForwardContentProcessor(creator.Facebook(), creator.Messenger())
	// default commands
	case COMMAND:
		return NewBaseCommandProcessor(creator.Facebook(), creator.Messenger())
	case HISTORY:
		return NewHistoryCommandProcessor(creator.Facebook(), creator.Messenger())
	// default contents
	case 0:
		return NewBaseContentProcessor(creator.Facebook(), creator.Messenger())
	// unknown
	default:
		return nil
	}
}

func (creator *BaseCreator) CreateCommandProcessor(_ ContentType, cmdName string) ContentProcessor {
	switch cmdName {
	// meta command
	case META:
		return NewMetaCommandProcessor(creator.Facebook(), creator.Messenger())
	// document command
	case DOCUMENT:
		return NewDocumentCommandProcessor(creator.Facebook(), creator.Messenger())
	// group commands
	case "group":
		return NewGroupCommandProcessor(creator.Facebook(), creator.Messenger())
	case INVITE:
		return NewInviteCommandProcessor(creator.Facebook(), creator.Messenger())
	case EXPEL:
		return NewExpelCommandProcessor(creator.Facebook(), creator.Messenger())
	case QUIT:
		return NewQuitCommandProcessor(creator.Facebook(), creator.Messenger())
	case RESET:
		return NewResetCommandProcessor(creator.Facebook(), creator.Messenger())
	case QUERY:
		return NewQueryCommandProcessor(creator.Facebook(), creator.Messenger())
	// unknown
	default:
		return nil
	}
}
