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
	. "github.com/dimchat/sdk-go/dkd"
	. "github.com/dimchat/sdk-go/sdk"
)

/**
 *  Base ContentProcessor Creator
 */
type BaseCreator struct {
	//ContentProcessorCreator
	TwinsHelper
}

//func (creator *BaseCreator) Init(facebook IFacebook, messenger IMessenger) ContentProcessorCreator {
//	if creator.TwinsHelper.Init(facebook, messenger) != nil {
//	}
//	return creator
//}

// Override
func (creator *BaseCreator) CreateContentProcessor(msgType MessageType) ContentProcessor {
	switch msgType {
	// forward content
	case ContentType.FORWARD:
		return NewForwardContentProcessor(creator.Facebook, creator.Messenger)
	// array content
	case ContentType.ARRAY:
		return NewArrayContentProcessor(creator.Facebook, creator.Messenger)
	// default commands
	case ContentType.COMMON:
		return NewBaseCommandProcessor(creator.Facebook, creator.Messenger)
	// default contents
	case ContentType.ANY:
		return NewBaseContentProcessor(creator.Facebook, creator.Messenger)
	case "*":
		// must return a default processor for unknown type
		return NewBaseContentProcessor(creator.Facebook, creator.Messenger)
	// unknown
	default:
		return nil
	}
}

func (creator *BaseCreator) CreateCommandProcessor(_ MessageType, cmdName string) ContentProcessor {
	switch cmdName {
	// unknown
	default:
		return nil
	}
}

//
//  Initialize base creator for CPU factory
//

type cpuHelper struct {
	//ContentProcessorHelper
}

// Override
func (helper cpuHelper) CreateContentProcessorFactory(facebook IFacebook, messenger IMessenger) ContentProcessorFactory {
	creator := &BaseCreator{
		TwinsHelper{
			Facebook:  facebook,
			Messenger: messenger,
		},
	}
	return NewContentProcessorFactory(creator)
}

func init() {
	helper := &cpuHelper{}
	SetContentProcessorHelper(helper)
}
