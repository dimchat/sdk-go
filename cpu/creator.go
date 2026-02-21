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
type BaseContentProcessorCreator struct {
	//ContentProcessorCreator
	*TwinsHelper
}

func NewBaseContentProcessorCreator(facebook Facebook, messenger Messenger) *BaseContentProcessorCreator {
	return &BaseContentProcessorCreator{
		TwinsHelper: NewTwinsHelper(facebook, messenger),
	}
}

// Override
func (creator *BaseContentProcessorCreator) CreateContentProcessor(msgType MessageType) ContentProcessor {
	switch msgType {
	// application customized
	case ContentType.APPLICATION, ContentType.CUSTOMIZED:
		return NewCustomizedContentProcessor(creator.Facebook, creator.Messenger)
	// forward content
	case ContentType.FORWARD:
		return NewForwardContentProcessor(creator.Facebook, creator.Messenger)
	// array content
	case ContentType.ARRAY:
		return NewArrayContentProcessor(creator.Facebook, creator.Messenger)
	// default commands
	case ContentType.COMMAND:
		return NewBaseCommandProcessor(creator.Facebook, creator.Messenger)
	// default contents
	case ContentType.ANY:
		// must return a default processor for unknown type
		return NewBaseContentProcessor(creator.Facebook, creator.Messenger)
	// unknown
	default:
		//panic("unsupported content type")
		return nil
	}
}

// Override
func (creator *BaseContentProcessorCreator) CreateCommandProcessor(_ MessageType, cmdName string) ContentProcessor {
	switch cmdName {
	// meta command
	case META:
		return NewMetaCommandProcessor(creator.Facebook, creator.Messenger)
	// documents command
	case DOCUMENTS:
		return NewDocumentCommandProcessor(creator.Facebook, creator.Messenger)
	// unknown
	default:
		//panic("unsupported command: " + cmdName)
		return nil
	}
}

//
//  Factories
//

func NewBaseContentProcessor(facebook Facebook, messenger Messenger) *BaseContentProcessor {
	return &BaseContentProcessor{
		TwinsHelper: NewTwinsHelper(facebook, messenger),
	}
}

func NewBaseCommandProcessor(facebook Facebook, messenger Messenger) *BaseCommandProcessor {
	return &BaseCommandProcessor{
		BaseContentProcessor: NewBaseContentProcessor(facebook, messenger),
	}
}

func NewForwardContentProcessor(facebook Facebook, messenger Messenger) *ForwardContentProcessor {
	return &ForwardContentProcessor{
		BaseContentProcessor: NewBaseContentProcessor(facebook, messenger),
	}
}

func NewArrayContentProcessor(facebook Facebook, messenger Messenger) *ArrayContentProcessor {
	return &ArrayContentProcessor{
		BaseContentProcessor: NewBaseContentProcessor(facebook, messenger),
	}
}

func NewMetaCommandProcessor(facebook Facebook, messenger Messenger) *MetaCommandProcessor {
	return &MetaCommandProcessor{
		BaseCommandProcessor: NewBaseCommandProcessor(facebook, messenger),
	}
}

func NewDocumentCommandProcessor(facebook Facebook, messenger Messenger) *DocumentCommandProcessor {
	return &DocumentCommandProcessor{
		MetaCommandProcessor: NewMetaCommandProcessor(facebook, messenger),
	}
}

func NewCustomizedContentProcessor(facebook Facebook, messenger Messenger) *CustomizedContentProcessor {
	return &CustomizedContentProcessor{
		BaseContentProcessor: NewBaseContentProcessor(facebook, messenger),
	}
}

//
//  Initialize base creator for CPU factory
//

type cpuHelper struct {
	//ContentProcessorHelper
}

// Override
func (cpuHelper) CreateContentProcessorFactory(facebook Facebook, messenger Messenger) ContentProcessorFactory {
	creator := NewBaseContentProcessorCreator(facebook, messenger)
	return NewGeneralContentProcessorFactory(creator)
}

func init() {
	helper := &cpuHelper{}
	SetContentProcessorHelper(helper)
}
