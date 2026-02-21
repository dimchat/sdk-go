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
)

type ContentProcessorMap = map[string]ContentProcessor

func NewContentProcessorMap(capacity int) ContentProcessorMap {
	return make(ContentProcessorMap, capacity)
}

// General ContentProcessor Factory
type GeneralContentProcessorFactory struct {
	//ContentProcessorFactory

	creator ContentProcessorCreator

	contentProcessors ContentProcessorMap
	commandProcessors ContentProcessorMap
}

func NewGeneralContentProcessorFactory(creator ContentProcessorCreator) *GeneralContentProcessorFactory {
	return &GeneralContentProcessorFactory{
		creator:           creator,
		contentProcessors: NewContentProcessorMap(32),
		commandProcessors: NewContentProcessorMap(32),
	}
}

// Override
func (factory *GeneralContentProcessorFactory) GetContentProcessor(content Content) ContentProcessor {
	msgType := content.Type()
	if command, ok := content.(Command); ok {
		cmd := command.CMD()
		cpu := factory.GetCommandProcessor(msgType, cmd)
		if cpu != nil {
			return cpu
		} else if _, ok = content.(GroupCommand); ok {
			// group command processor
			cpu = factory.GetCommandProcessor(msgType, "group")
			if cpu != nil {
				return cpu
			}
		}
	}
	// content processor
	return factory.GetContentProcessorForType(msgType)
}

// Override
func (factory *GeneralContentProcessorFactory) GetContentProcessorForType(msgType MessageType) ContentProcessor {
	cpu := factory.contentProcessors[msgType]
	if cpu == nil {
		cpu = factory.creator.CreateContentProcessor(msgType)
		if cpu != nil {
			factory.contentProcessors[msgType] = cpu
		}
	}
	return cpu
}

// private
func (factory *GeneralContentProcessorFactory) GetCommandProcessor(msgType MessageType, cmdName string) ContentProcessor {
	cpu := factory.commandProcessors[cmdName]
	if cpu == nil {
		cpu = factory.creator.CreateCommandProcessor(msgType, cmdName)
		if cpu != nil {
			factory.commandProcessors[cmdName] = cpu
		}
	}
	return cpu
}
